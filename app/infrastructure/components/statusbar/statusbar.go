package statusbar

import (
	"strings"

	"charm.land/lipgloss/v2"
	statedomain "github.com/julioguillermo/jgshell/state/domain"
)

func StatusBar(status statedomain.Status, width int) string {
	os := GetOS(status.OS())
	shell := GetShell(status.Shell())
	user := GetUser(status.User())
	pwd := GetPwdHome(status)
	git := GetGit(status)

	left := user + " >> " + pwd
	right := os
	if shell != "" {
		right = shell + " " + right
	}
	if git != "" {
		right = git + " " + right
	}

	size := lipgloss.Width(left) + lipgloss.Width(right)

	space := strings.Repeat(" ", max(width-size-2, 3))

	return lipgloss.NewStyle().
		Width(width).
		Render(" " + left + space + right)
}
