#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

type kubectl >/dev/null 2>&1 || { echo >&2 "Need kubecuddle but it's not installed."; exit 1; }


###############################################################################
# MAIN
if kubectl version --short | tail -1 | awk '{ print $3 }' | grep 1.5 > /dev/null
then
  printf "Generating synthetic testbed, this can take up to xx seconds!\n"
  if ! kubectl get ns | grep reshifter-testbed > /dev/null
  then
    printf "I need a namespace called 'reshifter-testbed' to proceed, please create one.\nABORTING.\n\n"
    exit 1
  fi
else
  printf "Sorry, unsupported version of Kubernetes, can only do 1.5 and above.\nABORTING.\n\n"
  exit 1
fi

kubectl create -f https://raw.githubusercontent.com/mhausenblas/kbe/master/specs/pods/pod.yaml
kubectl create -f https://raw.githubusercontent.com/mhausenblas/kbe/master/specs/labels/pod.yaml
kubectl create -f https://raw.githubusercontent.com/mhausenblas/kbe/master/specs/services/rc.yaml
kubectl create -f https://raw.githubusercontent.com/mhausenblas/kbe/master/specs/services/svc.yaml
kubectl create -f https://raw.githubusercontent.com/mhausenblas/kbe/master/specs/healthz/pod.yaml
kubectl create -f https://raw.githubusercontent.com/mhausenblas/kbe/master/specs/volumes/pod.yaml
kubectl create secret generic apikey --from-literal=thekey=supersecret
kubectl create -f https://raw.githubusercontent.com/mhausenblas/kbe/master/specs/secrets/pod.yaml
kubectl create -f https://raw.githubusercontent.com/mhausenblas/kbe/master/specs/jobs/job.yaml

printf "\nDONE=====================================================================\n"
