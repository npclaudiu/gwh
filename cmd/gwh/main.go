package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/npclaudiu/gwh/v1"
	"github.com/urfave/cli/v2"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	cwd, err := os.Getwd()

	if err != nil {
		log.Fatal("failed to load working directory", err)
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
				Usage:   "directory containing the warehouse directory",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "initialize warehouse",
				Action: func(cliCtx *cli.Context) error {
					log.Debug("initializing warehouse...")

					location := cliPrefixFlag(cliCtx, cwd)
					_, err := gwh.Open(location)

					if err != nil {
						log.Fatal("init failed", "error", err)
					}

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

func cliPrefixFlag(cliCtx *cli.Context, cwd string) string {
	prefix := cliFlag(cliCtx, "prefix", "")

	if prefix != "" {
		if !filepath.IsAbs(prefix) {
			prefix = filepath.Join(cwd, prefix)
		}
	} else {
		prefix = cwd
	}

	return prefix
}
