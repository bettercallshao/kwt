package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

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
					Usage:   "starts executor for a <menu>",
					Aliases: []string{"s"},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "master",
							Value:   "http://localhost:7171",
							Usage:   "use a different master than default, use cautiously!",
							Aliases: []string{"m"},
						},
						&cli.StringFlag{
							Name:    "channel",
							Value:   "0",
							Usage:   "select a channel.",
							Aliases: []string{"c"},
						},
					},
					Action: func(c *cli.Context) error {
						if c.NArg() != 1 {
							return cli.Exit("error: expect exactly 1 argument", -1)
						}

						master := c.String("master")
						channel := c.String("channel")
						menuName := c.Args().Get(0)

						if start(master, channel, menuName) != nil {
							return cli.Exit("error: failed to start", -1)
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
