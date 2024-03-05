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
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/db_sql"
	"github.com/vkhobor/go-opencv/importing"
	"github.com/vkhobor/go-opencv/jobs"
	"github.com/vkhobor/go-opencv/scraper"

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

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	queries := db_sql.New(db)
	jobCreator := &jobs.JobCreator{
		Queries: queries,
		Scrape: func(args jobs.ScrapeArgs, ctx context.Context) <-chan jobs.ScrapedVideo {
			return jobs.MapChannel(scraper.ScrapeToChannel(args.SearchQuery, ctx), func(id string) jobs.ScrapedVideo {
				return jobs.ScrapedVideo{ID: id}
			})
		},
		VImport: func(refs []string, vid ...jobs.DownlodedVideo) <-chan jobs.ImportedVideo {
			output := make(chan jobs.ImportedVideo)
			go func() {
				for _, video := range vid {
					val, err := importing.HandleVideoFromPath(video.SavePath, config.WorkDirImages, 1, "", refs)
					fmt.Printf("Imported video to %v\n", err)
					if err != nil {
						output <- jobs.ImportedVideo{
							Error: err,
						}
						continue
					}
					frames := make([]jobs.Frame, 0)
					for _, v := range val.FileNames {
						frames = append(frames, jobs.Frame{FrameNumber: 0, Path: v})
					}
					output <- jobs.ImportedVideo{
						DownlodedVideo:  video,
						ExtractedFrames: frames,
					}
				}
				close(output)
			}()
			return output
		},
		Download: func(vid ...jobs.ScrapedVideo) <-chan jobs.DownlodedVideo {
			output := make(chan jobs.DownlodedVideo)
			go func() {
				for _, video := range vid {
					time.Sleep(time.Second * 45)
					path, _, err := importing.DownloadVideo(video.ID)
					if err != nil {
						output <- jobs.DownlodedVideo{
							Error: err,
						}
					}
					output <- jobs.DownlodedVideo{
						ScrapedVideo: video,
						SavePath:     path,
					}
				}
				close(output)
			}()
			return output
		},
	}
	jobCreator.RunJobPool()

	srv := &http.Server{Addr: ":3010", Handler: api.NewRouter(queries, jobCreator)}

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
