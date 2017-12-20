# Honestbee Devops Tools

Misc Devops Tools used by Honestbee

## Golang Dependency management

Using glide, get all vendored dependencies with:

```
glide up
```

# Index

- [postgres](postgres/): script to batch backup / restore postgres databases
- [github-tf](github-tf/): script to read github teams and team_memberships and apply go templates (alpha)
- [ghac](ghac/): script to read specific `teams.yaml` data structure and apply go tempates (alpha)
- [helmns](helmns/): helm wrappers to handle a tiller per namespace deployments
- [tf-state-parser](tf-state-parser/): sample script on reading terraform state and sample to generate state refactoring commands
