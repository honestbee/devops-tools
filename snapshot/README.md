## rds-snapper

A tool to manage rds snapshots

## Usage

```bash
NAME:
   rds-snapper - golang tools to manage RDS snapshots

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     export   Export snapshots list to csv file
     clear    Clear snapshot of specific dbName and only a specified limit number
     create   Create new snapshot and name it with commit SHA
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Export

```bash
NAME:
   main export - Export snapshots list to csv file

USAGE:
   main export [command options] [arguments...]

OPTIONS:
   --aws-access-key AWS_ACCESS_KEY  AWS Access Key AWS_ACCESS_KEY [$AWS_ACCESS_KEY_ID, $AWS_ACCESS_KEY]
   --aws-secret-key AWS_SECRET_KEY  AWS Secret Key AWS_SECRET_KEY [$AWS_SECRET_ACCESS_KEY, $AWS_SECRET_KEY]
   --aws-region AWS_REGION          AWS Region AWS_REGION [$PLUGIN_AWS_REGION, $ AWS_REGION]
   --dbName value                   origin of snapshots [$PLUGIN_DBNAME]
   --file value                     file to save snapshots list [$PLUGIN_FILE]
```

### Clear

```bash
NAME:
   main clear - Keep specified number of snapshots and clean up the rest

USAGE:
   main clear [command options] [arguments...]

OPTIONS:
   --aws-access-key AWS_ACCESS_KEY  AWS Access Key AWS_ACCESS_KEY [$AWS_ACCESS_KEY_ID, $AWS_ACCESS_KEY]
   --aws-secret-key AWS_SECRET_KEY  AWS Secret Key AWS_SECRET_KEY [$AWS_SECRET_ACCESS_KEY, $AWS_SECRET_KEY]
   --aws-region AWS_REGION          AWS Region AWS_REGION [$PLUGIN_AWS_REGION, $ AWS_REGION]
   --dbName value                   origin of snapshots [$PLUGIN_DBNAME]
   --limit value                    number of snapshots to keep (default: 0) [$PLUGIN_LIMIT]
```

### Create

```bash
NAME:
   main create - Create new snapshot and name it with commit SHA

USAGE:
   main create [command options] [arguments...]

OPTIONS:
   --aws-access-key AWS_ACCESS_KEY  AWS Access Key AWS_ACCESS_KEY [$AWS_ACCESS_KEY_ID, $AWS_ACCESS_KEY]
   --aws-secret-key AWS_SECRET_KEY  AWS Secret Key AWS_SECRET_KEY [$AWS_SECRET_ACCESS_KEY, $AWS_SECRET_KEY]
   --aws-region AWS_REGION          AWS Region AWS_REGION [$PLUGIN_AWS_REGION, $ AWS_REGION]
   --dbName value                   origin of snapshots [$PLUGIN_DBNAME]
   --suffix value                   suffix to add to snapshot name (if not specified, would be a random string) [$PLUGIN_SUFFIX]
```