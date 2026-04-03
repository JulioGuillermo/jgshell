# Autocomplete for PowerShell
# Usage: pwsh -File autocomplete_powershell.ps1 -Line "<line>" -Cursor <cursor>
# Output: completions, one per line

param(
    [string]$Line,
    [int]$Cursor = 0
)

if ([string]::IsNullOrEmpty($Line)) { exit 0 }
if ($Cursor -eq 0) { $Cursor = $Line.Length }

# Usar TabExpansion2 para obtener completions reales
try {
    $result = [System.Management.Automation.CommandCompletion]::CompleteInput(
        $Line,
        $Cursor,
        @{}
    )

    if ($null -ne $result.CompletionMatches -and $result.CompletionMatches.Count -gt 0) {
        $result.CompletionMatches | ForEach-Object {
            $_.CompletionText
        } | Sort-Object -Unique
        exit 0
    }
} catch {
    # Fallback si falla
}

# Fallback: comandos y archivos
$word = ($Line -split '\s+')[-1]
if ($Line -match '\s$') { $word = "" }

if (-not [string]::IsNullOrEmpty($word)) {
    # Comandos
    Get-Command "$word*" -ErrorAction SilentlyContinue | ForEach-Object { $_.Name }

    # Archivos/directorios
    $dir = Split-Path $word -Parent 2>$null
    if ([string]::IsNullOrEmpty($dir)) { $dir = "." }
    $leaf = Split-Path $word -Leaf 2>$null
    Get-ChildItem -Path $dir -Filter "$leaf*" -ErrorAction SilentlyContinue | ForEach-Object { $_.Name }
} | Sort-Object -Unique
