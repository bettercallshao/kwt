package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/bettercallshao/kwt/pkg/alias"
	"github.com/bettercallshao/kwt/pkg/menu"
	"github.com/bettercallshao/kwt/pkg/version"
)

func main() {
	log.SetPrefix("[kwt] ")

	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} [global options] command subcommand [options]{{if .Commands}}

VERSION:
   {{.Version}}

COMMANDS:
{{range .Commands}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}
{{if .Subcommands}}{{range .Subcommands}}{{"\t"}}{{join .Names ", "}}{{"\t"}}{{.Usage}}
{{end}}{{end}}{{end}}{{end}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}
`

	app := &cli.App{
		Name:                 "kwt",
		Usage:                "Run commands easily.",
		Version:              version.Version,
		EnableBashCompletion: true,
		Commands: append(
			[]*cli.Command{
				{
					Name:  "start",
					Usage: "Starts executor for a menu",
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
					Name:  "ingest",
					Usage: "Ingests menu locally from a source",
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
    local cur opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"

    local cmdline
    cmdline=$(printf "%q " "${COMP_WORDS[@]:0:$COMP_CWORD}")

    if [[ "$cur" == "-"* ]]; then
        opts=$(eval "${cmdline} ${cur} --generate-bash-completion")
    else
        opts=$(eval "${cmdline} --generate-bash-completion")
    fi

    # Replace cut the line short to the word before ':'
    opts=$(echo "${opts}" | sed 's/^\([^:]*\):.*/\1/')

    COMPREPLY=( $(compgen -W "${opts}" -S '' -- "${cur}") )
    return 0
}

complete -o bashdefault -o default -o nospace -F _kwt_bash_autocomplete ` + c.String("command")
						if c.String("shell") == "zsh" {
							script = `# zsh completion script for kwt
_kwt_zsh_autocomplete() {
    local -a opts
    local curcontext="$curcontext" state line
    typeset -A opt_args

    cur="${words[CURRENT]}"
    if [[ "$cur" == "-"* ]]; then
        opts=("${(@f)$( ${words[1,CURRENT-1]} ${cur} --generate-bash-completion )}")
    else
        opts=("${(@f)$( ${words[1,CURRENT-1]} --generate-bash-completion )}")
    fi

    _describe 'kwt options' opts
    return
}

compdef _kwt_zsh_autocomplete kwt`
						}
						fmt.Println(script)
						return nil
					},
				},
				{
					Name:  "shorthands",
					Usage: "Generate no space short hands as shell alias",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "prefix",
							Value:   "kk",
							Usage:   "The string to prefix short hands",
							Aliases: []string{"p"},
						},
					},
					Action: func(c *cli.Context) error {
						prefix := c.String("prefix")
						for _, command := range c.App.Commands {
							first := shorten(command)
							if first == "" {
								continue
							}
							fmt.Printf("alias %s%s='kwt %s'\n", prefix, first, first)
							for _, subcommand := range command.Subcommands {
								second := shorten(subcommand)
								if second == "" {
									continue
								}
								fmt.Printf("alias %s%s%s='kwt %s %s'\n",
									prefix, first, second, first, second)
								for _, flag := range subcommand.Flags {
									if stringFlag, ok := flag.(*cli.StringFlag); ok {
										if len(stringFlag.Aliases) == 1 {
											third := stringFlag.Aliases[0]
											fmt.Printf("alias %s%s%s%s='kwt %s %s -%s'\n",
												prefix, first, second, third, first, second, third)
										}
									}
								}
							}
						}
						return nil
					},
				},
			},
			commands()...,
		),
	}

	store := alias.New()
	for _, command := range app.Commands {
		alias.Avoid(store, command.Aliases)
	}
	for _, command := range app.Commands {
		if len(command.Aliases) == 0 {
			command.Aliases = alias.Pick(store, command.Name)
		}
		subStore := alias.New()
		for _, subcommand := range command.Subcommands {
			alias.Avoid(subStore, subcommand.Aliases)
		}
		for _, subcommand := range command.Subcommands {
			if len(subcommand.Aliases) == 0 {
				subcommand.Aliases = alias.Pick(subStore, subcommand.Name)
			}
			flagStore := alias.New()
			for _, flag := range subcommand.Flags {
				if boolFlag, ok := flag.(*cli.BoolFlag); ok {
					alias.Avoid(flagStore, boolFlag.Aliases)
				} else if stringFlag, ok := flag.(*cli.StringFlag); ok {
					alias.Avoid(flagStore, stringFlag.Aliases)
				}
			}
			for _, flag := range subcommand.Flags {
				if boolFlag, ok := flag.(*cli.BoolFlag); ok {
					if len(boolFlag.Aliases) == 0 {
						boolFlag.Aliases = alias.Pick(flagStore, boolFlag.Name)
					}
				} else if stringFlag, ok := flag.(*cli.StringFlag); ok {
					if len(stringFlag.Aliases) == 0 {
						stringFlag.Aliases = alias.Pick(flagStore, stringFlag.Name)
					}
				}
			}
		}
	}

	app.Run(os.Args)
}
