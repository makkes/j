#!/bin/bash

set -uo pipefail

function rollback_and_exit() {
    if [[ -n $OLD_INST ]] ; then
        echo "Cancelling upgrade and rolling back to old version..."
        rm -rf "$DEST_BASE"
        mv "$OLD_INST" "$DEST_BASE"
    fi
    exit "$1"
}

function cleanup_and_exit() {
    if [[ -n $OLD_INST ]] ; then
        rm -rf "$OLD_INST"
    fi
    exit "$1"
}

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

OLD_INST=
if [[ -d $DEST_BASE ]] ; then
    OLD_INST=~/.j-$(date +%s)
    echo j is already installed in $DEST_BASE. Trying to replace it...
    mv "$DEST_BASE" "$OLD_INST" || cleanup_and_exit 1
fi

echo Downloading j release bundle...

version="v1.0.7"
curl -s -L -o /tmp/j.tar.gz https://github.com/makkes/j/releases/download/${version}/j_${version}_"$(uname -o|sed 's/^GNU\///'|awk '{print tolower($0)}')"_"$(uname -m)".tar.gz || rollback_and_exit 1
mkdir ~/.j || exit 1
tar -C ~/.j -xzf /tmp/j.tar.gz || rollback_and_exit 1

if grep -qc '/.j/j.sh' ${BASH_PROFILE}; then
    echo "j source string already in ${BASH_PROFILE}"
else
    echo "source $DEST_BASE/j.sh" >> $BASH_PROFILE || rollback_and_exit 1
fi

if grep -qc '/.j/j_completion' ${BASH_PROFILE}; then
    echo "bash completion string already in ${BASH_PROFILE}"
else
echo "source $DEST_BASE/j_completion" >> $BASH_PROFILE || rollback_and_exit 1
fi

echo j has been installed to $DEST_BASE.
echo
echo Close and reopen your terminal to start using j.

cleanup_and_exit 0
