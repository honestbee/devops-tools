# GitHub TF

Golang script to extract an organisation into TF Configs.

Restrictions:

- This generates TF configs which depend on a private GitHub team module (which is trivial and may be open sourced upon request)
- This currently only extracts teams and team membership (not repositories)

Example usage:

```bash
export GITHUB_ORGANIZATION=honestbee
export GITHUB_TOKEN=<token>

# get debug output (keeps track of github rate limit)
export LOG_LEVEL=debug

# get team membership filtered by slug names team1 and team2
bin/github-tf -f team1 -f team2
```

```bash
bin/github-tf --help
NAME:
   github-tf - Download GitHub teams to TF config

USAGE:
   github-tf [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR:
   Honestbee DevOps

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --org organization, -o organization  organization to generate tf config for [$GITHUB_ORGANIZATION]
   --token token, -t token              token to access GitHub API [$GITHUB_TOKEN]
   --log-level value                    Log level (panic, fatal, error, warn, info, or debug) (default: "error") [$LOG_LEVEL]
   --team-filter slugs, -f slugs        slugs to filter teams by
   --help, -h                           show help
   --version, -v                        print the version
```