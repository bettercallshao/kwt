package main

import (
	"fmt"
	"strings"

	"github.com/bettercallshao/kut/pkg/cmd"
	"github.com/bettercallshao/kut/pkg/menu"
	"github.com/urfave/cli/v2"
)

func commands() []*cli.Command {
	commands := make([]*cli.Command, 0)

	dryFlag := &cli.BoolFlag{
		Name:  "dry",
		Usage: "render command but don't run",
	}

	for _, name := range menu.List() {
		menu, err := menu.Load(name)
		if err != nil {
			continue
		}

		subCommands := make([]*cli.Command, 0)
		for _, action := range menu.Actions {

			flags := make([]cli.Flag, 0)
			for _, param := range action.Params {
				flags = append(
					flags,
					&cli.StringFlag{
						Name:  param.Name,
						Value: param.Value,
						Usage: param.Help,
					},
				)
			}
			flags = append(flags, dryFlag)

			subCommands = append(
				subCommands,
				&cli.Command{
					Name:   action.Name,
					Usage:  action.Help,
					Flags:  flags,
					Action: act,
				},
			)
		}

		commands = append(
			commands,
			&cli.Command{
				Name:        menu.Name,
				Usage:       menu.Help,
				Subcommands: subCommands,
			},
		)
	}

	return commands
}

func act(c *cli.Context) error {
	dry := c.Bool("dry")

	// the cli lib concatenates command hierarchy, menu name is the second
	name := strings.Split(c.App.Name, " ")[1]

	loaded, err := menu.Load(name)
	if err != nil {
		cli.Exit("failed to open menu "+name, -1)
	}

	action := menu.Action{}
	for _, item := range loaded.Actions {
		if item.Name == c.Command.Name {
			action = item
			break
		}
	}
	if action.Name == "" {
		cli.Exit("failed to find action "+c.Command.Name, -1)
	}

	for idx, param := range action.Params {
		action.Params[idx].Value = c.String(param.Name)
	}

	input, err := menu.Render(action)
	if err != nil {
		cli.Exit(fmt.Sprintf("action failed to render %+v", action), -1)
	}

	if dry {
		fmt.Printf("template: %s\n", action.Template)
		fmt.Printf("rendered: %s\n", input)
		return nil
	}

	cmd.Run(input, nil)

	return nil
}
