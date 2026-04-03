#!/bin/sh
# Autocomplete for POSIX sh (dash, ash, mksh, etc.)
# Usage: sh autocomplete_sh.sh "<line>" <cursor>
# Output: completions, one per line

LINE="$1"
CURSOR="${2:-0}"

[ -z "$LINE" ] && exit 0

# Extraer palabra actual
_get_word() {
    _word="${LINE##* }"
    case "$LINE" in
        *" ") _word="" ;;
    esac
    echo "$_word"
}

# Extraer comando
_get_cmd() {
    echo "${LINE%% *}"
}

word=$(_get_word)
cmd=$(_get_cmd)
[ -z "$cmd" ] && exit 0

# Fallback: buscar en PATH
if [ -n "$word" ]; then
    # Comandos en PATH
    oldIFS="$IFS"
    IFS=:
    for dir in $PATH; do
        [ -d "$dir" ] && ls -1 "$dir" 2>/dev/null | grep "^${word}"
    done | sort -u
    IFS="$oldIFS"

    # Directorios
    for f in ${word}*/; do
        [ -d "$f" ] && echo "$f"
    done 2>/dev/null

    # Archivos
    for f in ${word}*; do
        [ -e "$f" ] && echo "$f"
    done 2>/dev/null
fi
