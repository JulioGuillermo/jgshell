# Markers:
#   START: \x1b]123;START\x07<user> <pwd> >>>
#   END:   \x1b]123;<exit_code>;<uuid>;DONE\x07

unsetopt PROMPT_SP
unset HISTFILE

# _jg_emit_start_marker() {
#     printf '\033]123;START;%s;%s;>>>\007' "$(whoami)" "$(pwd)"
# }

# add-zsh-hook preexec _jg_emit_start_marker 2>/dev/null || {
#     autoload -Uz add-zsh-hook 2>/dev/null
#     add-zsh-hook preexec _jg_emit_start_marker
# }

PS1='$(printf "\033]123;$?;DONE\007")'
