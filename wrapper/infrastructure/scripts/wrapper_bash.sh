set +o history

_jg_prompt() {
    local EXIT_CODE=$?
    printf '\033]JGSHELL;%s;DONE\007' "$EXIT_CODE"
}

PS1='$(_jg_prompt) '
