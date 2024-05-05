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

	_ "github.com/mattn/go-sqlite3"
	"github.com/vkhobor/go-opencv/api"
	"github.com/vkhobor/go-opencv/db_sql"
	"github.com/vkhobor/go-opencv/jobs"

	"github.com/spf13/cobra"
)

func run(ctx context.Context, w io.Writer, args []string, port int) error {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.SetDefault(logger)

	slog.Info("Opening database")
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		return err
	}
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}
	slog.Info("Migrating database")
	m, err := migrate.NewWithDatabaseInstance(
		"file://./db_sql/migrations",
		"sqlite3", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	slog.Info("Setup dependencies")
	queries := db_sql.New(db)

	wake := make(chan struct{}, 1)
	jobManager := jobs.JobManager{
		Queries:          queries,
		Wake:             wake,
		AutoWakePeriod:   time.Minute * 2,
		ScrapeThrottle:   time.Second * 5,
		ImportThrottle:   time.Second * 0,
		DownloadThrottle: time.Second * 15,
	}

	go jobManager.Run()
	wake <- struct{}{}

	portString := fmt.Sprintf(":%d", port)
	srv := &http.Server{Addr: portString, Handler: api.NewRouter(queries, wake)}
	slog.Info("Server started", "port", port)

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
	var port int

	var cmdPrint = &cobra.Command{
		Use: "serve",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)
			defer cancel()
			if err := run(ctx, os.Stdout, args, port); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			return nil
		},
	}

	cmdPrint.Flags().IntVarP(&port, "port", "p", 8080, "Specify the port")
	cmdPrint.MarkFlagRequired("port")

	return cmdPrint
}
