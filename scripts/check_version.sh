#!/bin/bash

# Adapted from https://unix.stackexchange.com/a/285928/416548

CURRENT_VERSION="$(go version | { read _ _ v _; echo ${v#go}; })"
REQUIRED_VERSION="1.16"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$CURRENT_VERSION" | sort -V | head -n1)" = "$REQUIRED_VERSION" ]; then
      echo "1"
else
      echo "0"
fi
