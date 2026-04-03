# Get shell
if [ -n "$PSVersionTable" ]; then
    SHELL="powershell"
# elif [ -n "$ZSH_VERSION" ]; then
#     SHELL="zsh"
# elif [ -n "$BASH_VERSION" ]; then
#     SHELL="bash"
# elif [ -n "$FISH_VERSION" ]; then
#     SHELL="fish"
# elif [ -n "$KSH_VERSION" ]; then
#     SHELL="ksh"
# elif [ -n "$ELVISH_VERSION" ]; then
#     SHELL="elvish"
# elif [ -n "$XONSH_VERSION" ]; then
#     SHELL="xonsh"
else
    # sh, dash, ash, mksh, nushell, etc. have no version variables
    # Detect from process name
    SHELL=$(ps -p $$ -o comm= 2>/dev/null | sed 's/-//' || echo "sh")
fi

echo "$SHELL"
