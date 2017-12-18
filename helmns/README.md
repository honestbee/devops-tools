This repository provides both `helmns` and `helmls` tools.

**`helmns`** `helm` cli wrapper forcing `--tiller-namespace` to match kubectl context

**`helmls`** a small script to list Helm releases across all `--tiller-namespaces`

# helmns

`helmns` is a [Helm](http://helm.sh) cli wrapper forcing `--tiller-namespace` to match kubectl context.

This is meant to simplify management for Kubernetes clusters with RBAC enabled and where tiller is heavily locked down per namespace.

# helmls

`helmls` is a small script to list Helm releases across all `--tiller-namespaces`.

-----

## Installation

Since `helmns` and `helmls` are written in Bash, they can run in shells that support POSIX standards.

- Download the `helmns` and `helmls` scripts
  
  ```bash
  for s in helmns helmls; do curl -sLo $s https://github.com/honestbee/devops-tools/tree/master/helmns/$s; chmod +x $s; done
  ```

- Either save all somewhere on your `PATH` (like `/usr/local/bin`):
  
  ```bash
  sudo mv helm* /usr/local/bin/
  ```

- Or move somewhere else and create symlinks yourself

-----

Heavily based on [ahmetb/kubens](https://github.com/ahmetb/kubectx)
