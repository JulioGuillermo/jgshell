package cmdcard

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	statedomain "github.com/julioguillermo/jgshell/state/domain"
)

var BorderColor = lipgloss.Color("#8800ff")

func spinner(duration time.Duration) string {
	Map := []rune("")
	// Map := []rune("⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏")
	// Map := []rune("")
	index := int(duration.Milliseconds()/20) % len(Map)
	return string(Map[index])
}

func getStatusCode(code int, duration time.Duration) string {
	if code == -10 {
		return lipgloss.NewStyle().
			Blink(true).
			Foreground(lipgloss.Color("#ffff55")).
			Render(spinner(duration))
	}
	if code == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff55")).
			Render("✓")
	}
	if code == 1 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff5555")).
			Render("✗")
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff5555")).
		Render(fmt.Sprintf("ERR %d", code))
}

func getDuration(duration time.Duration) string {
	hour := duration.Hours()
	min := duration.Minutes()
	sec := duration.Seconds()
	mil := duration.Microseconds()

	color := "#00ff00"
	switch {
	case hour > 0:
		color = "#ff5555"
	case min > 0:
		color = "#ffff55"
	case sec > 0:
		color = "#00ff55"
	}

	dur := fmt.Sprintf("%02d:%02d:%02d.%03d", int(hour), int(min)%60, int(sec)%60, mil%1000)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Render(dur)
}

func getFancyPWD(cmd statedomain.Cmd) string {
	if cmd.USER == "" {
		return cmd.PWD
	}

	prefix := fmt.Sprintf("/home/%s", cmd.USER)

	if !strings.HasPrefix(cmd.PWD, prefix) {
		return cmd.PWD
	}

	return " " + strings.TrimPrefix(cmd.PWD, prefix)
}

func getUserPWD(cmd statedomain.Cmd) string {
	if cmd.USER == "" && cmd.PWD == "" {
		return ""
	}
	user := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff8855")).
		Render(cmd.USER)
	pwd := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ffaa")).
		Render(getFancyPWD(cmd))
	return fmt.Sprintf("%s ➜ %s", user, pwd)
}

func getTopBorder(width int, cmd statedomain.Cmd) string {
	duration := cmd.GetDuration()

	left := getUserPWD(cmd)
	right := fmt.Sprintf("%s ➜ %s", getStatusCode(cmd.ExitCode, duration), getDuration(duration))

	leftSize := lipgloss.Width(left)
	rightSize := lipgloss.Width(right)
	totalSize := leftSize + rightSize

	border := strings.Repeat("─", max(width-totalSize-10, 1))
	style := lipgloss.NewStyle().
		Foreground(BorderColor)

	return style.Render("╭[ ") +
		left +
		style.Render(" ]"+border+"( ") +
		right +
		style.Render(" )╮")
}

func CmdCard(cmd statedomain.Cmd, width int) string {
	output := cmd.CleanOuput()
	borderTop := getTopBorder(width, cmd)

	return borderTop + "\n" +
		lipgloss.NewStyle().
			Width(width).
			Border(lipgloss.RoundedBorder(), true).
			BorderTop(false).
			BorderForeground(BorderColor).
			Render(cmd.Cmd+"\n"+output)
}
