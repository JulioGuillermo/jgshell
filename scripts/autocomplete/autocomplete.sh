# Completion bridge script para jgshell
# Variables de entrada: GO_SHELL, GO_LINE, GO_CURSOR (pasadas desde Go)
# Salida: lista de completions, uno por línea

SHELL_TYPE="%{GO_SHELL}"
LINE="%{GO_LINE}"
CURSOR="%{GO_CURSOR}"


# Extraer la palabra actual (lo que va después del último espacio)
_get_word() {
    local word="${LINE##* }"
    # Si la línea termina en espacio, la palabra está vacía
    [[ "$LINE" =~ \ $ ]] && word=""
    echo "$word"
}

# Extraer el comando (primera palabra)
_get_cmd() {
    echo "${LINE%% *}"
}

# === BASH ===
_completion_bash() {
    local word=$(_get_word)
    local cmd=$(_get_cmd)
    [ -z "$cmd" ] && return

    # Cargar sistema de bash-completion si existe
    if [ -f /usr/share/bash-completion/bash_completion ]; then
        source /usr/share/bash-completion/bash_completion 2>/dev/null
    fi

    # Intentar cargar completions para el comando
    if [ -z "$(complete -p "$cmd" 2>/dev/null)" ]; then
        if declare -F _completion_loader &>/dev/null; then
            _completion_loader "$cmd" 2>/dev/null
        fi
    fi

    # Obtener la función de completion
    local comp_line
    comp_line=$(complete -p "$cmd" 2>/dev/null)
    local comp_func=""

    if [ -n "$comp_line" ]; then
        # Parsear: complete -F _func cmd o complete -c cmd
        if [[ "$comp_line" =~ -F[[:space:]]+([^[:space:]]+) ]]; then
            comp_func="${BASH_REMATCH[1]}"
        fi
    fi

    if [ -n "$comp_func" ] && declare -f "$comp_func" >/dev/null 2>&1; then
        local -a words
        read -ra words <<< "$LINE"
        COMP_LINE="$LINE"
        COMP_POINT="$CURSOR"
        COMP_WORDS=("${words[@]}")
        COMP_CWORD=$((${#words[@]} - 1))
        COMP_TYPE=9  # Menú completion

        eval "$comp_func" 2>/dev/null

        if [ ${#COMPREPLY[@]} -gt 0 ]; then
            printf '%s\n' "${COMPREPLY[@]}" | sort -u
            return
        fi
    fi

    # Fallback: comandos, directorios, archivos
    if [ -n "$word" ]; then
        {
            compgen -c -- "$word" 2>/dev/null
            compgen -d -- "$word" 2>/dev/null
            compgen -f -- "$word" 2>/dev/null
        } | sort -u
    fi
}

# === ZSH ===
_completion_zsh() {
    local word=$(_get_word)
    local cmd=$(_get_cmd)
    [ -z "$cmd" ] && return

    # zsh no puede usar su sistema de completion desde bash
    # Fallback a completions básicas
    if [ -n "$word" ]; then
        {
            compgen -c -- "$word" 2>/dev/null
            compgen -d -- "$word" 2>/dev/null
            compgen -f -- "$word" 2>/dev/null
        } | sort -u
    fi
}

# === FISH ===
_completion_fish() {
    local word=$(_get_word)
    local cmd=$(_get_cmd)
    [ -z "$cmd" ] && return

    # fish tiene su propio sistema, no accesible desde bash
    # Fallback: buscar en PATH
    if [ -n "$word" ]; then
        local IFS=:
        for dir in $PATH; do
            [ -d "$dir" ] && ls -1 "$dir" 2>/dev/null | grep "^${word}"
        done | sort -u

        compgen -d -- "$word" 2>/dev/null
        compgen -f -- "$word" 2>/dev/null
    fi
}

# === NUSHELL ===
_completion_nushell() {
    # nushell no es accesible desde bash
    # Fallback genérico
    local word=$(_get_word)
    if [ -n "$word" ]; then
        {
            compgen -c -- "$word" 2>/dev/null
            compgen -d -- "$word" 2>/dev/null
            compgen -f -- "$word" 2>/dev/null
        } | sort -u
    fi
}

# === ELVISH / XONSH / OTROS ===
_completion_other() {
    local word=$(_get_word)
    if [ -n "$word" ]; then
        {
            compgen -c -- "$word" 2>/dev/null
            compgen -d -- "$word" 2>/dev/null
            compgen -f -- "$word" 2>/dev/null
        } | sort -u
    fi
}

# === MAIN ===
_run() {
    [ -z "$LINE" ] && return
    [ -z "$CURSOR" ] && CURSOR=0
    case "$SHELL_TYPE" in
        bash)    _completion_bash ;;
        zsh)     _completion_zsh ;;
        fish)    _completion_fish ;;
        nushell) _completion_nushell ;;
        *)       _completion_other ;;
    esac
}

_run
