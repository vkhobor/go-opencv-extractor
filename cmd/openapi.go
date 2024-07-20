package cmd

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vkhobor/go-opencv/api"
	"github.com/vkhobor/go-opencv/config"
	database "github.com/vkhobor/go-opencv/db"
	"github.com/vkhobor/go-opencv/mlog"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/spf13/cobra"
)

func runOpenApiServer(ctx context.Context, _ io.Writer, _ []string, port int, output string) error {
	db, _, err := sqlmock.New()
	if err != nil {
		return err
	}
	defer db.Close()

	dbQueries := database.New(db)

	dirConfig, err := config.NewDirectoryConfig("./")
	if err != nil {
		return err
	}

	emptyChan := make(chan struct{})
	portString := fmt.Sprintf(":%d", port)

	mlog.Log().Info("Starting router", "port", port)
	router := api.NewRouter(dbQueries, emptyChan, dirConfig, config.ProgramConfig{})

	savingOpenApiDone := make(chan struct{})

	go func() {
		mlog.Log().Info("Server starting", "port", port)
		l, err := net.Listen("tcp", portString)
		if err != nil {
			log.Fatal(err)
		}

		wg := &sync.WaitGroup{}
		wg.Add(2)
		go mustSaveOpenApiSpecs(wg, port, output, "/openapi-3.0.json")
		go mustSaveOpenApiSpecs(wg, port, output, "/openapi.json")

		go func() {
			wg.Wait()
			close(savingOpenApiDone)
		}()

		if err := http.Serve(l, router); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	select {
	case <-ctx.Done():
	case <-savingOpenApiDone:
	}
	return nil
}

func mustSaveOpenApiSpecs(wg *sync.WaitGroup, port int, output string, apiPath string) {
	if err := saveOpenApiSpecs(wg, port, output, apiPath); err != nil {
		log.Fatal(err)
	}
}

func saveOpenApiSpecs(wg *sync.WaitGroup, port int, output string, apiPath string) error {
	defer wg.Done()

	url := fmt.Sprintf("http://localhost:%d%v", port, apiPath)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making request to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-200 status code: %d", resp.StatusCode)
	}

	outputFile := fmt.Sprintf("%s%v", output, apiPath)
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("error saving specs to file: %w", err)
	}

	slog.Info("Saved openapi spec", "file", outputFile)

	return nil
}

func NewRunOpenApi() *cobra.Command {

	var cmdPrint = &cobra.Command{
		Use: "openapi",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)
			defer cancel()

			port, err := cmd.Flags().GetInt("port")
			if err != nil {
				return err
			}

			output, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}

			mlog.Log().Info("Starting openapi cmd", "port", port, "output", output)

			if err := runOpenApiServer(ctx, os.Stdout, args, port, output); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			return nil
		},
	}

	cmdPrint.Flags().IntP("port", "p", 8080, "Specify the port")
	cmdPrint.MarkFlagRequired("port")

	cmdPrint.Flags().StringP("output", "o", "./doc", "Specify the folder to output the openapi specs")
	cmdPrint.MarkFlagRequired("output")

	return cmdPrint
}
