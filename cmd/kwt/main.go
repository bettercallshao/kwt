package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/bettercallshao/kwt/pkg/menu"
	"github.com/bettercallshao/kwt/pkg/version"
)

func main() {
	log.SetPrefix("[kwt] ")

	app := &cli.App{
		Name:                 "kwt",
		Usage:                "Run commands easily.",
		Version:              version.Version,
		EnableBashCompletion: true,
		Commands: append(
			[]*cli.Command{
				{
					Name:    "start",
					Usage:   "Starts executor for a menu",
					Aliases: []string{"s"},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "master",
							Value:   "http://localhost:7171",
							Usage:   "Use a different master than default, use cautiously!",
							Aliases: []string{"r"},
						},
						&cli.StringFlag{
							Name:    "channel",
							Value:   "0",
							Usage:   "Select a channel.",
							Aliases: []string{"c"},
						},
						&cli.StringFlag{
							Name:     "menu",
							Usage:    "Menu to use",
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
					Usage:   "Ingests menu locally from a source",
					Aliases: []string{"i"},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "source",
							Usage:    "URL or file path to source",
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
				{
					// https://github.com/urfave/cli/blob/master/docs/v2/manual.md#default-auto-completion
					Name:  "complete",
					Usage: "Generate shell completion script",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "command",
							Value:   "kwt",
							Usage:   "The command to attach completion to",
							Aliases: []string{"c"},
						},
						&cli.StringFlag{
							Name:    "shell",
							Value:   "bash",
							Usage:   "The shell type",
							Aliases: []string{"s"},
						},
					},
					Action: func(c *cli.Context) error {
						script := `# bash completion script for kwt
_kwt_bash_autocomplete() {
	if [[ "${COMP_WORDS[0]}" != "source" ]]; then
	local cur opts base
	COMPREPLY=()
	cur="${COMP_WORDS[COMP_CWORD]}"
	if [[ "$cur" == "-"* ]]; then
		opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} ${cur} --generate-bash-completion )
	else
		opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
	fi
	COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
	return 0
	fi
}
complete -o bashdefault -o default -o nospace -F _kwt_bash_autocomplete ` + c.String("command")
						if c.String("shell") == "zsh" {
							script = `# zsh completion script for kwt
_CLI_ZSH_AUTOCOMPLETE_HACK=1
_kwt_zsh_autocomplete() {

	local -a opts
	local cur
	cur=${words[-1]}
	if [[ "$cur" == "-"* ]]; then
	opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[@]:0:#words[@]-1} ${cur} --generate-bash-completion)}")
	else
	opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[@]:0:#words[@]-1} --generate-bash-completion)}")
	fi

	if [[ "${opts[1]}" != "" ]]; then
	_describe 'values' opts
	else
	_files
	fi

	return
}
compdef _kwt_zsh_autocomplete ` + c.String("command")
						}
						fmt.Println(script)
						return nil
					},
				},
			},
			commands([]string{"s", "i"})...,
		),
	}

	app.Run(os.Args)
}
