#!/bin/bash

OS="unknown"

# GET OS
if [[ "$OSTYPE" == "linux-android"* ]]; then
    OS="android"
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS="linux"
elif [[ "$OSTYPE" == "freebsd"* ]]; then
    OS="freebsd"
elif [[ "$OSTYPE" == "openbsd"* ]]; then
    OS="openbsd"
elif [[ "$OSTYPE" == "netbsd"* ]]; then
    OS="netbsd"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    if [[ $(uname -m) == "iPhone"* ]] || [[ $(uname -m) == "iPad"* ]]; then
        OS="ios"
    else
        OS="mac"
    fi
elif [[ "$OS" == "Windows_NT" ]]; then
    OS="windows"
fi

# Get user
USER=$(id -u -n 2>/dev/null || echo $USER)
if [ -z "$USER" ]; then
    USER=$(whoami)
fi

# Get current directory
CURRENT_DIR=$(pwd 2>/dev/null || echo $PWD)
if [ -z "$CURRENT_DIR" ]; then
    CURRENT_DIR="."
fi
CURRENT_DIR=$(echo "$CURRENT_DIR" | tr -d '\r\n')

# Get shell
if [ -n "$PSVersionTable" ]; then
    CURRENT_SHELL="powershell"
elif [ -n "$BASH_VERSION" ]; then
    CURRENT_SHELL="bash"
elif [ -n "$ZSH_VERSION" ]; then
    CURRENT_SHELL="zsh"
elif [ -n "$KSH_VERSION" ]; then
    CURRENT_SHELL="ksh"
else
    # Fallback: intentar sacar el nombre del proceso padre
    CURRENT_SHELL=$(ps -p $$ -o comm= 2>/dev/null | sed 's/-//' || echo "sh")
fi

echo "OS: ($OS)"
echo "User: ($USER)"
echo "Dir: ($CURRENT_DIR)"
echo "SHELL: ($CURRENT_SHELL)"
echo "=== GIT START ==="
git status -sb 2>/dev/null || echo "NO GIT"
echo "=== GIT END ==="
