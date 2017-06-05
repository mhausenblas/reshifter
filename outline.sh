#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

case "$1" in
"create")
    make create
    ;;
"up")
    while true; do
      make --silent
      sleep 1
    done
    ;;
*)
    echo Unknown command, try 'create' or 'up'
    ;;
esac
