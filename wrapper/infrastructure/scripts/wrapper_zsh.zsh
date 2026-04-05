unsetopt PROMPT_SP
unset HISTFILE

PS1='$(printf "\033]JGSHELL;$?;DONE\007")'
