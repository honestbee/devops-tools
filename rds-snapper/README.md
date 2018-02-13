# rds-snapper

A tool to manage rds snapshots.
This is designed to run during `drone pipeline` in which we could cleanup old snapshots and create a new one before deploying new code to production.

## Usage

### Commandline use case

```bash
NAME:
   rds-snapper - golang tools to manage RDS snapshots

USAGE:
   rds-snapper [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     export   Export snapshots list to csv file or stdout
     clear    Clear snapshot of specific dbname and keep only a specified number of snapshots
     create   Create new snapshot
     help, h  Shows a list of commands or help for one command
```

#### Export

```bash
NAME:
   rds-snapper export - Export snapshots list to csv file or stdout

USAGE:
   rds-snapper export [command options] [arguments...]

EXAMPLE:
   # print manual snapshots list to output.csv file.
   rds-snapper export --file "output.csv"
```

#### Clear

```bash
NAME:
   rds-snapper clear - Clear snapshot of specific dbname and only a specified keep number

USAGE:
   rds-snapper clear [command options] [arguments...]

EXAMPLE:
   # Clean up hb-staging rds's snapshots, only keep 5 latest copy.
   rds-snapper clear --dbname "hb-staging" --keep "5"
```

#### Create

```bash
NAME:
   rds-snapper create - Create new snapshot

USAGE:
   rds-snapper create [command options] [arguments...]

EXAMPLE:
   # Create new snapshot named "hb-staging-aaeec89" for hb-staging rds instance, the suffix here is github commit reference.
   # Option `--keep` could be provided to clean up old snapshots as well.
   rds-snapper create --dbname "hb-staging" --suffix "aaeec89" [--keep 5]
```

### Drone use case

```yaml
  # print list of snapshots to stdout
  export-snapshot:
    image: quay.io/honestbee/rds-snapper
    pull: true
    secrets:
      - source: snapshot_aws_access_key_id
        target: aws_access_key_id
      - source: snapshot_aws_secret_access_key
        target: aws_secret_access_key
    action: "export"

  # Clear <db-name>'s snapshots and keep only <keep> latest copies
  clear-snapshot:
    image: quay.io/honestbee/rds-snapper
    pull: true
    secrets:
      - source: snapshot_aws_access_key_id
        target: aws_access_key_id
      - source: snapshot_aws_secret_access_key
        target: aws_secret_access_key
    action: "clear"
    dbname: "<db-name>"
    keep: <numbers-to-keep>

  # Create new <db-name> snapshot
  create-snapshot:
    image: quay.io/honestbee/rds-snapper
    pull: true
    secrets:
      - source: snapshot_aws_access_key_id
        target: aws_access_key_id
      - source: snapshot_aws_secret_access_key
        target: aws_secret_access_key
    action: "create"
    dbname: "<db-name>"
    suffix: "<snapshot-name-suffix>"
    keep: <numbers-to-keep>
```
