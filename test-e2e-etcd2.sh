#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

docker run --rm -d -p 2379:2379 --name test-etcd --dns 8.8.8.8 quay.io/coreos/etcd:v2.3.8 --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379

sleep 3s

curl localhost:2379/v2/keys/foo -XPUT -d value="bar"
curl localhost:2379/v2/keys/that/here -XPUT -d value="moar"
curl localhost:2379/v2/keys/this:also -XPUT -d value="escaped"
