package cmd

import (
	"context"
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vkhobor/go-opencv/api"
	"github.com/vkhobor/go-opencv/db_sql"

	"github.com/spf13/cobra"
)

func run(ctx context.Context, w io.Writer, args []string) error {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		return err
	}
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./db_sql/migrations",
		"sqlite3", driver)
	if err != nil {
		return err
	}

	m.Up()
	queries := db_sql.New(db)

	srv := &http.Server{Addr: ":3010", Handler: api.NewRouter(queries)}

	go func() {
		httpError := srv.ListenAndServe()
		if httpError != nil && httpError != http.ErrServerClosed {
			log.Fatal(httpError)
		}
	}()

	<-ctx.Done()

	gracefulTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	fmt.Printf("Shutting down server...\n")
	err = srv.Shutdown(gracefulTimeout)
	if err != nil {
		fmt.Printf("Error shutting down server: %v\n", err)
		return err
	}

	return nil
}

func NewRunserver() *cobra.Command {
	var cmdPrint = &cobra.Command{
		Use: "serve",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
			defer cancel()
			if err := run(ctx, os.Stdout, args); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			return nil
		},
	}

	return cmdPrint
}
