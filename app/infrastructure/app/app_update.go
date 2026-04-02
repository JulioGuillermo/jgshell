package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/julioguillermo/jgshell/app/infrastructure/components/statusbar"
)

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if _, c := a.input.Update(msg); c != nil {
		cmds = append(cmds, c)
	}

	switch msg := msg.(type) {
	case tickMsg:
		cmds = append(cmds, doTick())
	case tea.KeyPressMsg:
		if c := a.HandleKeyPress(msg.String()); c != nil {
			cmds = append(cmds, c)
		}
	case tea.WindowSizeMsg:
		if c := a.HandleWindowSize(msg); c != nil {
			cmds = append(cmds, c)
		}

		state := statusbar.StatusBar(a.state, a.width)
		input := a.input.View(a.width, a.height)

		height := a.height
		if a.state.ShowInput() {
			height -= lipgloss.Height(input)
		}
		if a.state.ShowStatusBar() {
			height -= lipgloss.Height(state)
		}

		a.cmdViewPort.Resize(a.width, height)
		a.state.SetSize(a.width-2, height)
	}

	if _, ok := msg.(tea.WindowSizeMsg); !ok {
		v, cmd := a.cmdViewPort.Update(a.state.GetHistory(), a.width, msg)
		a.cmdViewPort = v
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}
