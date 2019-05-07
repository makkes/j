#!/bin/bash

set -euo pipefail

DEST_BASE=~/.j

BASH_PROFILE=
if [[ -w ~/.bashrc ]] ; then
    BASH_PROFILE=~/.bashrc
elif [[ -w ~/.bash_profile ]] ; then
    BASH_PROFILE=~/.bash_profile
else
    echo no bash profile found
    exit 1
fi

if [[ -d $DEST_BASE ]] ; then
    echo j is already installed in $DEST_BASE
    exit 1
fi

echo Downloading j release bundle...

curl -s -L -o /tmp/j.tar.gz https://github.com/makkes/j/releases/download/v1.0.2/j-v1.0.2.tar.gz
mkdir ~/.j
tar -C ~/.j -xzf /tmp/j.tar.gz

if [[ ! $(grep -qc '/.j/j.sh' ${BASH_PROFILE}) ]] ; then
    echo "j source string already in ${BASH_PROFILE}"
else
    echo "source $DEST_BASE/j.sh" >> $BASH_PROFILE
fi

if [[ ! $(grep -qc '/.j/j_completion' ${BASH_PROFILE}) ]] ; then
    echo "bash completion string already in ${BASH_PROFILE}"
else
echo "source $DEST_BASE/j_completion" >> $BASH_PROFILE
fi

echo j has been installed to $DEST_BASE.
echo
echo Close and reopen your terminal to start using j.
