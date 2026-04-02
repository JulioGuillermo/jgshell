package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/julioguillermo/jgshell/app/infrastructure/components/statusbar"
)

func (a *App) View() tea.View {
	list := a.cmdViewPort.View()
	elements := []string{list}

	if a.state.ShowInput() {
		input := a.input.View(a.width, a.height)
		elements = append(elements, input)
	}

	if a.state.ShowStatusBar() {
		state := statusbar.StatusBar(a.state, a.width)
		elements = append(elements, state)
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
