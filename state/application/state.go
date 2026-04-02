package stateapplication

import (
	"sync"

	shelldomain "github.com/julioguillermo/jgshell/shell/domain"
	statedomain "github.com/julioguillermo/jgshell/state/domain"
)

type State struct {
	mu   *sync.Mutex
	cond *sync.Cond

	shell     shelldomain.Shell
	history   []statedomain.Cmd
	isRunning bool
}

func NewState(shell shelldomain.Shell) *State {
	mu := &sync.Mutex{}
	cond := sync.NewCond(mu)
	state := &State{
		mu:   mu,
		cond: cond,

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
	s.shell.OnClose(func(shelldomain.Shell) {
		f(s)
	})
}
