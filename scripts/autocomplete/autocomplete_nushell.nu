# Autocomplete for nushell
# Usage: nu autocomplete_nushell.nu "<line>" <cursor>
# Output: completions, one per line

let line = $env.JG_LINE? | default ""
let cursor = $env.JG_CURSOR? | default 0

if ($line | is-empty) { exit }

# Extraer palabra actual
let parts = ($line | split row " ")
let word = ($parts | last | default "")
let cmd = ($parts | first | default "")

if ($cmd | is-empty) { exit }

# Nushell no tiene una API pública de completions programática
# Fallback: comandos y archivos

let completions = (
    # Comandos internos
    ^which $word 2> /dev/null | get name? | default []

    union (
        # Archivos en PATH
        $env.PATH
        | split row (if $nu.os-info.family == "windows" { ";" } else { ":" })
        | each { |dir|
            try { ls $dir | where name =~ $"^($word)" | get name } catch { [] }
        }
        | flatten
    )

    union (
        # Archivos locales
        try { ls | where name =~ $"^($word)" | get name } catch { [] }
    )
)

if ($completions | is-not-empty) {
    $completions | uniq | sort
}
