# Drone KubeConfig

> Retrieve and create Drone secrets for Kubernetes service accounts

Quickly configure Drone to deploy across multiple Kubernetes clusters. 

This script prepares repositories for Kubernetes deployments using the [drone-helm plugin](github.com/ipedrazas/drone-helm/) - which expects  `API_SERVER` and `KUBERNETES_TOKEN` as secrets.

## Pre-reqs

To keep the tool simple, it expects `kubectl` and `drone` binaries to exist in the `$PATH` and be fully configured.

Tested with Kubernetes 1.7.12 and Drone 0.8.0

## Golang binary

```
NAME:
   drone-kfg - retrieve and create drone secrets for kubernetes service accounts

USAGE:
   drone-kfg [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --repository value, -r value                           repository name (e.g. octocat/hello-world)
   --context PREFIX=KUBE_CONTEXT, -c PREFIX=KUBE_CONTEXT  PREFIX=KUBE_CONTEXT pairs for retrieving drone secrets
   --service-account value, -s value                      Kubernetes service account for drone (default: "drone")
   --namespace value, -n value                            Kubernetes namespace of service account (default: "kube-system")
   --timeout DURATION, -t DURATION                        DURATION before commands are cancelled (default: 5s)
   --help, -h                                             show help
   --version, -v                                          print the version
```

Example usage:

```bash
bin/drone-kfg -c PRODUCTION_1A=1a.prod.mycluster.com -c PRODUCTION_1B=1b.prod.mycluster.com octocat/hello-world
```

This will create the following secrets for `octocat/hello-world`:

- `PRODUCTION_1A_API_SERVER`: The API server for the 1A kubernetes production cluster
- `PRODUCTION_1A_KUBERNETES_TOKEN`: The token for the drone Service Account in the 1A kubernetes production cluster
- `PRODUCTION_1B_API_SERVER`: The API server for the 1B kubernetes production cluster
- `PRODUCTION_1B_KUBERNETES_TOKEN`: The token for the drone Service Account in the 1B kubernetes production cluster

## Bash script (deprecated)

Early iteration of drone configuration script - Set up Kube API and Token for a drone repo of choice

```
./drone-kfg.sh -h
USAGE:
  drone-kfg <PREFIX> <KUBE_CONTEXT> <REPO_NAME> : Add Drone secrets for Kubernetes for repo_name
  drone-kfg -h,--help                      : show this message
```

This will create the following secrets for `<repo_name>`:

- `<PREFIX>_API_SERVER`: The Kubernetes API server
- `<PREFIX>_KUBERNETES_TOKEN`: The token for the drone Service Account
