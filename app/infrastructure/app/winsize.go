package app

import tea "charm.land/bubbletea/v2"

func (a *App) HandleWindowSize(size tea.WindowSizeMsg) tea.Cmd {
	a.width = size.Width
	a.height = size.Height
	return nil
}
