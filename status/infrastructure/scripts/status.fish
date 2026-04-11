#!/usr/bin/fish

set -l OS "unknown"

set -l KERNEL (uname -s)
set -l ARCH (uname -m)

if test "$KERNEL" = "Linux"
    if set -q PREFIX; and string match -q "*com.termux*" "$PREFIX"
        set OS "android"
    else
        set OS "linux"
    end
else if test "$KERNEL" = "Darwin"
    if string match -q "iPhone*" "$ARCH"; or string match -q "iPad*" "$ARCH"
        set OS "ios"
    else
        set OS "mac"
    end
else if string match -q "*BSD" "$KERNEL"
    set OS (string lower "$KERNEL")
else if string match -q "*NT*" "$OS"
    set OS "windows"
end

# Get user
set -l USER (id -u -n 2>/dev/null; or echo "$USER")
if test -z "$USER"
    set USER (whoami)
end

# Get current directory - Usamos string collect para evitar que se rompa con espacios
set -l DIR (pwd 2>/dev/null; or echo "$PWD" | string collect)
if test -z "$DIR"
    set DIR "."
end
set DIR (string trim "$DIR")

set -l raw (git status -sb 2>/dev/null; or echo "NO GIT")
set -l GIT_STATUS (string join "\n" $raw)


printf "OS: (%s)\nUser: (%s)\nDir: (%s)\n=== GIT START ===\n" "$OS" "$USER" "$DIR"
printf "$GIT_STATUS"
printf "\n=== GIT END ==="
