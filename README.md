# Honestbee Devops Tools

Misc Devops Tools used by Honestbee

## Golang Dependency management

Using dep, get all vendored dependencies with:

```bash
dep ensure
```

## Index

- [postgres](postgres/): script to batch backup / restore postgres databases
- [github-tf](github-tf/): script to read github teams and team_memberships and apply go templates (alpha)
- [ghac](ghac/): script to read specific `teams.yaml` data structure and apply go tempates (alpha)
- [helmns](helmns/): helm wrappers to handle a tiller per namespace deployments
- [tf-state-parser](tf-state-parser/): sample script on reading terraform state and sample to generate state refactoring commands

## Release howto

### Introduction

Basically, we use tools such as [mbt](https://github.com/mbtproject/mbt) and [gox](https://github.com/mitchellh/gox), which are combined into one single [docker image](https://quay.io/repository/honestbee/monorepo) to create binaries and upload to github release automatically via drone pipeline.

### Step to release

- Work on your PR, change your tool version appropriately. See [sample](ghac/.mbt.yml)

- Get it merged to master

- Tag your master release with the following format `v<YEAR>.<MONTH>.<RELEASE-NUMBER>`. For e.g: `v2018.02.0` (released on Feb 2018, first release of this month)

```bash
# Ensure that you are on up-to-date master branch
git tag -a v2018.02.0 -m "Release on Feb 2018"
git push origin v2018.02.0
```

- Screenshot:

![](https://i.imgur.com/ogxkjtI.jpg)
