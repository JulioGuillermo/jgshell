package controllerdomain

import executordomain "github.com/julioguillermo/jgshell/executor/domain"

type ShellController interface {
	// Simple Shell
	WrapShell() error
	SetSize(width, height int) error
	OnClose(f func(ShellController))
	Close() error

	// Runner
	IsRunning() bool
	Run(string) error
	GetHistory() []*executordomain.Cmd
}
