package main

import (
	"fmt"
	"strings"

	"github.com/bettercallshao/kwt/pkg/alias"
	"github.com/bettercallshao/kwt/pkg/cmd"
	"github.com/bettercallshao/kwt/pkg/menu"
	"github.com/urfave/cli/v2"
)

const DRY = "dry"

func commands(avoid []string) []*cli.Command {
	// Two aliases are reserved: help -> h, version -> v
	hv := []string{"h", "v"}
	commands := make([]*cli.Command, 0)
	store := alias.New()
	alias.Avoid(store, hv)
	alias.Avoid(store, avoid)

	dryFlag := &cli.BoolFlag{
		Name:    DRY,
		Usage:   "render command but don't run",
		Aliases: []string{"d"},
	}

	menuMap := menu.Map()
	// we need the ordering provided by List()
	for _, name := range menu.List() {
		menu := menuMap[name]

		subCommands := make([]*cli.Command, 0)
		subStore := alias.New()
		alias.Avoid(subStore, hv)

		for _, action := range menu.Actions {

			flags := make([]cli.Flag, 0)
			flagStore := alias.New()
			alias.Avoid(flagStore, hv)
			alias.Avoid(flagStore, dryFlag.Aliases)

			for _, param := range action.Params {

				flags = append(
					flags,
					&cli.StringFlag{
						Name:    param.Name,
						Value:   param.Value,
						Usage:   param.Help,
						Aliases: alias.Pick(flagStore, param.Name),
					},
				)
			}
			flags = append(flags, dryFlag)

			subCommands = append(
				subCommands,
				&cli.Command{
					Name:    action.Name,
					Usage:   action.Help,
					Flags:   flags,
					Action:  act,
					Aliases: alias.Pick(subStore, action.Name),
				},
			)
		}

		commands = append(
			commands,
			&cli.Command{
				Name:        menu.Name,
				Usage:       menu.Help + " " + menu.Path,
				Subcommands: subCommands,
				Aliases:     alias.Pick(store, menu.Name),
			},
		)
	}

	return commands
}

func act(c *cli.Context) error {
	dry := c.Bool(DRY)

	// only HelpName field under command can we find parent command
	name := strings.Split(c.Command.HelpName, " ")[1]

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
