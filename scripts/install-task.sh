#!/usr/bin/env bash

curl -sL https://taskfile.dev/install.sh | sh
mv bin/task $HOME/.local/bin/
rm -R bin