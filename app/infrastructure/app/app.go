package app

import (
	tea "charm.land/bubbletea/v2"
	"github.com/julioguillermo/jgshell/app/infrastructure/components/cmdcard"
	"github.com/julioguillermo/jgshell/app/infrastructure/components/input"
	statedomain "github.com/julioguillermo/jgshell/state/domain"
)

type App struct {
	state statedomain.State

	cmdViewPort *cmdcard.CmdViewPort

	width  int
	height int

	input *input.Input
}

func NewApp(state statedomain.State) *App {
	a := &App{
		state:       state,
		cmdViewPort: &cmdcard.CmdViewPort{},
		width:       80,
		height:      24,
	}
	a.input = input.New(state, a.onSend)
	return a
}

func (a *App) Init() tea.Cmd {
	return doTick()
}

func (a *App) onSend(msg string) {
	a.state.Send(msg)
	a.cmdViewPort.GoToBottom()
}
