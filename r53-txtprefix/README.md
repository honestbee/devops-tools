# r53 txtprefix

Prefix all r53 of type TXT with prefix

## Building

Use make targets

## Usage

**Note**: Make targets automatically load `.env` file into the enivronment

```bash
NAME:
   r53-txtprefix - Prefix all r53 of type TXT with provided string

USAGE:
   prefix [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR:
   Honestbee DevOps

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --hosted-zone-id id         Hosted zone id for route53 [$AWS_HOSTED_ZONE_ID]
   --prefix prefix             prefix for external-dns TXT records (default: "prefix.")
   --log-level value           Log level (panic, fatal, error, warn, info, or debug) (default: "error") [$LOG_LEVEL]
   --region region             default region to use (default: "ap-southeast-1") [$AWS_REGION]
   --access-key-id key         default key to use
   --secret-access-key secret  default secret to use
   --help, -h                  show help
   --version, -v               print the version
```

## Result

```bash
INFO[0001] Skipping (already prefixed) "prefix.foo.honestbee.com."
INFO[0002] Prefixing "bar.honestbee.com." ...
```

```json
"Change" {
  "ChangeBatch": {
    "Changes": [{
        "Action": "DELETE",
        "ResourceRecordSet": {
          "Name": "bar.honestbee.com.",
          "ResourceRecords": [{
              "Value": "\"heritage=external-dns,external-dns/owner=ap-southeast-1a.honestbee.com\""
            }],
          "TTL": 300,
          "Type": "TXT"
        }
      },{
        "Action": "CREATE",
        "ResourceRecordSet": {
          "Name": "prefix.bar.honestbee.com.",
          "ResourceRecords": [{
              "Value": "\"heritage=external-dns,external-dns/owner=ap-southeast-1a.honestbee.com\""
            }],
          "TTL": 300,
          "Type": "TXT"
        }
      }],
    "Comment": "Update prefix"
  },
  "HostedZoneId": "..."
}
```
