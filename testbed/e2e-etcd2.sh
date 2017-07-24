#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

type docker >/dev/null 2>&1 || { echo >&2 "Need Docker but it's not installed."; exit 1; }
type http >/dev/null 2>&1 || { echo >&2 "Need http command but it's not installed. You can get it from https://httpie.org"; exit 1; }
type jq >/dev/null 2>&1 || { echo >&2 "Need jq command but it's not installed. You can get it from https://stedolan.github.io/jq/"; exit 1; }

etcd2up () {
  dr=$(docker run --rm -d -p 2379:2379 --name test-etcd --dns 8.8.8.8 quay.io/coreos/etcd:v2.3.8 --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379)
  sleep 3s
}

etcd2secureup () {
  dr=$(docker run --rm -d -v $(pwd)/certs/:/etc/ssl/certs -p 2379:2379 --name test-etcd --dns 8.8.8.8 quay.io/coreos/etcd:v2.3.8 \
  --ca-file /etc/ssl/certs/ca.pem --cert-file /etc/ssl/certs/server.pem --key-file /etc/ssl/certs/server-key.pem \
  --advertise-client-urls https://0.0.0.0:2379 --listen-client-urls https://0.0.0.0:2379)
  sleep 3s
}

etcddown () {
  dk=$(docker kill test-etcd)
}

populate() {
  curl http://localhost:2379/v2/keys/kubernetes.io/namespaces/kube-system -XPUT -d value="."
  curl http://localhost:2379/v2/keys/openshift.io -XPUT -d value="."
}

populatesecure() {
  curl --cacert $(pwd)/certs/ca.pem --cert $(pwd)/certs/client.p12 --pass reshifter -L https://localhost:2379/v2/keys/kubernetes.io/namespaces/kube-system -XPUT -d value="."
  curl --cacert $(pwd)/certs/ca.pem --cert $(pwd)/certs/client.p12 --pass reshifter -L https://localhost:2379/v2/keys/openshift.io -XPUT -d value="."
}

doversion() {
  printf "\n=========================================================================\n"
  printf "Getting ReShifter version:\n"
  http localhost:8080/v1/version
}

dobackup() {
  printf "=========================================================================\n"
  printf "Performing backup operation:\n"
  todaybucket=reshifter-test-$(date "+%Y-%m-%d")
  bid=$(http POST localhost:8080/v1/backup endpoint=$1 remote=play.minio.io:9000 bucket=$todaybucket | jq -r .backupid)
  printf "got backup ID %s\n" $bid
}

dorestore() {
  printf "\n=========================================================================\n"
  printf "Performing restore operation:\n"
  todaybucket=reshifter-test-$(date "+%Y-%m-%d")
  http POST localhost:8080/v1/restore endpoint=$1 backupid=$bid remote=play.minio.io:9000 bucket=$todaybucket
}

cleanup () {
  printf "=========================================================================\n"
  printf "Cleaning up:\nremoving local backup file\ntearing down etcd2\n"
  pkill -f reshifter
  etcddown
}

###############################################################################
# MAIN
printf "Test plan: etcd2 and secure etcd2, this can take up to 30 seconds!\n"

# using the Minio play backend:
export ACCESS_KEY_ID=Q3AM3UQ867SPQQA43P2F
export SECRET_ACCESS_KEY=zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG

# main test plan etcd2:
printf "\n=========================================================================\n"
printf "Ramping up etcd2 and populating it with a few keys:\n"
etcd2up
populate
printf "\n=========================================================================\n"
printf "Launching ReShifter in the background:\n"
reshifter &
RESHIFTER_PID=$!
sleep 3s
doversion
dobackup http://localhost:2379
etcddown
etcd2up
dorestore http://localhost:2379
cleanup
printf "\nDONE=====================================================================\n"

sleep 3s

# main test plan etcd2 secure:
export RS_ETCD_CLIENT_CERT=$(pwd)/certs/client.pem
export RS_ETCD_CLIENT_KEY=$(pwd)/certs/client-key.pem
export RS_ETCD_CA_CERT=$(pwd)/certs/ca.pem
printf "\n=========================================================================\n"
printf "Ramping up secure etcd2 and populating it with a few keys:\n"
etcd2secureup
populatesecure
printf "\n=========================================================================\n"
printf "Launching ReShifter in the background:\n"
reshifter &
RESHIFTER_PID=$!
sleep 3s
doversion
dobackup https://localhost:2379
etcddown
etcd2secureup
dorestore https://localhost:2379
cleanup
printf "\nDONE=====================================================================\n"
