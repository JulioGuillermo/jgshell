package app

import (
	tea "charm.land/bubbletea/v2"
	"github.com/julioguillermo/jgshell/app/infrastructure/components/cmdcard"
	"github.com/julioguillermo/jgshell/app/infrastructure/components/input"
	"github.com/julioguillermo/jgshell/app/infrastructure/components/menu"
	controllerdomain "github.com/julioguillermo/jgshell/controller/domain"
	statusdomain "github.com/julioguillermo/jgshell/status/domain"
	syntaxdomain "github.com/julioguillermo/jgshell/syntax/domain"
)

type App struct {
	ctl              controllerdomain.ShellController
	status           *statusdomain.Status
	statusDepricated bool

	cmdViewPort *cmdcard.CmdViewPort
	highlighter syntaxdomain.Highlighter

	width  int
	height int

	input            *input.Input
	autocomplete     *menu.Autocomplete
	showAutocomplete bool
	history          *menu.History
	showHistory      bool
}

func NewApp(ctl controllerdomain.ShellController, highlighter syntaxdomain.Highlighter) *App {
	a := &App{
		ctl:              ctl,
		cmdViewPort:      cmdcard.NewCmdViewPort(80, 24, highlighter),
		highlighter:      highlighter,
		width:            80,
		height:           24,
		statusDepricated: true,
	}

	a.input = input.New(ctl, a.onSend, highlighter)

	a.autocomplete = menu.NewAutocomplete()
	a.autocomplete.OnClose = func() {
		a.showAutocomplete = false
	}
	a.autocomplete.OnSelect = func(item string) {
		a.input.InsertAutocomplete(item)
	}

	a.history = menu.NewHistory()
	a.history.OnClose = func() {
		a.showHistory = false
	}
	a.history.OnSelect = func(item string) {
		a.input.SetValue(item)
	}

	return a
}

func (a *App) Init() tea.Cmd {
	return tea.Batch(
		doTick(),
		a.input.Init(),
	)
}

func (a *App) onSend(msg string) {
	a.ctl.Run(msg)
	a.cmdViewPort.GoToBottom()
	a.statusDepricated = true
}
