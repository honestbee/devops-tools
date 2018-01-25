#!/bin/bash

set -eou pipefail

usage() {
  cat <<"EOF"
USAGE:
  drone-kfg <PREFIX> <KUBE_CONTEXT> <REPO_NAME> : Add Drone secrets for Kubernetes for repo_name
  drone-kfg -h,--help                      : show this message
EOF
  exit 1
}

SA="drone"
NS="kube-system"

set_secrets() {
  local prefix="${1}"
  local context="${2}"
  local repo_name="${3}"

  local cluster=$(kubectl config view -o=jsonpath="{.contexts[?(@.name==\"${context}\")].context.cluster}")
  local api_server=$(kubectl config view -o=jsonpath="{.clusters[?(@.name==\"${cluster}\")].cluster.server}")

  local secret=$(kubectl --context $context get -n $NS sa $SA -o jsonpath="{.secrets[].name}")
  local kube_token=$(kubectl --context $context get -n $NS secret $secret -o jsonpath="" | base64 -D)
  # currently no support for tls verification in drone helm!
  # local kube_ca=$(kubectl --context $context get -n $NS secret $secret -o jsonpath="{.data.ca\.crt}")

  drone secret add --repository ${repo_name} --name ${prefix}_API_SERVER --value ${api_server}
  drone secret add --repository ${repo_name} --name ${prefix}_KUBERNETES_TOKEN --value ${kube_token}
}

main() {
  if [[ "$#" -eq 1 ]]; then
    # if [[ "${1}" == '-h' || "${1}" == '--help' ]]; then
      usage
    # fi
  elif [[ "$#" -ne 3 ]]; then
    echo "exactly 3 arguments required"
    usage
  else 
    set_secrets "${1}" "${2}" "${3}"
  fi
}

main "$@"

