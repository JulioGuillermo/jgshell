package app

import (
	tea "charm.land/bubbletea/v2"
	"github.com/julioguillermo/jgshell/app/infrastructure/components/cmdcard"
	"github.com/julioguillermo/jgshell/app/infrastructure/components/input"
	statedomain "github.com/julioguillermo/jgshell/state/domain"
	syntaxdomain "github.com/julioguillermo/jgshell/syntax/domain"
)

type App struct {
	state            statedomain.State
	status           statedomain.Status
	statusDepricated bool

	cmdViewPort *cmdcard.CmdViewPort
	highlighter syntaxdomain.Highlighter

	width  int
	height int

	input *input.Input
}

func NewApp(state statedomain.State, highlighter syntaxdomain.Highlighter) *App {
	a := &App{
		state:            state,
		cmdViewPort:      cmdcard.NewCmdViewPort(80, 24, highlighter),
		highlighter:      highlighter,
		width:            80,
		height:           24,
		statusDepricated: true,
	}
	a.input = input.New(state, a.onSend, highlighter)
	return a
}

func (a *App) Init() tea.Cmd {
	return tea.Batch(
		doTick(),
		a.input.Init(),
	)
}

func (a *App) onSend(msg string) {
	a.state.Send(msg)
	a.cmdViewPort.GoToBottom()
	a.statusDepricated = true
}
