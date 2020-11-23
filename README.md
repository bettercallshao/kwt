# kut
[![Release](https://img.shields.io/github/release/bettercallshao/kut.svg)](https://github.com/bettercallshao/kut/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](/LICENSE.md)
[![CircleCI](https://circleci.com/gh/bettercallshao/kut.svg?style=shield)](https://circleci.com/gh/bettercallshao/kut)

Run commands easily.

* [What is it](#what-is-it)
* [Installation](#installation)
* [Quick start](#quick-start)
* [Develop](#develop)
* [TODO](#todo)

## What is it

Kut reads from a repository of YAMLs of templated commands, and enables the user to render and execute the commands with arguments provided at run-time. It can help to
* Avoid re-typing long commands over and over for routine tasks.
* Version control and share commands.

## Installation

Linux & Mac through [Homebrew](https://brew.sh/).
```
brew install bettercallshao/tap/kut
```

Windows through [Scoop](https://scoop.sh/).
```
scoop bucket add bettercallshao https://github.com/bettercallshao/scoop-bucket
scoop install bettercallshao/kut
```

Or download latest zip from [releases](https://github.com/bettercallshao/kut/releases), extract, and put the binary files on your system path.

## Quick start

Kick start by ingesting the demo menus.
```
kut i -s https://raw.githubusercontent.com/bettercallshao/kut-menus/master/python-demo.yaml
kut i -s https://raw.githubusercontent.com/bettercallshao/kut-menus/master/developer-demo.yaml
```

See a list of commands by running `kut h`.
```
NAME:
   kut - run a kut executer.

USAGE:
   kut [global options] command [command options] [arguments...]

VERSION:
   v0.4.11-20201115140037

COMMANDS:
   start, s           starts executor for a menu
   ingest, i          ingests menu locally from a source
   developer-demo, d  Developer commands for demo
   python-demo, p     Python commands for demo
   help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```
* `start`, `ingest`, `help`, `--help`, `--version` are global arguments.
* `python-demo` and `developer-demo` are command definitions.

Check what is in `python-demo` by `kut p -h`.
```
NAME:
   kut python-demo - Python commands for demo

USAGE:
   kut python-demo command [command options] [arguments...]

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

Check what `forex-rate` does with `kut p f -h`.
```
NAME:
   kut python-demo forex-rate - Print forex rates

USAGE:
   kut python-demo forex-rate [command options] [arguments...]

OPTIONS:
   --base value, -b value     Base currency (default: "USD")
   --symbols value, -s value  Comma separated currency symbol list (default: "CAD,GBP")
   --dry, -d                  render command but don't run (default: false)
   --help, -h                 show help (default: false)
```
* `--dry`, `--help` are flags defined by `kut`.
* `--base`, `--symbols` are input arguments, with default values.

Lets try EUR to JPY and AUD.
```
$ kut p f -b EUR -s JPY,AUD
On 2020-11-13
1 EUR can buy 123.88 JPY
1 EUR can buy 1.63 AUD
```

Use the `--dry` flag to see what actually ran.
```
$ kut p f -b EUR -s AUD -d
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

First the command template was printed, then the command rendered with input arguments. The command is defined in `$HOME/.kut/menus/python-demo.yaml`.
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

Add to this file or create more YAMLs in `$HOME/.kut/menus/` to add more commands.

## Develop

To build, install golang and run `make`. The CI is powered by [GoReleaser](https://goreleaser.com/) and CircleCI.

## TODO

This project builds `kut` and `kutd`. `kutd` is a web frontend allowing access to the same commands from the browser. Documentation is to be completed, for the meantime, there is an outdated [blog](https://medium.com/@bettercallshao/kut-free-ui-for-everyone-a262a82c5bab).