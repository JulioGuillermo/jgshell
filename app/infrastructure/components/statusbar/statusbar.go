package statusbar

import (
	"strings"

	"charm.land/lipgloss/v2"
	statedomain "github.com/julioguillermo/jgshell/state/domain"
)

func StatusBar(status statedomain.Status, width int) string {
	os := GetOS(status.OS())
	user := GetUser(status.User())
	pwd := GetPwdHome(status)
	git := GetGit(status)

	left := user + " >> " + pwd
	right := os
	if git != "" {
		right = git + " " + os
	}

	size := lipgloss.Width(left) + lipgloss.Width(right)

	space := strings.Repeat(" ", max(width-size, 3))

	return lipgloss.NewStyle().
		Width(width).
		Render(left + space + right)
}
