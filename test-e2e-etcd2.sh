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
  dr=$(docker run -d  -v $(pwd)/certs/:/etc/ssl/certs -p 2379:2379 --name test-etcd --dns 8.8.8.8 quay.io/coreos/etcd:v2.3.8 \
  --ca-file /etc/ssl/certs/ca.pem --cert-file /etc/ssl/certs/server.pem --key-file /etc/ssl/certs/server-key.pem \
  --advertise-client-urls https://0.0.0.0:2379 --listen-client-urls https://0.0.0.0:2379)
  sleep 3s
}

etcddown () {
  dr=$(docker kill test-etcd)
}

restartetcd () {
  printf "\n=========================================================================\n"
  printf "Restarting etcd2 ...\n"
  etcddown
  $1
}

populate() {
  curl localhost:2379/v2/keys/foo -XPUT -d value="bar"
  curl localhost:2379/v2/keys/that/here -XPUT -d value="moar"
  curl localhost:2379/v2/keys/this:also -XPUT -d value="escaped"
  curl localhost:2379/v2/keys/other -XPUT -d value="value"
}

populatesecure() {
  curl --cacert $(pwd)/certs/ca.pem --cert $(pwd)/certs/client.p12 --pass reshifter -L https://127.0.0.1:2379/v2/keys/foo -XPUT -d value="bar"
  curl --cacert $(pwd)/certs/ca.pem --cert $(pwd)/certs/client.p12 --pass reshifter -L https://127.0.0.1:2379/v2/keys/that/here -XPUT -d value="moar"
  curl --cacert $(pwd)/certs/ca.pem --cert $(pwd)/certs/client.p12 --pass reshifter -L https://127.0.0.1:2379/v2/keys/this:also -XPUT -d value="escaped"
  curl --cacert $(pwd)/certs/ca.pem --cert $(pwd)/certs/client.p12 --pass reshifter -L https://127.0.0.1:2379/v2/keys/other -XPUT -d value="value"
}

doversion() {
  printf "\n=========================================================================\n"
  printf "Getting ReShifter version:\n"
  http localhost:8080/v1/version
}

dobackup() {
  printf "=========================================================================\n"
  printf "Performing backup operation:\n"
  bid=$(http localhost:8080/v1/backup | jq -r .backupid)
  printf "got backup ID %s\n" $bid
}

dorestore() {
  printf "\n=========================================================================\n"
  printf "Performing restore operation:\n"
  http localhost:8080/v1/restore?archive=$bid
}

cleanup () {
  printf "=========================================================================\n"
  printf "Cleaning up:\nremoving local backup file\ntearing down etcd2\n"
  rm ./$bid.zip
  kr=$(kill $RESHIFTER_PID)
  etcddown
}

###############################################################################
# MAIN
printf "Test plan: etcd2 and secure etcd2, this can take up to 30 seconds!\n"

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
dobackup
restartetcd etcd2up
dorestore
cleanup
printf "\nDONE=====================================================================\n"

sleep 5s

# main test plan etcd2 secure:
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
dobackup
restartetcd etcd2secureup
dorestore
cleanup
printf "\nDONE=====================================================================\n"
