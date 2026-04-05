#!/bin/sh

export HISTSIZE=0

PS1='$(printf "\033]JGSHELL;$?;DONE\007")'
