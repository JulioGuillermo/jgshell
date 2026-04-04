#!/usr/bin/env pwsh
# Markers:
#   START: \x1b]123;START\x07<user> <pwd> >>>
#   END:   \x1b]123;<exit_code>;<uuid>;DONE\x07

if (Get-Command Stop-OhMyPosh -ErrorAction SilentlyContinue) { Stop-OhMyPosh }
$global:PromptPurge = $true

function global:prompt {
    $exitCode = if ($null -eq $LASTEXITCODE) { [int][bool]$? } else { $LASTEXITCODE }
    return "\033]123;$exitCode;DONE\007"
}
