package app

import (
	tea "charm.land/bubbletea/v2"
)

func (a *App) HandleKeyPress(keypress string) tea.Cmd {
	switch keypress {
	case "alt+shift+e", "ctrl+e":
		return tea.Quit
	}

	return nil
}
