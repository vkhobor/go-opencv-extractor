package cmd

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"sync"

	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vkhobor/go-opencv/api"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/mlog"

	"github.com/DATA-DOG/go-sqlmock"
)

func RunOpenApiServer(ctx context.Context, _ io.Writer, _ []string, port int, output string) error {
	db, _, err := sqlmock.New()
	if err != nil {
		return err
	}
	defer db.Close()

	conf := config.ServerConfig{
		BlobStorage: "/",
	}
	dirConfig, err := conf.GetDirectoryConfig()
	if err != nil {
		return err
	}

	emptyChan := make(chan struct{})
	portString := fmt.Sprintf(":%d", port)

	mlog.Log().Info("Starting router", "port", port)
	router := api.NewRouter(db, emptyChan, dirConfig, config.ServerConfig{})

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
