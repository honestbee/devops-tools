# Backstore - Mass backup & restore postgresql databases

## Introduction
It's pretty painful to manually dump & restore a large number of databases at once, so this tool is intent to help Operation team or DevOps team automate that task at a certain level.

## Features
- Dump & restore single database
- Dump & restore many databases via CSV input

## Usage
```bash
NAME:
   backstore - posgres database backup & restore

USAGE:
   backstore [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR:
   Tuan Nguyen

COMMANDS:
     dump, d     dump database
     restore, r  restore database
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dbname <database name>  <database name> to dump [$DBNAME]
   --dbhost <database host>  <database host> to connect to [$DBHOST]
   --dbuser <username>       <username> to authenticate with [$DBUSER]
   --dbpassword <password>   <password> to authenticate with (optional) [$DBPASSWORD]
   --config FILE             Load config from FILE [$DBCONFIG]
   --help, -h                show help
   --version, -v             print the version
```

## Notes
- This tool is still at basic level & created for learning purpose, so it might not have many advanced feature like [barkup](https://github.com/keighl/barkup), [go-sync](https://github.com/webdevops/go-sync) or [wal-g](https://github.com/wal-g/wal-g)
- PRs are welcome!
