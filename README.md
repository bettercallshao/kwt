# kut
[![CircleCI](https://circleci.com/gh/bettercallshao/kut.svg?style=shield)](https://circleci.com/gh/bettercallshao/kut)

A GUI to help make running commands easier.

## Background

### Motivation

This project aims to help in the following scenarios

* You are a developer and you have long commands that you repeated type, and you need a easier way to compose the commands.

* You are a developer and you are often asked to help analysts fix their computer because it requires running a few commands and the analysts are not comfortable dealing with terminals, but of course the command is so simple it's not worth making a UI for it.

### Approach

* Use Go instead of Node or Python to allow easy installation on Windows.

* Pacakge all frontend resources in binary to allow minimal setup.

* Group each functional group as a "Menu" which is a defined by a YAML file saved in `~/.kut/menus`, e.g.

```yaml
name: windows-shell
help: Collection of common windows shell commands
actions:
- name: ping
  help: Ping a party in network
  template: ping {{.host}}
  params:
  - name: host
    help: Host to ping
    value: google.com
- name: ip
  help: Show my ip adresses
  template: ipconfig
```

* The user executes an Action by providing Params, which is rendered into the Template (Go Template) and executed in a shell environment.

* The user can view the results of the execution in the frontend.

## Web Mode Usage

### Installation

* Download a distributable zip from Github releases.

* Unzip the package and move the executables to a place where they can be executed.

### Ingesting Menu

* Run `kutd` or `kutd.exe` or create a shortcut to run. This is the web server and message exchange server and it needs to keep running.

* Visit http://localhost:7171 .

* At bottom of the page, locate the Ingest Menu form, enter into source

    ```
    https://raw.githubusercontent.com/bettercallshao/kut-menus/master/windows-shell.yaml
    ```

    or

    ```
    https://raw.githubusercontent.com/bettercallshao/kut-menus/master/linux-shell.yaml
    ```

    Note source could be either a http link or a local file.

* The official menus are found in https://github.com/bettercallshao/kut-menus .

* Press ingest to observe a new Menu item.

### Starting Channel

A predefined number of Channels are created in the web server (Master), a command runner (Executor) can connect to a channel to create an execution environment to run commands.

* Run `kut` or `kut.exe` in the desired working directory with arguments to claim a Channel with a Menu name and default Channel 0, e.g.

    ```
    kut start --menu linux-shell
    ```

* On Windows, this is best done by creating a shortcut to `kutd.exe` on the Desktop, then editing the properties to specify a desired starting directory and arguments. Name the shortcut properly and use a lot of them.

* `kut` is built with `github.com/urfave/cli` and help is available with

    ```
    kut -h
    ```

* This needs to keep running for Actions to execute.

### Executing Action

* Visit http://localhost:7171 and select Channel 0.

* Observe the list of Actions, select one.

* Click Execute and observe output at the bottom.

## Standalone Usage

The `kut` client can be run standalone to render and execute commands based on commandline input.

* After installation, menu can be ingested locally for standalone consumption by using the `ingest` command, e.g.

    ```
    kut ingest --source https://raw.githubusercontent.com/bettercallshao/kut-menus/master/linux-shell.yaml
    ```

* Actions in menus are listed in the help page of `kut`.
    ```
    kut -h
    ```

* Run a action by using the menu name as the primary command, action name as sub command, and parameters as flags, e.g.
    ```
    kut linux-shell ping --host twitter.com
    ```

## Building

* Install Go.

* Install go-assets-builder.
    ```
    go install github.com/jessevdk/go-assets-builder
    ```

* Get source code. This uses `go mod` so it lives out of the go path.
    ```
    git clone git@github.com:bettercallshao/kut.git
    ```

* __[Windows]__ Install cygwin. cygwin is the only `make` environment tested, others may work as well.

* Make distributable package.
    ```
    GOOS=windows GOARCH=amd64 make clean package
    ```

* Locate the zip files in `dist`.
