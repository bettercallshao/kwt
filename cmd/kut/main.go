package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/bettercallshao/kut/pkg/menu"
	"github.com/bettercallshao/kut/pkg/version"
)

func main() {
	log.SetPrefix("[kut] ")

	app := &cli.App{
		Name:    "kut",
		Usage:   "run a kut executer.",
		Version: version.Version,
		Commands: append(
			[]*cli.Command{
				{
					Name:    "start",
					Usage:   "starts executor for a menu",
					Aliases: []string{"s"},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "master",
							Value:   "http://localhost:7171",
							Usage:   "use a different master than default, use cautiously!",
							Aliases: []string{"r"},
						},
						&cli.StringFlag{
							Name:    "channel",
							Value:   "0",
							Usage:   "select a channel.",
							Aliases: []string{"c"},
						},
						&cli.StringFlag{
							Name:     "menu",
							Usage:    "the menu to use",
							Aliases:  []string{"m"},
							Required: true,
						},
					},
					Action: func(c *cli.Context) error {
						master := c.String("master")
						channel := c.String("channel")
						menuName := c.String("menu")

						if start(master, channel, menuName) != nil {
							return cli.Exit("error: failed to start", -1)
						}

						return nil
					},
				},
				{
					Name:    "ingest",
					Usage:   "ingests menu locally from a source",
					Aliases: []string{"i"},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "source",
							Usage:    "source to ingest",
							Aliases:  []string{"s"},
							Required: true,
						},
					},
					Action: func(c *cli.Context) error {
						source := c.String("source")

						if menu.Ingest(source) != nil {
							return cli.Exit("error: failed to ingest", -1)
						}

						return nil
					},
				},
			},
			commands()...,
		),
	}

	app.Run(os.Args)
}
