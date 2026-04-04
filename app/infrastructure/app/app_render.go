package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/julioguillermo/jgshell/app/infrastructure/components/statusbar"
)

func (a *App) View() tea.View {
	if cmd := a.ctl.LastCmd(); cmd != nil && cmd.IsFullScreen() {
		return tea.NewView(cmd.Output)
	}

	list := a.cmdViewPort.View()
	elements := []string{list}

	if input := a.RenderInput(); input != "" {
		elements = append(elements, input)
	}
	if autocomplete := a.RenderAutocomplete(); autocomplete != "" {
		elements = append(elements, autocomplete)
	}
	if history := a.RenderHistory(); history != "" {
		elements = append(elements, history)
	}
	if status := a.RenderStatus(); status != "" {
		elements = append(elements, status)
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

func (a *App) RenderInput() string {
	if a.ctl.IsRunning() || a.showHistory {
		return ""
	}
	return a.input.View(a.width, a.height)
}

func (a *App) RenderAutocomplete() string {
	if a.ctl.IsRunning() || !a.showAutocomplete || a.showHistory {
		return ""
	}
	return a.autocomplete.View(a.width, a.height)
}

func (a *App) RenderHistory() string {
	if a.ctl.IsRunning() || !a.showHistory {
		return ""
	}
	return a.history.View(a.width, a.height)
}

func (a *App) RenderStatus() string {
	if a.status == nil {
		return ""
	}
	return statusbar.StatusBar(a.status, a.width)
}

func (a *App) FreeHeight() int {
	height := a.height
	if input := a.RenderInput(); input != "" {
		height -= lipgloss.Height(input)
	}
	if autocomplete := a.RenderAutocomplete(); autocomplete != "" {
		height -= lipgloss.Height(autocomplete)
	}
	if history := a.RenderHistory(); history != "" {
		height -= lipgloss.Height(history)
	}
	if status := a.RenderStatus(); status != "" {
		height -= lipgloss.Height(status)
	}
	return height
}
