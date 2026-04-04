package app

import (
	tea "charm.land/bubbletea/v2"
)

func (a *App) HandleKeyPress(msg tea.KeyMsg) tea.Cmd {
	keys := msg.String()
	switch keys {
	case "alt+shift+e", "ctrl+e":
		return tea.Quit
	}

	if a.ctl.IsRunning() {
		a.sendKey(msg)
		return nil
	}

	switch keys {
	case "tab":
		if !a.showAutocomplete {
			line := a.input.GetCurrentLine()
			if line != "" {
				items, _ := a.ctl.GetAutocomplete(a.input.GetCurrentLine(), a.input.GetCurrentLinePosition())
				a.autocomplete.SetItems(items)
				a.showAutocomplete = true
			}
		}
	}

	return nil
}
