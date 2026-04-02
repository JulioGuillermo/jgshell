package statusbar

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	statedomain "github.com/julioguillermo/jgshell/state/domain"
)

func GetGit(status statedomain.Status) string {
	//                            
	git := status.Git()
	if git == nil {
		return ""
	}

	gitIcon := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff5500")).
		Render("  ")

	branch := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffff00")).
		Render(git.BranchLocal)

	var sb strings.Builder
	sb.WriteString("(" + gitIcon + branch)

	if git.Untracked > 0 {
		sb.WriteString(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00ff00")).
				Render(fmt.Sprintf("  %d", git.Untracked)),
		)
	}

	if git.Modified > 0 {
		sb.WriteString(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffff00")).
				Render(fmt.Sprintf("  %d", git.Modified)),
		)
	}

	if git.Deleted > 0 {
		sb.WriteString(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff5500")).
				Render(fmt.Sprintf("  %d", git.Deleted)),
		)
	}

	if git.Staged > 0 {
		sb.WriteString(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00ffaa")).
				Render(fmt.Sprintf("  %d", git.Staged)),
		)
	}

	if git.Conflicts > 0 {
		sb.WriteString(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff0000")).
				Render(fmt.Sprintf("  %d", git.Conflicts)),
		)
	}

	if git.Ahead > 0 {
		sb.WriteString(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00ffff")).
				Render(fmt.Sprintf("  %d", git.Ahead)),
		)
	}

	if git.Behind > 0 {
		sb.WriteString(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00ffff")).
				Render(fmt.Sprintf("  %d", git.Behind)),
		)
	}

	return sb.String() + " )"
}
