# kwt
[![Release](https://img.shields.io/github/release/bettercallshao/kwt.svg)](https://github.com/bettercallshao/kwt/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](/LICENSE.md)
[![CircleCI](https://circleci.com/gh/bettercallshao/kwt.svg?style=shield)](https://circleci.com/gh/bettercallshao/kwt)

Run commands easily.

* [What is it](#what-is-it)
* [Installation](#installation)
* [Quick start](#quick-start)
* [Web interface](#web-interface)
* [Develop](#develop)
* [Further info](#further-info)

## What is it

Kwt reads from a repository of YAMLs of templated commands, and enables the user to render and execute the commands with arguments provided at run-time. It can help to
* Avoid re-typing long commands over and over for routine tasks.
* Version control and share commands.

## Installation

Linux & Mac through [Homebrew](https://brew.sh/).
```
brew install bettercallshao/tap/kwt
```

Windows through [Scoop](https://scoop.sh/).
```
scoop bucket add bettercallshao https://github.com/bettercallshao/scoop-bucket
scoop install bettercallshao/kwt
```

Or download latest zip from [releases](https://github.com/bettercallshao/kwt/releases), extract, and put the binary files on your system path.

## Quick start

Kick start by ingesting the demo menus.
```
kwt i -s https://raw.githubusercontent.com/bettercallshao/kwt-menus/master/python-demo.yaml
kwt i -s https://raw.githubusercontent.com/bettercallshao/kwt-menus/master/developer-demo.yaml
```

See a list of commands by running `kwt h`.
```
NAME:
   kwt - Run commands easily.

USAGE:
   kwt [global options] command [command options] [arguments...]

VERSION:
   v0.5.2-20201123005537

COMMANDS:
   start, s           Starts executor for a menu
   ingest, i          Ingests menu locally from a source
   developer-demo, d  Developer commands for demo
   python-demo, p     Python commands for demo
   help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```
* `start`, `ingest`, `help`, `--help`, `--version` are global arguments.
* `python-demo` and `developer-demo` are command definitions.

Check what is in `python-demo` by `kwt p -h`.
```
NAME:
   kwt python-demo - Python commands for demo

USAGE:
   kwt python-demo command [command options] [arguments...]

COMMANDS:
   uuid, u                 Generate a UUID
   forex-rate, f           Print forex rates
   mortgage-calculator, m  Calculate mortgage payment
   bit-expander, b         Convert between decimal, hex, and bit representations
   csv-to-markdown, c      Convert a CSV to markdown table
   help, h                 Shows a list of commands or help for one command

OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

Check what `forex-rate` does with `kwt p f -h`.
```
NAME:
   kwt python-demo forex-rate - Print forex rates

USAGE:
   kwt python-demo forex-rate [command options] [arguments...]

OPTIONS:
   --base value, -b value     Base currency (default: "USD")
   --symbols value, -s value  Comma separated currency symbol list (default: "CAD,GBP")
   --dry, -d                  render command but don't run (default: false)
   --help, -h                 show help (default: false)
```
* `--dry`, `--help` are flags defined by `kwt`.
* `--base`, `--symbols` are input arguments, with default values.

Lets try EUR to JPY and AUD.
```
$ kwt p f -b EUR -s JPY,AUD
On 2020-11-13
1 EUR can buy 123.88 JPY
1 EUR can buy 1.63 AUD
```

Use the `--dry` flag to see what actually ran.
```
$ kwt p f -b EUR -s AUD -d
template: python3 -u -c "
import urllib.request
import urllib.parse
import json

url = 'https://api.exchangeratesapi.io/latest?base={{.base}}&symbols={{.symbols}}'
r = json.load(urllib.request.urlopen(url))

date = r['date']
rates = r['rates']

print(f'On {date}')
print('\n'.join(
  f'1 {{.base}} can buy {rates[symbol]:.2f} {symbol}'
  for symbol in rates
))
"

rendered: python3 -u -c "
import urllib.request
import urllib.parse
import json

url = 'https://api.exchangeratesapi.io/latest?base=EUR&symbols=AUD'
r = json.load(urllib.request.urlopen(url))

date = r['date']
rates = r['rates']

print(f'On {date}')
print('\n'.join(
  f'1 EUR can buy {rates[symbol]:.2f} {symbol}'
  for symbol in rates
))
"
```

First the command template was printed, then the command rendered with input arguments. The command is defined in `$HOME/.kwt/menus/python-demo.yaml`.
```yaml
name: python-demo
version: v0.1.0
help: Python commands for demo
actions:
- name: forex-rate
  help: Print forex rates
  template: |
    python3 -u -c "
    import urllib.request
    import urllib.parse
    import json

    url = 'https://api.exchangeratesapi.io/latest?base={{.base}}&symbols={{.symbols}}'
    r = json.load(urllib.request.urlopen(url))

    date = r['date']
    rates = r['rates']

    print(f'On {date}')
    print('\n'.join(
      f'1 {{.base}} can buy {rates[symbol]:.2f} {symbol}'
      for symbol in rates
    ))
    "
  params:
  - name: base
    help: Base currency
    value: USD
  - name: symbols
    help: Comma separated currency symbol list
    value: CAD,GBP
```

Add to this file or create more YAMLs in `$HOME/.kwt/menus/` to add more commands.

## Web interface

Kwt can also be run in conjunction with kwtd to give a web based user interface to the menus. Kwtd is installed as part of the kwt package and runs without arguments.
```
[kwtd] 2020/11/24 01:42:01 version: v0.5.2-20201123005537
[kwtd] 2020/11/24 01:42:01 starting kwtd ...
[kwtd] 2020/11/24 01:42:01 listening on http://127.0.0.1:7171
```

It is recommended to install kwtd as a start up service for convenience with the official helper menus.

For Windows (see help for more commands).
```
kwt i -s https://raw.githubusercontent.com/bettercallshao/kwt-menus/master/windows-kwtd.yaml
kwt windows-kwtd startup-add
```

For Mac (see help for more commands).
```
kwt i -s https://raw.githubusercontent.com/bettercallshao/kwt-menus/master/mac-kwtd.yaml
kwt mac-kwtd startup-add
```

Once kwtd is running, visit [http://127.0.0.1:7171](http://127.0.0.1:7171) in browser to find three sections.
* Channels - each channel is a placeholder for a kwt executor to connect to. If visited without active connection, it shows a blank message.
* Menus - each available menu can be viewed as a JSON.
* Ingestion - ingesting menus same as kwt.

As an example, we will run the `csv-to-markdown` command in the web interface. First open a terminal (with python3 available) and connect a kwt executor to kwtd (on channel 0 by default) declaring the `python-demo` menu.
```
kwt s -m python-demo
```

Logs are printed and the command should block and occupy the terminal. Now visit channel 0 on the page, click on `csv-to-markdown`, copy the following into the `data` param, press execute, then toggle markdown to render it.
```
Name,Icon,Website
Facebook,[![Website shields.io](https://img.shields.io/website-up-down-green-red/http/shields.io.svg)](),https://facebook.com/
Twitter,[![Website shields.io](https://img.shields.io/website-up-down-green-red/http/shields.io.svg)](),https://twitter.com/home
```

![Drag Racing](https://i.imgur.com/hQcheIU.gif)

## Develop

To build, install golang and run `make`. The CI is powered by [GoReleaser](https://goreleaser.com/) and CircleCI.

## Further info

I wrote a blog series on kwt https://bettercallshao.com/tags/kwt/

Please contact me via https://bettercallshao.com/author/
