#!/bin/bash

SCRIPT_PATH="./check_version.sh"
OUTPUT=$("$SCRIPT_PATH")

if [[ $OUTPUT == 1 ]]; then
    go install github.com/golang/mock/mockgen@v1.6.0
else
    GO111MODULE=on go get github.com/golang/mock/mockgen@v1.6.0
fi