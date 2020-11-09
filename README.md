# DevCLi

## Building binary

First of all we need golang 1.14+ [install](https://golang.org)

```bash
# Run bash script
chmod +x ./build.sh
./build.sh

# Or buil with golang
go build -o dev main.go
```

## Usage

`dev help` Print help

```bash
Dev cli

Usage:
  dev [flags]
  dev [command]

Available Commands:
  help        Help about any command
  install     Install project
  list        List python projects
  update      Update projects
  version     Print the version

Flags:
  -h, --help   help for dev

Use "dev [command] --help" for more information about a command.
```

`dev install` Installing git repos to $HOME/.dev

```bash
# Single repo clone
dev install https://github.com/TheYkk/logger

#  Multiple repo clone
dev install https://github.com/TheYkk/logger https://github.com/TheYkk/synator
```

`dev list` Listing python scripts in $HOME/.dev

```bash
dev list

# ---SCRIPTS---
# zz => hi.py
# zz => main.py
# synator => handlers.py
```

`dev update` Updates git repos in $HOME/.dev

```bash
dev update

# Updating cobra
# already up-to-date cobra
# Updating devcli.git
# Updating gepp
# already up-to-date gepp
# Updating synator
# already up-to-date synator
```

`dev version` Print build version

```bash
dev version

# Version:  v20201022-2cadb1c
```