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

echo "OS: ($OS)"
echo "User: ($USER)"
echo "Dir: ($CURRENT_DIR)"
echo "=== GIT START ==="
git status -sb 2>/dev/null || echo "NO GIT"
echo "=== GIT END ==="
