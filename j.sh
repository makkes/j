#!/bin/bash

function j() {
    JUMPCMD=~/.j/jump
    DIR=`${JUMPCMD} "$@"`
    builtin cd "$DIR"
}
