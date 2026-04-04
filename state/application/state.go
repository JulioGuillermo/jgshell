package stateapplication

import (
	"sync"

	shelldomain "github.com/julioguillermo/jgshell/shell/domain"
	statedomain "github.com/julioguillermo/jgshell/state/domain"
)

type State struct {
	locker sync.Locker
	cond   *sync.Cond

	shell     shelldomain.FullShell
	history   []statedomain.Cmd
	isRunning bool
}

func NewState(shell shelldomain.FullShell) *State {
	mu := &sync.Mutex{}
	cond := sync.NewCond(mu)
	state := &State{
		locker: mu,
		cond:   cond,

		shell:   shell,
		history: []statedomain.Cmd{},
	}
	state.startReader()
	return state
}

func (s *State) Close() error {
	return s.shell.Close()
}

func (s *State) OnClose(f func(s statedomain.State)) {
	s.shell.OnClose(func(shelldomain.FullShell) {
		f(s)
	})
}

func (s *State) SetSize(width, height int) {
	s.shell.SetSize(uint16(height), uint16(width))
}
