# Drone KubeConfig

> Retrieve and create Drone secrets for Kubernetes service accounts

Quickly configure Drone to deploy across multiple Kubernetes clusters.

This script prepares repositories for Kubernetes deployments using the [drone-helm plugin](github.com/ipedrazas/drone-helm/) - which expects  `API_SERVER` and `KUBERNETES_TOKEN` as secrets.

## Pre-reqs

To keep the tool simple, it expects `kubectl` and `drone` binaries to exist in the `$PATH` and be fully configured.

Tested with Kubernetes 1.7.15 and Drone 0.8.0

## Golang binary

```bash
NAME:
   drone-kubeconfig - create drone secrets for kubernetes service accounts

USAGE:
   drone-kubeconfig [global options] command [command options] [arguments...]

VERSION:
   0.1.1

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --repository value, -r value       repository name (e.g. octocat/hello-world)
   --context PREFIX, -c PREFIX        map PREFIX with equivalent `KUBE_CONTEXT` (e.g STAGING_1A=<kube-content-name>)
   --service-account value, -s value  Kubernetes service account for drone (default: "drone")
   --namespace value, -n value        Kubernetes namespace of service account (default: "kube-system")
   --timeout DURATION, -t DURATION    DURATION before commands are cancelled (default: 5s)
   --help, -h                         show help
   --version, -v                      print the version

```

Example usage:

```bash
bin/drone-kubeconfig -c PRODUCTION_1A=1a.prod.mycluster.com -c PRODUCTION_1B=1b.prod.mycluster.com octocat/hello-world
```

This will create the following secrets for `octocat/hello-world`:

- `PRODUCTION_1A_API_SERVER`: The API server for the 1A kubernetes production cluster
- `PRODUCTION_1A_KUBERNETES_TOKEN`: The token for the drone Service Account in the 1A kubernetes production cluster
- `PRODUCTION_1B_API_SERVER`: The API server for the 1B kubernetes production cluster
- `PRODUCTION_1B_KUBERNETES_TOKEN`: The token for the drone Service Account in the 1B kubernetes production cluster
