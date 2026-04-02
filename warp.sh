[ -z $WARP_BOOTSTRAPPED ] && printf "\\e]9278;f;{\"hook\": \"InitSubshell\", \"value\": { \"shell\": \"%s\", \"uname\": \"%s\" }}\\a" $([ $FISH_VERSION ] && echo "fish" || { echo $0 | command -p grep -q zsh && echo "zsh"; } || { echo $0 | command -p grep -q bash && echo "bash"; } || echo "unknown") $(uname)
export WARP_HONOR_PS1=0
command -p stty raw
unset PROMPT_COMMAND
HISTCONTROL=ignorespace
HISTIGNORE=" *"
WARP_IS_SUBSHELL=1
WARP_SESSION_ID="$(command -p date +%s)$RANDOM"
_hostname=$(command -pv hostname >/dev/null 2>&1 && command -p hostname 2>/dev/null || command -p uname -n)
_user=$(command -pv whoami >/dev/null 2>&1 && command -p whoami 2>/dev/null || echo $USER)
_msg=$(printf "{\"hook\": \"InitShell\", \"value\": {\"session_id\": $WARP_SESSION_ID, \"shell\": \"bash\", \"user\": \"%s\", \"hostname\": \"%s\", \"is_subshell\": true}}" "$_user" "$_hostname" | command -p od -An -v -tx1 | command -p tr -d " \n")
if [[ "$OS" == Windows_NT ]]; then
   WARP_IN_MSYS2=true
else
   WARP_IN_MSYS2=false
fi
WARP_USING_WINDOWS_CON_PTY=false
if [ "$WARP_USING_WINDOWS_CON_PTY" = true ]; then
   printf '"'"'\e]9278;d;%s\x07'"'"' "$_msg"
else
   printf '"'"'\e\x50\x24\x64%s\x9c'"'"' "$_msg"
fi
unset _hostname _user _msg
