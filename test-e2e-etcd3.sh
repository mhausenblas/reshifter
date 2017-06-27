#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

type docker >/dev/null 2>&1 || { echo >&2 "Need Docker but it's not installed."; exit 1; }
type etcdctl >/dev/null 2>&1 || { echo >&2 "Need etcdctl but it's not installed."; exit 1; }
type http >/dev/null 2>&1 || { echo >&2 "Need http command but it's not installed. You can get it from https://httpie.org"; exit 1; }
type jq >/dev/null 2>&1 || { echo >&2 "Need jq command but it's not installed. You can get it from https://stedolan.github.io/jq/"; exit 1; }

etcd3up () {
  dr=$(docker run --rm -d -p 2379:2379 --name test-etcd --dns 8.8.8.8 --env ETCD_DEBUG quay.io/coreos/etcd:v3.1.0 /usr/local/bin/etcd  \
  --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379 --listen-peer-urls http://0.0.0.0:2380)
  sleep 3s
}

etcd3secureup () {
  dr=$(docker run --rm -d -v $(pwd)/certs/:/etc/ssl/certs -p 2379:2379 --name test-etcd --dns 8.8.8.8 quay.io/coreos/etcd:v3.1.0 /usr/local/bin/etcd \
  --ca-file /etc/ssl/certs/ca.pem --cert-file /etc/ssl/certs/server.pem --key-file /etc/ssl/certs/server-key.pem \
  --advertise-client-urls https://0.0.0.0:2379 --listen-client-urls https://0.0.0.0:2379)
  sleep 3s
}

etcddown () {
  dr=$(docker kill test-etcd)
}

populate() {
  etcdctl --endpoints=http://localhost:2379 put /kubernetes.io/namespaces/kube-system "."
}

populatesecure() {
  # etcdctl --endpoints=http://localhost:2379 set /kubernetes.io/namespaces/kube-system .
  # curl --cacert $(pwd)/certs/ca.pem --cert $(pwd)/certs/client.p12 --pass reshifter -L https://localhost:2379/v2/keys/foo -XPUT -d value="bar"
  echo TBD
}

doversion() {
  printf "\n=========================================================================\n"
  printf "Getting ReShifter version:\n"
  http localhost:8080/v1/version
}

dobackup() {
  printf "=========================================================================\n"
  printf "Performing backup operation:\n"
  bid=$(http POST localhost:8080/v1/backup endpoint=$1 | jq -r .backupid)
  printf "got backup ID %s\n" $bid
}

dorestore() {
  printf "\n=========================================================================\n"
  printf "Performing restore operation:\n"
  http POST localhost:8080/v1/restore endpoint=$1 archive=$bid
}

cleanup () {
  printf "=========================================================================\n"
  printf "Cleaning up:\nremoving local backup file\ntearing down etcd3\n"
  rm ./$bid.zip
  pkill -f reshifter
  etcddown
}

###############################################################################
# MAIN
printf "Test plan: etcd3 and secure etcd3, this can take up to 30 seconds!\n"
export ETCDCTL_API=3

# main test plan etcd3:
printf "\n=========================================================================\n"
printf "Ramping up etcd3 and populating it with a few keys:\n"
etcd3up
populate
sleep 2
printf "\n=========================================================================\n"
printf "Launching ReShifter in the background:\n"
DEBUG=true reshifter &
sleep 3s
doversion
dobackup http://localhost:2379
# etcddown
# etcd3up
# dorestore http://localhost:2379
cleanup
printf "\nDONE=====================================================================\n"

# sleep 3s

# main test plan etcd3 secure:
# export RS_ETCD_CLIENT_CERT=$(pwd)/certs/client.pem
# export RS_ETCD_CLIENT_KEY=$(pwd)/certs/client-key.pem
# export RS_ETCD_CA_CERT=$(pwd)/certs/ca.pem
# printf "\n=========================================================================\n"
# printf "Ramping up secure etcd3 and populating it with a few keys:\n"
# etcd3secureup
# populatesecure
# printf "\n=========================================================================\n"
# printf "Launching ReShifter in the background:\n"
# reshifter &
# RESHIFTER_PID=$!
# sleep 3s
# doversion
# dobackup https://localhost:2379
# etcddown
# etcd3secureup
# dorestore https://localhost:2379
# cleanup
# printf "\nDONE=====================================================================\n"
