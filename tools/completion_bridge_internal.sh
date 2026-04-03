#!/bin/bash
# Completion bridge script para ejecutar DENTRO del shell del PTY
# Este script se ejecuta en el contexto del shell existente, no en un bash nuevo
# Uso desde jgshell (vía PTY): source <(echo "$(cat completion_bridge_internal.sh)"; echo "get_completions \"\$LINE\" \$CURSOR")

# --- Funciones de completion (se cargan en el shell) ---

_completion_bridge_bash() {
    local word="${1:-}"
    [ -z "$word" ] && return
    
    # Cargar bash-completion
    [ -f /usr/share/bash-completion/bash_completion ] && source /usr/share/bash-completion/bash_completion 2>/dev/null
    
    # Obtener comando
    local cmd="${COMP_WORDS[0]:-}"
    [ -z "$cmd" ] && return
    
    # Cargar completion del comando si no existe
    [ -z "$(complete -p "$cmd" 2>/dev/null)" ] && declare -F _completion_loader &>/dev/null && _completion_loader "$cmd" 2>/dev/null
    
    # Obtener función de completion
    local comp_func=$(complete -p "$cmd" 2>/dev/null | awk '{print $(NF-1)}')
    
    if [ -n "$comp_func" ] && declare -f "$comp_func" >/dev/null 2>&1; then
        eval "$comp_func" 2>/dev/null
        [ ${#COMPREPLY[@]} -gt 0 ] && printf '%s\n' "${COMPREPLY[@]}" | sort -u && return
    fi
    
    # Fallback
    compgen -c "$word" 2>/dev/null
    compgen -d "$word" 2>/dev/null
    compgen -f -- "$word" 2>/dev/null
}

_completion_bridge_zsh() {
    local word="$1"
    [ -z "$word" ] && return
    
    autoload -Uz compinit 2>/dev/null
    
    # Intentar compadd
    compadd -A -q - "$word" 2>/dev/null
    
    # Si no hay resultados, intentar función de completion
    local cmd="${COMP_WORDS[0]:-}"
    [ -z "$cmd" ] && return
    
    local compfunc="_${cmd}"
    if declare -f "$compfunc" >/dev/null 2>&1; then
        eval "$compfunc" 2>/dev/null
    fi
}

# Función principal que llama según el shell
_bridge_completions() {
    local line="$1"
    local cursor="$2"
    
    if [ -n "$ZSH_VERSION" ]; then
        _completion_bridge_zsh "$(echo "$line" | awk '{print $NF}')"
    elif [ -n "$BASH_VERSION" ]; then
        _completion_bridge_bash "$(echo "$line" | awk '{print $NF}')"
    else
        compgen -c "$(echo "$line" | awk '{print $NF}')" 2>/dev/null
    fi
}

# Ejecutar y output
_bridge_completions "$1" "$2"