## rds-snapper

A tool to manage rds snapshots.
This is designed to run during `drone pipeline` in which we could cleanup old snapshots and create a new one before deploying new code to production.

## Usage

```bash
NAME:
   rds-snapper - golang tools to manage RDS snapshots

USAGE:
   rds-snapper [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     export   Export snapshots list to csv file
     clear    Clear snapshot of specific dbName and only a specified limit number
     create   Create new snapshot and name it with commit SHA
     help, h  Shows a list of commands or help for one command
```

### Export

```bash
NAME:
   rds-snapper export - Export snapshots list to csv file

USAGE:
   rds-snapper export [command options] [arguments...]

EXAMPLE:
   # print manual snapshots list to output.csv file.
   rds-snapper export --file "output.csv"
```


### Clear

```bash
NAME:
   rds-snapper clear - Keep specified number of snapshots and clean up the rest

USAGE:
   rds-snapper clear [command options] [arguments...]

EXAMPLE:
   # Clean up hb-staging rds's snapshots, only keep 5 latest copy.
   rds-snapper clear --dbName "hb-staging" --limit "5"
```

### Create

```bash
NAME:
   rds-snapper create - Create new snapshot and name it with commit SHA

USAGE:
   rds-snapper create [command options] [arguments...]

EXAMPLE:
   # Create new snapshot named "hb-staging-aaeec89" for hb-staging rds instance, the suffix here is github commit reference.
   rds-snapper create --dbName "hb-staging" --suffix "aaeec89"
```