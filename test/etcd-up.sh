#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

docker run --rm -d -p 2379:2379 --name test-etcd elcolio/etcd:2.0.10

curl localhost:2379/v2/keys/foo -XPUT -d value="bar"
curl localhost:2379/v2/keys/that/here -XPUT -d value="moar"
