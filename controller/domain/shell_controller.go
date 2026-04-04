package controllerdomain

import (
	executordomain "github.com/julioguillermo/jgshell/executor/domain"
	statusdomain "github.com/julioguillermo/jgshell/status/domain"
)

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
	LastCmd() *executordomain.Cmd

	// Features
	GetStatus() (*statusdomain.Status, error)
	GetAutocomplete(line string, cursor int) ([]string, error)
	GetCmdHistory() []string
	Filter(start string) []string
	FilterLast(start string) string
}
