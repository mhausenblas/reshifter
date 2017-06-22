#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

type http >/dev/null 2>&1 || { echo >&2 "Need http command but it's not installed. You can get it from https://httpie.org"; exit 1; }
type jq >/dev/null 2>&1 || { echo >&2 "Need jq command but it's not installed. You can get it from https://stedolan.github.io/jq/"; exit 1; }


etcd2up () {
  dr=$(docker run --rm -d -p 2379:2379 --name test-etcd --dns 8.8.8.8 quay.io/coreos/etcd:v2.3.8 --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379)
  sleep 3s
}

etcd2down () {
  dr=$(docker kill test-etcd)
}


printf "\n=========================================================================\n"
printf "Ramping up etcd2 and populating it with a few keys:\n"
etcd2up
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
printf "Restarting etcd2 ...\n"
etcd2down
etcd2up

printf "\n=========================================================================\n"
printf "Performing restore operation:\n"
http localhost:8080/v1/restore?archive=$bid

printf "=========================================================================\n"
printf "Cleaning up:\nremoving local backup file\ntearing down etcd2\n"
rm ./$bid.zip
kr=$(kill $RESHIFTER_PID)
etcd2down

printf "\nDONE=====================================================================\n"
