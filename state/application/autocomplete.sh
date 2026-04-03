#!/bin/bash
# Completion bridge script para jgshell
# Variables de entrada: GO_LINE, GO_CURSOR (pasadas desde Go)
# Salida: lista de completions, uno por línea

LINE="%{GO_LINE}"
CURSOR="%{GO_CURSOR}"

[ -z "$LINE" ] && exit 0
[ -z "$CURSOR" ] && CURSOR=0

# Detectar shell actual
_detect_shell() {
    if [ -n "$ZSH_VERSION" ]; then echo "zsh"
    elif [ -n "$BASH_VERSION" ]; then echo "bash"
    elif [ -n "$FISH_VERSION" ]; then echo "fish"
    else echo "sh"
    fi
}

# === SH (POSIX) ===
_completion_sh() {
    local word="${LINE##* }"
    [ -z "$word" ] && return

    local IFS=:
    for dir in $PATH; do
        [ -d "$dir" ] && ls -1 "$dir" 2>/dev/null | grep "^${word}" | head -20
    done | sort -u

    compgen -d "$word" 2>/dev/null
    compgen -f -- "$word" 2>/dev/null
}

# === BASH ===
_completion_bash() {
    local word="${LINE##* }"
    [[ "$LINE" =~ \ $ ]] && word=""

    [ -f /usr/share/bash-completion/bash_completion ] && source /usr/share/bash-completion/bash_completion 2>/dev/null

    local cmd="${LINE%% *}"
    [ -z "$cmd" ] && exit 0

    [ -z "$(complete -p "$cmd" 2>/dev/null)" ] && declare -F _completion_loader &>/dev/null && _completion_loader "$cmd" 2>/dev/null

    local comp_func=$(complete -p "$cmd" 2>/dev/null | awk '{print $(NF-1)}')

    if [ -n "$comp_func" ] && declare -f "$comp_func" >/dev/null 2>&1; then
        local -a words=($LINE)
        COMP_LINE="$LINE"
        COMP_POINT="$CURSOR"
        COMP_WORDS=("${words[@]}")
        COMP_CWORD=$((${#words[@]} - 1))
        eval "$comp_func" 2>/dev/null
        [ ${#COMPREPLY[@]} -gt 0 ] && printf '%s\n' "${COMPREPLY[@]}" | sort -u && return
    fi

    [ -n "$word" ] && { compgen -c "$word"; compgen -d "$word"; compgen -f -- "$word"; } | sort -u
}

# === ZSH ===
_completion_zsh() {
    local word="${LINE##* }"
    local cmd="${LINE%% *}"
    [ -z "$cmd" ] && exit 0

    autoload -Uz compinit 2>/dev/null

    local compfunc="_${cmd}"
    if declare -f "$compfunc" >/dev/null 2>&1; then
        compset -n $(( ${#LINE} - ${#word} ))
        eval "$compfunc" 2>/dev/null
        [ ${#completions[@]} -gt 0 ] && printf '%s\n' "${completions[@]}" | sort -u && return
    fi

    compgen -d "$word" 2>/dev/null
    compgen -f -- "$word" 2>/dev/null
}

# === FISH ===
_completion_fish() {
    local word="${LINE##* }"
    local cmd="${LINE%% *}"
    [ -z "$cmd" ] && exit 0

    local completions_dir="$HOME/.config/fish/completions"

    if [ -f "${completions_dir}/${cmd}.fish" ]; then
        source "${completions_dir}/${cmd}.fish" 2>/dev/null
    fi

    local IFS=:
    for dir in $PATH; do
        ls -1 "$dir" 2>/dev/null | grep "^${word}" | head -10
    done | sort -u
}

# === MAIN ===
SHELL=$(_detect_shell)

case "$SHELL" in
    zsh)     _completion_zsh ;;
    bash)    _completion_bash ;;
    fish)    _completion_fish ;;
    sh|*)    _completion_sh ;;
esac
