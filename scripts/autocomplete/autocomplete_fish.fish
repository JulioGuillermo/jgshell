#!/usr/bin/env fish
# Autocomplete for fish
# Usage: fish autocomplete_fish.fish "<line>" <cursor>
# Output: completions, one per line

set LINE "$argv[1]"
set CURSOR "$argv[2]"

if test -z "$LINE"
    exit 0
end

if test -z "$CURSOR"
    set CURSOR (string length "$LINE")
end

# Extraer palabra actual
set word (string match -r '[^ ]*$' "$LINE" | string trim)
if string match -q '* ' "$LINE"
    set word ""
end

# Extraer comando
set cmd (string split ' ' "$LINE")[1]

if test -z "$cmd"
    exit 0
end

# Fish tiene un sistema de completions nativo accesible via __fish_complete
# Método 1: usar complete -C para obtener completions
set completions (complete -C "$LINE" 2>/dev/null)

if test -n "$completions"
    printf '%s\n' $completions | sort -u
    exit 0
end

# Método 2: fallback manual
if test -n "$word"
    # Comandos en PATH
    for dir in (string split ':' $PATH)
        if test -d "$dir"
            ls -1 "$dir" 2>/dev/null | string match -r "^$word.*"
        end
    end | sort -u

    # Directorios
    set dirs (string match -r "^$word.*/" */ 2>/dev/null)
    if test -n "$dirs"
        printf '%s\n' $dirs
    end
end
