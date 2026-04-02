package statusbar

import (
	"strings"

	"charm.land/lipgloss/v2"
	statedomain "github.com/julioguillermo/jgshell/state/domain"
)

func GetPwd(pwd string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff0088")).
		Render(pwd)
}

func GetPwdHome(status statedomain.Status) string {
	prefix := "/home/" + status.User()
	pwd := status.Dir()
	hasPrefix := false

	if after, ok := strings.CutPrefix(pwd, prefix); ok {
		prefix = " "
		pwd = after
		hasPrefix = true
	}
	pwd = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff0088")).
		Render(pwd)

	if !hasPrefix {
		return pwd
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff8800")).
		Render(prefix) + pwd
}
