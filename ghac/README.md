# GitHub Access Control

Define GitHub teams and team memberships in yaml, render scripts and declarative configurations using templates.

## Usage

```
NAME:
   ghac - Manage GitHub Teams and Team membership in yaml

USAGE:
   ghac [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR:
   Honestbee DevOps

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-level value                      Log level (panic, fatal, error, warn, info, or debug) (default: "error") [$LOG_LEVEL]
   --source yaml, -s yaml                 Source yaml (default: "teams.yaml")
   --destination directory, -d directory  Destination directory to render in - must exist (default: "output/teams-config/")
   --template value, -t value             Desired template used to render output (default: "templates/team.tf.tpl")
   --filter regex, -f regex               regex filter on teams (Slug only for now)
   --help, -h                             show help
   --version, -v                          print the version
```

## Templating

Templates are go-templates and include Sprig function library

Read [Sprig guide for template writing](http://masterminds.github.io/sprig/)

## Why

Initial approach was to define teams as Terraform module which accepts a map or a list to generate team_membership resources.

As team_memberships change, terraform state would cascade when items of the list are removed.

This script is to generate TF Configs without lists, and to simplify managing these configurations by keeping the team
definition in a simple yaml data structure.
