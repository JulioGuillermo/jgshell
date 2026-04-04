package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/julioguillermo/jgshell/app/infrastructure/components/statusbar"
)

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	a.UpdateStatus()

	if !a.ctl.IsRunning() && !a.showAutocomplete {
		_, c := a.input.Update(msg)
		if c != nil {
			cmds = append(cmds, c)
		}
	}

	if a.showAutocomplete {
		_, c := a.autocomplete.Update(msg, a.width)
		if c != nil {
			cmds = append(cmds, c)
		}
	}

	switch msg := msg.(type) {
	case tickMsg:
		cmds = append(cmds, doTick())
	case tea.KeyMsg:
		if c := a.HandleKeyPress(msg); c != nil {
			cmds = append(cmds, c)
		}
	case tea.WindowSizeMsg:
		if c := a.HandleWindowSize(msg); c != nil {
			cmds = append(cmds, c)
		}
	case tea.PasteMsg:
		a.sendPaste(msg)
	}

	height := a.height
	if !a.ctl.IsRunning() {
		input := a.input.View(a.width, a.height)
		height -= lipgloss.Height(input)
	}
	if a.status != nil {
		state := statusbar.StatusBar(a.status, a.width)
		height -= lipgloss.Height(state)
	}
	if a.showAutocomplete {
		autocomplete := a.autocomplete.View(a.width, a.height)
		height -= lipgloss.Height(autocomplete)
	}

	a.cmdViewPort.Resize(a.width, height)
	a.ctl.SetSize(a.width-2, height-2)

	if _, ok := msg.(tea.WindowSizeMsg); !ok {
		v, cmd := a.cmdViewPort.Update(a.ctl.GetHistory(), a.width, msg)
		a.cmdViewPort = v
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a *App) UpdateStatus() {
	if !a.ctl.IsRunning() && (a.statusDepricated || a.status == nil) {
		a.statusDepricated = false
		a.status, _ = a.ctl.GetStatus()
	}
}
