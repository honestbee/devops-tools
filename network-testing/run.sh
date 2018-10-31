#!/usr/bin/env bash
set -eu

kubectl apply -f manifests/

until $(kubectl get pods -l name=iperf3-server -o jsonpath='{.items[0].status.containerStatuses[0].ready}'); do
    echo "server is not ready yet..."
    sleep 5
done

pods=$(kubectl get pods -l name=iperf3-client -o name | cut -d'/' -f2)

for pod in ${pods}; do
    until $(kubectl get pod ${pod} -o jsonpath='{.status.containerStatuses[0].ready}'); do
        echo "${pod} is not ready yet..."
        sleep 5
    done
    ip=$(kubectl get pod ${pod} -o jsonpath='{.status.hostIP}')
    kubectl exec -it ${pod} -- iperf3 -c iperf3-server -T "Client on ${ip}" $@
    echo
done

kubectl delete -f manifests/
