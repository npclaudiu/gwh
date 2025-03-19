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
		Description:            "Git Analitics Warehouse",
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Suggest:                true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "prefix",
				Aliases: []string{"p"},
				Value:   "",
				Usage:   "warehouse lookup directory (defaults to current working directory)",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "initialize warehouse",
				Action: func(cliCtx *cli.Context) error {
					log.Debug("initializing warehouse...")

					if err := gwh.Init(ctx, &gwh.InitOptions{
						Prefix: cliPrefixFlag(cliCtx, cwd),
					}); err != nil {
						log.Fatal("init failed", "error", err)
					}

					return nil
				},
			},
			{
				Name: "repository",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "add repository to warehouse",
						Action: func(cliCtx *cli.Context) error {
							log.Debug("adding repository...")

							if err := gwh.AddRepository(ctx, &gwh.AddRepositoryOptions{
								Prefix: cliPrefixFlag(cliCtx, cwd),
							}); err != nil {
								log.Fatal("adding repository failed", "error", err)
							}

							return nil
						},
					},
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

func cliFlag(cliCtx *cli.Context, flagName string, defaultValue string) string {
	if cliCtx.IsSet(flagName) {
		return cliCtx.String(flagName)
	}

	return defaultValue

}

func cliPrefixFlag(cliCtx *cli.Context, defaultValue string) string {
	prefix := cliFlag(cliCtx, "prefix", "")

	if prefix != "" {
		if !filepath.IsAbs(prefix) {
			prefix = filepath.Join(defaultValue, prefix)
		}
	} else {
		prefix = defaultValue
	}

	return prefix
}
