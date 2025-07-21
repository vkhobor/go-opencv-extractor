package commands

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"time"

	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gopkg.in/natefinch/lumberjack.v2"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vkhobor/go-opencv/internal/api"
	"github.com/vkhobor/go-opencv/internal/background"
	"github.com/vkhobor/go-opencv/internal/config"
	"github.com/vkhobor/go-opencv/internal/mlog"
	pathutils "github.com/vkhobor/go-opencv/internal/path"
)

func RunServer(ctx context.Context, w io.Writer, args []string, programConfig config.ServerConfig) error {
	pathutils.MustEnsurePath(programConfig.LogFolder, true)
	logFile := &lumberjack.Logger{
		Filename:   programConfig.LogFolder + "/log.log",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
	}
	multiWriter := io.MultiWriter(w, logFile)

	logger := slog.New(slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	mlog.SetDefault(logger)

	// Separate handler for slog.
	// If any package uses slog internally we filter for LevelError
	loggerSlog := slog.New(slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	slog.SetLogLoggerLevel(slog.LevelError)
	slog.SetDefault(loggerSlog)

	err := pathutils.EnsurePath(programConfig.Db, false)
	if err != nil {
		return err
	}

	mlog.Log().Info("Opening database", "file", programConfig.Db)
	dbconn, err := sql.Open("sqlite3", programConfig.Db)
	if err != nil {
		return err
	}
	driver, err := sqlite3.WithInstance(dbconn, &sqlite3.Config{})
	if err != nil {
		return err
	}

	mlog.Log().Info("Migrating database")
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

	mlog.Log().Info("Setup dependencies")

	dirConfig, err := programConfig.GetDirectoryConfig()
	if err != nil {
		return err
	}

	downloadedChan := make(chan struct {
		ID       string
		JobID    string
		FilterID string
	}, 100)
	defer close(downloadedChan)
	wakeJobs := make(chan struct{}, 1)
	defer close(wakeJobs)

	jobManager := background.DbMonitor{
		Config:      dirConfig,
		Wake:        wakeJobs,
		SqlDB:       dbconn,
		ImportInput: downloadedChan,
	}

	mlog.Log().Info("Starting jobs")
	go jobManager.Start()
	jobManager.Wake <- struct{}{}

	portString := fmt.Sprintf(":%d", programConfig.Port)
	router := api.NewRouter(dbconn, jobManager.Wake, dirConfig, programConfig)
	srv := &http.Server{Addr: portString, Handler: router}
	mlog.Log().Info("Server started", "port", programConfig.Port)

	go func() {
		httpError := srv.ListenAndServe()
		if httpError != nil && httpError != http.ErrServerClosed {
			mlog.Log().Error("Cannot listen and serve", "httpError", httpError)
			panic(httpError)
		}
	}()

	<-ctx.Done()

	gracefulTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	mlog.Log().Info("Shutting down server...")
	err = srv.Shutdown(gracefulTimeout)
	if err != nil {
		mlog.Log().Error("Error shutting down server", "error", err)
		return err
	}

	return nil
}
