# kut

A GUI to help make running commands easier.

## Background

### Motivation

This project aims to help in the following scenarios

* You are a developer and you have long commands that you repeated type, and you need a easier way to compose the commands.

* You are a developer and you are often asked to help analysts fix their software because it requires running a few commands and the analysts are not comfortable dealing with terminals, but of course the command is so simple it's not worth making a UI for it.

### Approach

* Use Go instead of Node or Python to allow easy installation on Windows.

* Pacakge all frontend resources in binary to allow minimal setup.

* Group each functional group as a "Menu" which is a defined by a YAML file, e.g.

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

## Usage

### Installation

* Download a distributable zip from Github releases.

* Unzip the package and move the executables to a place where they can be executed.

### Ingesting Menu

* Run `kutd` or `kutd.exe` or create a shortcut to run. This is the web server and message exchange server and it needs to keep running.

* Visit http://localhost:7171 .

* At bottom of the page, locate the Ingest Menu form, enter

    ```
    name: windows-shell
    source: https://raw.githubusercontent.com/bettercallshao/kut-menus/master/windows-shell.yaml
    ```

    or

    ```
    name: linux-shell
    source: https://raw.githubusercontent.com/bettercallshao/kut-menus/master/linux-shell.yaml
    ```

    Note source could be either a http link or a local file.

* Press ingest to observe a new Menu item.

### Starting Channel

A predefined number of Channels are created in the web server (Master), a command runner (Executor) can connect to a channel to create an execution environment to run commands.

* Run `kut` or `kut.exe` in the desired working directory with arguments to claim a Channel with a Menu name and default Channel 0, e.g.

    ```
    kut start linux-shell
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

## Building

* Install Go.

* Install go-assets-builder.
    ```
    go install go-assets-builder
    ```

* Get source code.
    ```
    go get github.com/bettercallshao/kut
    ```

* Locate the directory in Go's home and cd into it.

* __[Windows]__ Install cygwin. cygwin is the only `make` environment tested, others may work as well.

* Make distributable package.
    ```
    PLATFORM=windows-x86_64 make package
    ```

* Locate the zip files in `dist`.
