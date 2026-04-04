package cmdcard

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	executordomain "github.com/julioguillermo/jgshell/executor/domain"
	syntaxdomain "github.com/julioguillermo/jgshell/syntax/domain"
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
	if code == -11 {
		return lipgloss.NewStyle().
			Blink(true).
			Foreground(lipgloss.Color("#ffAA55")).
			Render("")
	}
	if code == -12 {
		return lipgloss.NewStyle().
			Blink(true).
			Foreground(lipgloss.Color("#ff00ff")).
			Render("")
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
	hour := int(duration.Hours())
	min := int(duration.Minutes()) % 60
	sec := int(duration.Seconds()) % 60
	mil := int(duration.Microseconds()) % 1000

	color := "#00ff00"
	switch {
	case hour > 0:
		color = "#ff5555"
	case min > 0:
		color = "#ffff55"
	case sec > 0:
		color = "#00ff55"
	}

	dur := fmt.Sprintf("%02d:%02d:%02d.%03d", hour, min, sec, mil)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Render(dur)
}

func getFancyPWD(cmd *executordomain.Cmd) string {
	if cmd.USER == "" {
		return cmd.PWD
	}

	prefix := fmt.Sprintf("/home/%s", cmd.USER)

	if !strings.HasPrefix(cmd.PWD, prefix) {
		return cmd.PWD
	}

	return " " + strings.TrimPrefix(cmd.PWD, prefix)
}

func getUserPWD(cmd *executordomain.Cmd) string {
	if cmd.USER == "" && cmd.PWD == "" {
		return ""
	}
	user := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff8855")).
		Render(cmd.USER)
	pwd := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ffaa")).
		Render(getFancyPWD(cmd))
	sh := ""
	if cmd.SH != "" {
		sh = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ffaa")).
			Render(cmd.SH) + " "
	}
	return fmt.Sprintf("%s[%s] ➜ %s", sh, user, pwd)
}

func getTopBorder(width int, cmd *executordomain.Cmd) string {
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

func CmdCard(cmd *executordomain.Cmd, width int, highlighter syntaxdomain.Highlighter) string {
	input := highlighter.Highlight(cmd.Cmd)
	output := cmd.CleanOuput()
	borderTop := getTopBorder(width, cmd)

	return borderTop + "\n" +
		lipgloss.NewStyle().
			Width(width).
			MaxWidth(width).
			Border(lipgloss.RoundedBorder(), true).
			BorderTop(false).
			BorderForeground(BorderColor).
			Render(input+"\n\n"+output)
}
