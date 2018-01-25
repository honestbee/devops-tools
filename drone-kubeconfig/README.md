# Drone KubeConfig

Set up Kube API and Token for a drone repo of choice

```
./drone-kfg.sh -h
USAGE:
  drone-kfg <PREFIX> <KUBE_CONTEXT> <REPO_NAME> : Add Drone secrets for Kubernetes for repo_name
  drone-kfg -h,--help                      : show this message
```

This will create the following secrets for `<repo_name>`:

- `<PREFIX>_API_SERVER`: The Kubernetes API server
- `<PREFIX>_KUBERNETES_TOKEN`: The token for the drone Service Account
