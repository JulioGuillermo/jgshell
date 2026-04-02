package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/julioguillermo/jgshell/app/infrastructure/components/statusbar"
)

func (a *App) View() tea.View {
	list := a.cmdViewPort.View()
	elements := []string{list}

	if !a.state.IsRunning() {
		input := a.input.View(a.width, a.height)
		state := statusbar.StatusBar(a.state, a.width)
		elements = append(elements, input, state)
	}

	v := tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			elements...,
		),
	)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion

	return v
}
