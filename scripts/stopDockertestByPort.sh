#!/usr/bin/env bash

# Source: https://stackoverflow.com/a/56953427/1033134
# usage:
#   scripts/stopByPort.sh 5434

for id in $(docker ps -q)
do
    if [[ $(docker port "${id}") == *"${1}"* ]]; then
        echo "stopping container ${id}"
        docker stop "${id}"
    fi
done