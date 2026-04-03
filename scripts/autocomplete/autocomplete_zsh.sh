#!/bin/zsh
# Autocomplete for zsh
# Usage: zsh autocomplete_zsh.sh "<line>" <cursor>
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

# Inicializar sistema de completions
autoload -Uz compinit 2>/dev/null
if ! compinit -u 2>/dev/null; then
    # Fallback si compinit falla
    if [ -n "$word" ]; then
        {
            print -l ${(f)"$(compgen -c -- "$word" 2>/dev/null)"}
            print -l ${(f)"$(compgen -d -- "$word" 2>/dev/null)"}
            print -l ${(f)"$(compgen -f -- "$word" 2>/dev/null)"}
        } | sort -u
    fi
    exit 0
fi

# Usar el sistema de completions de zsh
# _main_complete ejecuta la completion para el comando actual
zle -C _complete main_complete _complete 2>/dev/null

# Método: usar _describe o llamar la función de completion directamente
compfunc="_${cmd//-/_}"

if declare -f "$compfunc" >/dev/null 2>&1; then
    # Preparar contexto de completion
    set -- ${(z)LINE}
    local -a compreply
    compset -P "${LINE%%${word}*}" 2>/dev/null
    compset -S '*' 2>/dev/null

    # Llamar la función de completion
    eval "$compfunc" 2>/dev/null

    if [ ${#compreply[@]} -gt 0 ]; then
        printf '%s\n' "${compreply[@]}" | sort -u
        exit 0
    fi
fi

# Fallback: usar _main_complete con captura
# Truco: usar zsh -c con compadd
replies=()
compadd -a replies 2>/dev/null

# Método alternativo: usar la función de completion del sistema
if declare -f "_${cmd}" >/dev/null 2>&1; then
    local -a matches
    _arguments "*: :(${(f)"$(compgen -c -- "$word")"})" 2>/dev/null
fi

# Último fallback
if [ -n "$word" ]; then
    {
        compgen -c -- "$word" 2>/dev/null
        compgen -d -- "$word" 2>/dev/null
        compgen -f -- "$word" 2>/dev/null
    } | sort -u
fi
