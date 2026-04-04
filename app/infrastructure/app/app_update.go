package app

import (
	tea "charm.land/bubbletea/v2"
)

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	a.UpdateStatus()

	if a.ToInput() {
		_, c := a.input.Update(msg)
		if c != nil {
			cmds = append(cmds, c)
		}
	}

	if a.ToAutocomplete() {
		_, c := a.autocomplete.Update(msg, a.width, a.height)
		if c != nil {
			cmds = append(cmds, c)
		}
	}

	if a.ToHistory() {
		_, c := a.history.Update(msg, a.width, a.height)
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

	height := a.FreeHeight()
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

func (a *App) ToInput() bool {
	return !a.ctl.IsRunning() && !a.showAutocomplete && !a.showHistory
}

func (a *App) ToAutocomplete() bool {
	return !a.ctl.IsRunning() && a.showAutocomplete && !a.showHistory
}

func (a *App) ToHistory() bool {
	return !a.ctl.IsRunning() && a.showHistory
}
