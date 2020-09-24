#!/usr/bin/env bash

curl -sL https://taskfile.dev/install.sh | sh
mv bin/task $HOME/.local/bin/
echo "task binary moved to $HOME/.local/bin"
echo "please do 'source ~/.bashrc' to reload \$PATH"
rm -R bin