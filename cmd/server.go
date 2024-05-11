package cmd

import (
	"context"
	"database/sql"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lmittmann/tint"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vkhobor/go-opencv/api"
	"github.com/vkhobor/go-opencv/config"
	pathutils "github.com/vkhobor/go-opencv/path"

	database "github.com/vkhobor/go-opencv/db"
	"github.com/vkhobor/go-opencv/domain"
	"github.com/vkhobor/go-opencv/jobs"

	"github.com/spf13/cobra"
)

func run(ctx context.Context, w io.Writer, args []string, programConfig config.ProgramConfig) error {
	logger := slog.New(tint.NewHandler(w, &tint.Options{
		Level: slog.LevelDebug,
	}))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.SetDefault(logger)

	slog.Info("Opening database", "file", programConfig.Db)
	path, err := pathutils.EnsurePath(programConfig.BlobStorage)
	if err != nil {
		return err
	}

	dbconn, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	driver, err := sqlite3.WithInstance(dbconn, &sqlite3.Config{})
	if err != nil {
		return err
	}

	slog.Info("Migrating database")
	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations",
		"sqlite3", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	slog.Info("Setup dependencies")
	queries := database.New(dbconn)
	domainQueries := domain.JobQueries{
		Queries: queries,
	}

	dirConfig, err := config.NewDirectoryConfig(programConfig.BlobStorage)
	if err != nil {
		return err
	}

	wake := make(chan struct{}, 1)
	jobManager := jobs.JobManager{
		Wake:           wake,
		AutoWakePeriod: time.Minute * 2,
		Scraper:        jobs.NewScraper(&domainQueries, time.Second*5),
		Importer:       jobs.NewImporter(&domainQueries, time.Second*0, dirConfig),
		Downloader:     jobs.NewDownloader(&domainQueries, time.Second*15),
	}

	slog.Info("Starting job manager")
	go jobManager.Start()
	wake <- struct{}{}

	portString := fmt.Sprintf(":%d", programConfig.Port)
	srv := &http.Server{Addr: portString, Handler: api.NewRouter(queries, wake, dirConfig)}
	slog.Info("Server started", "port", programConfig.Port)

	go func() {
		httpError := srv.ListenAndServe()
		if httpError != nil && httpError != http.ErrServerClosed {
			log.Fatal(httpError)
		}
	}()

	<-ctx.Done()

	gracefulTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	slog.Info("Shutting down server...")
	err = srv.Shutdown(gracefulTimeout)
	if err != nil {
		slog.Error("Error shutting down server", "error", err)
		return err
	}

	return nil
}

func NewRunserver() *cobra.Command {
	viperConf := config.MustNewDefaultViperConfig()

	var cmdPrint = &cobra.Command{
		Use: "serve",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)
			defer cancel()

			config := config.ProgramConfig{}
			err := viperConf.Unmarshal(&config)
			if err != nil {
				return err
			}
			if err := config.Validate(); err != nil {
				slog.Error("Invalid configuration", "config", config)
				return fmt.Errorf("invalid configuration")
			}

			slog.Info("Starting serve", "configuration", config)

			if err := run(ctx, os.Stdout, args, config); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			return nil
		},
	}

	cmdPrint.Flags().IntP("port", "p", 8080, "Specify the port")
	cmdPrint.MarkFlagRequired("port")
	viperConf.BindPFlag("port", cmdPrint.Flags().Lookup("port"))

	cmdPrint.Flags().StringP("db", "d", "~/test", "Address of the sqlite database")
	viperConf.BindPFlag("db", cmdPrint.Flags().Lookup("db"))

	cmdPrint.Flags().StringP("blob-storage", "s", "~/test", "Specify where to store files")
	viperConf.BindPFlag("blob_storage", cmdPrint.Flags().Lookup("blob-storage"))

	return cmdPrint
}
