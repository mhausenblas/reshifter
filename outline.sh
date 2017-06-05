#!/usr/bin/env bash

## Note that this script depends on
## https://github.com/emcrisostomo/fswatch

set -o errexit
set -o errtrace
set -o nounset
# set -o pipefail

type fswatch >/dev/null 2>&1 || { echo >&2 "Need fswatch but it's not installed, so exiting …"; exit 1; }

case "$1" in
"create")
    make create
    ;;
"up")
    while true; do
      fswatch -0 $(pwd)/app/* | make
    done
    ;;
*)
    echo Unknown command, try 'create' or 'up'
    ;;
esac
