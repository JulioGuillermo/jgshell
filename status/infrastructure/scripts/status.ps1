$OSName = "unknown"

# GET OS
$ostype = $env:OSTYPE
if ($ostype -like "linux-android*")
{
    $OSName = "android"
} elseif ($ostype -like "linux-gnu*")
{
    $OSName = "linux"
} elseif ($ostype -like "freebsd*")
{
    $OSName = "freebsd"
} elseif ($ostype -like "openbsd*")
{
    $OSName = "openbsd"
} elseif ($ostype -like "netbsd*")
{
    $OSName = "netbsd"
} elseif ($ostype -like "darwin*")
{
    $arch = ""
    try
    {
        $arch = (& uname -m 2>$null | Out-String).Trim()
    } catch
    {
        $arch = ""
    }

    if ($arch -like "iPhone*" -or $arch -like "iPad*")
    {
        $OSName = "ios"
    } else
    {
        $OSName = "mac"
    }
} elseif ($env:OS -eq "Windows_NT")
{
    $OSName = "windows"
} else
{
    if ([System.Runtime.InteropServices.RuntimeInformation]::IsOSPlatform([System.Runtime.InteropServices.OSPlatform]::Windows))
    {
        $OSName = "windows"
    } elseif ([System.Runtime.InteropServices.RuntimeInformation]::IsOSPlatform([System.Runtime.InteropServices.OSPlatform]::Linux))
    {
        $OSName = "linux"
    } elseif ([System.Runtime.InteropServices.RuntimeInformation]::IsOSPlatform([System.Runtime.InteropServices.OSPlatform]::OSX))
    {
        $OSName = "mac"
    }
}

# Get user
$UserName = ""
try
{
    $UserName = (& id -u -n 2>$null | Out-String).Trim()
} catch
{
    $UserName = ""
}

if ([string]::IsNullOrWhiteSpace($UserName))
{
    $UserName = $env:USER
}
if ([string]::IsNullOrWhiteSpace($UserName))
{
    $UserName = $env:USERNAME
}
if ([string]::IsNullOrWhiteSpace($UserName))
{
    $UserName = [System.Environment]::UserName
}

# Get current directory
$Dir = ""
try
{
    $Dir = (Get-Location).Path
} catch
{
    $Dir = ""
}
if ([string]::IsNullOrWhiteSpace($Dir))
{
    if ($null -ne $PWD)
    {
        $Dir = $PWD.Path
    }
}
if ([string]::IsNullOrWhiteSpace($Dir))
{
    $Dir = "."
}
$Dir = [regex]::Replace($Dir, "[\r\n]", "")

$GitStatus = "NO GIT"
try
{
    $gitOutput = (& git status -sb 2>$null | Out-String).TrimEnd()
    if (-not [string]::IsNullOrWhiteSpace($gitOutput))
    {
        $GitStatus = $gitOutput
    }
} catch
{
    $GitStatus = "NO GIT"
}

# [Console]::Write(("OS: ({0})`nUser: ({1})`nDir: ({2})`n=== GIT START ===`n{3}`n=== GIT END ===" -f $OSName, $UserName, $Dir, $GitStatus))
printf "OS: (%s)`nUser: (%s)`nDir: (%s)`n=== GIT START ===`n%s`n=== GIT END ===" $OSName $UserName $Dir $GitStatus
