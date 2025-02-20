package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/charmbracelet/log"
	gwh "github.com/npclaudiu/gwh/v1"
	"github.com/urfave/cli/v2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	cwd, err := os.Getwd()

	if err != nil {
		log.Fatal("failed to load working directory", err)
	}

	if !filepath.IsAbs(cwd) {
		log.Fatal("working directory must be absolute")
	}

	app := &cli.App{
		Name:                   "gwh",
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Suggest:                true,
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "initialize warehouse",
				Action: func(cliCtx *cli.Context) error {
					log.Debug("initializing warehouse...")

					// Handle "--cwd" option.
					//
					cwdOption := cliCtx.String("cwd")

					if cwdOption != "" {
						if !filepath.IsAbs(cwdOption) {
							cwdOption = filepath.Join(cwd, cwdOption)
						}
					} else {
						cwdOption = cwd
					}

					// Handle "--sync/--no-sync" option (TODO).
					//

					// Initialize the warehouse.
					//
					err := gwh.Init(ctx, &gwh.InitOptions{
						Cwd: cwdOption,
					})

					if err != nil {
						log.Fatal("init failed", "error", err)
					}

					return nil
				},
			},
			{
				Name:  "sync",
				Usage: "synchronize warehouse with repository",
				Action: func(cliCtx *cli.Context) error {
					log.Debug("synchronize warehouse with repository...")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
