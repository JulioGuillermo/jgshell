package wrapperdomain

import shelldomain "github.com/julioguillermo/jgshell/shell/domain"

type ShellWrapper interface {
	WrapShell(shelldomain.Shell) error
}
