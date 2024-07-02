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
	"github.com/vkhobor/go-opencv/background"
	"github.com/vkhobor/go-opencv/config"
	pathutils "github.com/vkhobor/go-opencv/path"
	"github.com/vkhobor/go-opencv/queries"
	"github.com/vkhobor/go-opencv/scraper"

	database "github.com/vkhobor/go-opencv/db"

	"github.com/spf13/cobra"
)

func run(ctx context.Context, w io.Writer, args []string, programConfig config.ProgramConfig) error {
	logger := slog.New(tint.NewHandler(w, &tint.Options{
		Level: slog.LevelDebug,
	}))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.SetDefault(logger)

	err := pathutils.EnsurePath(programConfig.Db, false)
	if err != nil {
		return err
	}

	slog.Info("Opening database", "file", programConfig.Db)
	dbconn, err := sql.Open("sqlite3", programConfig.Db)
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
	dbQueries := database.New(dbconn)

	highLevelQueries := queries.Queries{
		Queries: dbQueries,
	}

	dirConfig, err := config.NewDirectoryConfig(programConfig.BlobStorage)
	if err != nil {
		return err
	}

	scrapeArgsChan := make(chan queries.Job)
	scrapedVideoChan := make(chan queries.ScrapedVideo)
	downloadedChan := make(chan queries.DownlodedVideo)
	importedChan := make(chan queries.ImportedVideo)

	downloader := background.Downloader{
		Queries:  &highLevelQueries,
		Throttle: time.Second * 5,
		Config:   dirConfig,
		Input:    scrapedVideoChan,
		Output:   downloadedChan,
	}

	importer := background.Importer{
		Queries:  &highLevelQueries,
		Throttle: time.Second * 5,
		Config:   dirConfig,
		Input:    downloadedChan,
		Output:   importedChan,
	}

	scraperJob := background.ScraperJob{
		Scraper: scraper.Scraper{
			Throttle: time.Second * 5,
			Domain:   "yewtu.be",
		},
		Queries: &highLevelQueries,
		Input:   scrapeArgsChan,
		Output:  scrapedVideoChan,
		Config:  dirConfig,
	}

	jobManager := background.DbMonitor{
		Wake:          make(chan struct{}, 1),
		Queries:       &highLevelQueries,
		ScrapeInput:   scrapeArgsChan,
		DownloadInput: scrapedVideoChan,
		ImportInput:   downloadedChan,
	}

	slog.Info("Starting jobs")

	go scraperJob.Start()
	go downloader.Start()
	go importer.Start()

	go jobManager.Start()
	jobManager.Wake <- struct{}{}

	portString := fmt.Sprintf(":%d", programConfig.Port)
	router := api.NewRouter(dbQueries, jobManager.Wake, dirConfig, programConfig)
	srv := &http.Server{Addr: portString, Handler: router}
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
