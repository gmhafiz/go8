#!/usr/bin/env bash

TASKPATH=$(which task)
if [ -z "$TASKPATH" ]
then
  curl -sL https://taskfile.dev/install.sh | sh
  mv bin/task $HOME/.local/bin/
  echo "task binary moved to $HOME/.local/bin"
  echo "please do 'source ~/.bashrc' to reload \$PATH"
  rm -R bin
else
  echo "Task has already been installed"
fi
