package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/urfave/cli/v2"
	"github.com/vkhobor/go-opencv/cmd"
	"github.com/vkhobor/go-opencv/config"
	"github.com/vkhobor/go-opencv/mlog"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "blob-storage",
					},
					&cli.IntFlag{
						Name: "port",
						Aliases: []string{
							"p",
						},
					},
					&cli.StringFlag{
						Name: "log-folder",
					},
					&cli.StringFlag{
						Name: "db",
					},
				},
				Usage: "starts the server",
				Action: func(cCtx *cli.Context) error {
					var k = koanf.New(".")

					home, err := os.UserHomeDir()
					if err != nil {
						return err
					}

					// Set up default configuration
					k.Load(structs.Provider(config.ServerConfig{
						Port:        7000,
						BlobStorage: home + "/go-extractor/blob",
						BaseUrl:     "http://localhost:7000",
						LogFolder:   home + "/go-extractor/log",
						Db:          home + "/go-extractor/db.sqlite3",
					}, "koanf"), nil)

					k.Load(env.ProviderWithValue("GO_EXTRACTOR", ".", func(s string, v string) (string, any) {
						// Strip out the MYVAR_ prefix and lowercase and get the key while also replacing
						// the _ character with . in the key (koanf delimeter).
						key := strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(s, "GO_EXTRACTOR")), "_", ".")

						// If there is a space in the value, split the value into a slice by the space.
						if strings.Contains(v, " ") {
							return key, strings.Split(v, " ")
						}

						// Otherwise, return the plain string.
						return key, v
					}), nil)

					for _, name := range cCtx.FlagNames() {
						if cCtx.IsSet(name) {
							val := cCtx.Value(name)
							normalized := strings.ReplaceAll(name, "-", "")
							k.Set(normalized, val)
						}
					}

					var config config.ServerConfig
					k.Unmarshal("", &config)

					args := cCtx.Args().Slice()
					ctx := cCtx.Context

					ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)
					defer cancel()

					mlog.Log().Info("Starting server", "configuration", config)
					if err := cmd.RunServer(ctx, os.Stdout, args, config); err != nil {
						fmt.Fprintf(os.Stderr, "%s\n", err)
						return err
					}
					return nil
				},
			},
			{
				Name: "openapi",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
					},
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
					},
				},
				Action: func(cCtx *cli.Context) error {
					args := cCtx.Args().Slice()
					ctx := cCtx.Context

					ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)
					defer cancel()

					port := cCtx.Int("port")
					output := cCtx.String("output")

					mlog.Log().Info("Starting openapi cmd", "port", port, "output", output)

					if err := cmd.RunOpenApiServer(ctx, os.Stdout, args, port, output); err != nil {
						fmt.Fprintf(os.Stderr, "%s\n", err)
						return err
					}

					return nil
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Action: func(cCtx *cli.Context) error {
					fmt.Println("version 1.0.0")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
