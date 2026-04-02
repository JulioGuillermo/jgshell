package statusbar

import (
	"charm.land/lipgloss/v2"
	statedomain "github.com/julioguillermo/jgshell/state/domain"
)

func StatusBar(state statedomain.State, width int) string {
	return lipgloss.NewStyle().
		Width(width).
		// Background(lipgloss.Color("#112233")).
		Foreground(lipgloss.Color("#00ff88")).
		Render("")
}
