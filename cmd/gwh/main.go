package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/npclaudiu/gwh/v1"
	"github.com/urfave/cli/v2"
)

const (
	kAppName    = "gwh"
	kAppVersion = "0.1.0" // TODO(npclaudiu): Inject at build.
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
		Name:                   kAppName,
		Version:                kAppVersion,
		Usage:                  "manage local Git warehouses",
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
					warehouse, err := gwh.Open(location)

					if err != nil {
						die("init", err)
					}

					defer warehouse.Close()

					return nil
				},
			},
			{
				Name:  "link",
				Usage: "link repository",
				Args:  true,
				Action: func(cliCtx *cli.Context) error {
					log.Debug("linking repository...")

					location := cliPrefixFlag(cliCtx, cwd)
					warehouse, err := gwh.Open(location)

					if err != nil {
						die("link", err)
					}

					defer warehouse.Close()

					args := cliCtx.Args()
					name := args.Get(0)
					path := args.Get(1)

					if err := warehouse.LinkRepository(name, path); err != nil {
						die("link", err)
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

func die(cmd string, err error) {
	msg := fmt.Sprintf("%s failed", cmd)
	log.Fatal(msg, "error", err)
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
