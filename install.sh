#!/bin/bash

set -euo pipefail

DEST_BASE=~/.j

if [[ -d $DEST_BASE ]] ; then
    echo j is already installed in $DEST_BASE
    exit 1
fi

curl -L -o /tmp/j.tar.gz https://github.com/makkes/j/releases/download/v1.0-rc1/j-v1.0.0-rc1.tar.gz
mkdir ~/.j
tar -C ~/.j -xzf /tmp/j.tar.gz

BASH_PROFILE=
if [[ -w ~/.bashrc ]] ; then
    BASH_PROFILE=~/.bashrc
elif [[ -w ~/.bash_profile ]] ; then
    BASH_PROFILE=~/.bash_profile
else
    echo no bash profile found
    exit 1
fi

echo "source $DEST_BASE/j.sh" >> $BASH_PROFILE
echo "source $DEST_BASE/j_completion" >> $BASH_PROFILE

