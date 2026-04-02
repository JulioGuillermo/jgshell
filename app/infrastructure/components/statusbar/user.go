package statusbar

import "charm.land/lipgloss/v2"

func GetUser(user string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffff00")).
		Render(user)
}
