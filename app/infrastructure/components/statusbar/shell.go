package statusbar

import "charm.land/lipgloss/v2"

func GetShell(shell string) string {
	//  
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#0088ff")).
		Render(" " + shell)
}
