#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

docker run --rm -d -p 2379:2379 --name test-etcd3 quay.io/coreos/etcd:v3.1.0 /usr/local/bin/etcd --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379
sleep 5s
curl localhost:2379/v2/keys/foo -XPUT -d value="some"
curl localhost:2379/v2/keys/that/here -XPUT -d value="moar"
