#!/usr/bin/env pwsh

if (Get-Command Stop-OhMyPosh -ErrorAction SilentlyContinue)
{
    Stop-OhMyPosh
}
$global:PromptPurge = $true

function global:prompt
{
    $exitCode = if ($null -eq $LASTEXITCODE)
    {
        [int][bool]$?
    } else
    {
        $LASTEXITCODE
    }
    printf "\033]JGSHELL;$exitCode;DONE\007"
    # return "`e]JGSHELL;$exitCode;DONE`a"
    return "`n"
}
