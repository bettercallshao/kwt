package main

import (
	"encoding/json"
	"github.com/urfave/cli/v2"
	"net/url"
	"os"
	"os/signal"
	"path"
)

func join(base string, paths ...string) string {
	u, err := url.Parse(base)
	if err != nil {
		panic(err)
	}

	u.Path = path.Join(u.Path, path.Join(paths...))

	return u.String()
}

func wsURL(base string) string {
	u, err := url.Parse(base)
	if err != nil {
		panic(err)
	}

	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}

	return u.String()
}

func addParam(base string, key string, value string) string {
	u, err := url.Parse(base)
	if err != nil {
		panic(err)
	}

	q := u.Query()
	q.Add(key, value)
	u.RawQuery = q.Encode()

	return u.String()
}

func cancelChan() chan os.Signal {
	cancel := make(chan os.Signal, 1)
	signal.Notify(cancel, os.Interrupt)
	return cancel
}

func jsonParagraph(target interface{}) string {
	raw, _ := json.MarshalIndent(target, "", "  ")
	return string(raw)
}

func shorten(command *cli.Command) string {
	if len(command.Name) == 1 {
		return command.Name
	} else if len(command.Aliases) == 1 {
		return command.Aliases[0]
	} else {
		return ""
	}
}
