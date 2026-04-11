package wrapperdomain

const (
	SimpleWrapper    = `printf "<<<JGSHELL_START;%s>>> %s <<<JGSHELL_END;%s>>>\r"`
	WrapperStartFast = `printf "\033]JGSHELL;START_FAST;%s;%s;>>>\007" "$(whoami)" "$(pwd)"`
	WrapperStart     = `printf "\033]JGSHELL;START;%s;%s;>>>\007" "$(whoami)" "$(pwd)"`
)

const (
	REWrapStartSimple = `<<<JGSHELL_START;([\w\d-_]+?)>>>`
	REWrapStartFast   = `\033]JGSHELL;START_FAST;(.+?);(.+?);>>>\007`
	REWrapStart       = `\033]JGSHELL;START;(.+?);(.+?);>>>\007`
)

const (
	REWrapDoneSimple = `<<<JGSHELL_END;([\w\d-_]+?)>>>`
	REWrapDone       = `\033]JGSHELL;(\d*);DONE\007`
)
