#!/bin/bash

function j() {
    JUMPCMD=jump
    DIR=`${JUMPCMD} "$@"`
    builtin cd "$DIR"
}
