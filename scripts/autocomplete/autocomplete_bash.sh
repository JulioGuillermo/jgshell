#!/bin/bash
# Autocomplete for bash
# Usage: bash autocomplete_bash.sh "<line>" <cursor>
# Output: completions, one per line

LINE="$1"
CURSOR="${2:-${#LINE}}"

[ -z "$LINE" ] && exit 0

# Extraer palabra actual
_get_word() {
    local word="${LINE##* }"
    [[ "$LINE" =~ \ $ ]] && word=""
    echo "$word"
}

_get_cmd() {
    echo "${LINE%% *}"
}

word=$(_get_word)
cmd=$(_get_cmd)
[ -z "$cmd" ] && exit 0

# Cargar bash-completion si existe
if [ -f /usr/share/bash-completion/bash_completion ]; then
    source /usr/share/bash-completion/bash_completion 2>/dev/null
fi

# Intentar cargar completions específicas del comando
if [ -z "$(complete -p "$cmd" 2>/dev/null)" ]; then
    if declare -F _completion_loader &>/dev/null; then
        _completion_loader "$cmd" 2>/dev/null
    fi
fi

# Obtener función de completion
comp_line=$(complete -p "$cmd" 2>/dev/null)
comp_func=""

if [ -n "$comp_line" ]; then
    if [[ "$comp_line" =~ -F[[:space:]]+([^[:space:]]+) ]]; then
        comp_func="${BASH_REMATCH[1]}"
    fi
fi

if [ -n "$comp_func" ] && declare -f "$comp_func" >/dev/null 2>&1; then
    read -ra words <<< "$LINE"
    COMP_LINE="$LINE"
    COMP_POINT="$CURSOR"
    COMP_WORDS=("${words[@]}")
    COMP_CWORD=$((${#words[@]} - 1))
    COMP_TYPE=9

    eval "$comp_func" 2>/dev/null

    if [ ${#COMPREPLY[@]} -gt 0 ]; then
        printf '%s\n' "${COMPREPLY[@]}" | sort -u
        exit 0
    fi
fi

# Fallback
if [ -n "$word" ]; then
    {
        compgen -c -- "$word" 2>/dev/null
        compgen -d -- "$word" 2>/dev/null
        compgen -f -- "$word" 2>/dev/null
    } | sort -u
fi
