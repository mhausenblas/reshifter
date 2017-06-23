#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

type docker >/dev/null 2>&1 || { echo >&2 "Need Docker but it's not installed."; exit 1; }
type http >/dev/null 2>&1 || { echo >&2 "Need http command but it's not installed. You can get it from https://httpie.org"; exit 1; }
type jq >/dev/null 2>&1 || { echo >&2 "Need jq command but it's not installed. You can get it from https://stedolan.github.io/jq/"; exit 1; }


etcd3up () {
  dr=$(docker run --rm -d -p 2379:2379 --name test-etcd --dns 8.8.8.8 quay.io/coreos/etcd:v3.1.0 /usr/local/bin/etcd --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379)
  sleep 3s
}

etcddown () {
  dr=$(docker kill test-etcd)
}


printf "\n=========================================================================\n"
printf "Ramping up etcd3 and populating it with a few keys:\n"
etcd3up
curl localhost:2379/v2/keys/foo -XPUT -d value="bar"
curl localhost:2379/v2/keys/that/here -XPUT -d value="moar"
curl localhost:2379/v2/keys/this:also -XPUT -d value="escaped"

printf "\n=========================================================================\n"
printf "Launching ReShifter in the background:\n"
reshifter &
RESHIFTER_PID=$!
sleep 3s

printf "\n=========================================================================\n"
printf "Getting ReShifter version:\n"
http localhost:8080/v1/version

printf "=========================================================================\n"
printf "Performing backup operation:\n"
bid=$(http localhost:8080/v1/backup | jq -r .backupid)
printf "got backup ID %s\n" $bid

printf "\n=========================================================================\n"
printf "Restarting etcd3 ...\n"
etcddown
etcd3up

printf "\n=========================================================================\n"
printf "Performing restore operation:\n"
http localhost:8080/v1/restore?archive=$bid

printf "=========================================================================\n"
printf "Cleaning up:\nremoving local backup file\ntearing down etcd3\n"
rm ./$bid.zip
kr=$(kill $RESHIFTER_PID)
etcddown

printf "\nDONE=====================================================================\n"
