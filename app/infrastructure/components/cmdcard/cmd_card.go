package cmdcard

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	executordomain "github.com/julioguillermo/jgshell/executor/domain"
	syntaxdomain "github.com/julioguillermo/jgshell/syntax/domain"
)

type CmdCard struct {
	Cmd         *executordomain.Cmd
	Highlighter syntaxdomain.Highlighter
}

func NewCmdCard(cmd *executordomain.Cmd, highlighter syntaxdomain.Highlighter) *CmdCard {
	return &CmdCard{
		Cmd:         cmd,
		Highlighter: highlighter,
	}
}

func (c *CmdCard) View(width int) string {
	input := c.Highlighter.Highlight(c.Cmd.Cmd)
	output := c.Cmd.CleanOuput()
	borderTop := c.getTopBorder(width)

	return borderTop + "\n" +
		lipgloss.NewStyle().
			Width(width).
			MaxWidth(width).
			Border(lipgloss.RoundedBorder(), true).
			BorderTop(false).
			BorderForeground(c.getBorderColor()).
			Render(input+"\n\n"+output)
}

func (c *CmdCard) getBorderColor() color.Color {
	if c.Cmd.ExitCode == 0 {
		return lipgloss.Color("#aa00ff")
	}
	if c.Cmd.ExitCode == -10 {
		return lipgloss.Color("#ffff00")
	}
	if c.Cmd.ExitCode == -11 {
		return lipgloss.Color("#ffAA88")
	}
	if c.Cmd.ExitCode == -12 {
		return lipgloss.Color("#0000ff")
	}
	return lipgloss.Color("#ff5555")
}

func (c *CmdCard) getTopBorder(width int) string {
	left := c.getUserPWD()
	right := fmt.Sprintf("%s ➜ %s", c.getStatusCode(), c.getDuration())

	leftSize := lipgloss.Width(left)
	rightSize := lipgloss.Width(right)
	totalSize := leftSize + rightSize

	border := strings.Repeat("─", max(width-totalSize-10, 1))
	style := lipgloss.NewStyle().
		Foreground(c.getBorderColor())

	return style.Render("╭[ ") +
		left +
		style.Render(" ]"+border+"( ") +
		right +
		style.Render(" )╮")
}

func (c *CmdCard) getUserPWD() string {
	if c.Cmd.USER == "" && c.Cmd.PWD == "" {
		return ""
	}
	user := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff8855")).
		Render(c.Cmd.USER)
	pwd := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ffaa")).
		Render(c.getFancyPWD())
	sh := ""
	if c.Cmd.SH != "" {
		sh = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ffaa")).
			Render(c.Cmd.SH) + " "
	}
	return fmt.Sprintf("%s[%s] ➜ %s", sh, user, pwd)
}

func (c *CmdCard) getFancyPWD() string {
	if c.Cmd.USER == "" {
		return c.Cmd.PWD
	}

	prefix := fmt.Sprintf("/home/%s", c.Cmd.USER)

	if !strings.HasPrefix(c.Cmd.PWD, prefix) {
		return c.Cmd.PWD
	}

	return " " + strings.TrimPrefix(c.Cmd.PWD, prefix)
}

func (c *CmdCard) getDuration() string {
	duration := c.Cmd.GetDuration()

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

func (c *CmdCard) getStatusCode() string {
	code := c.Cmd.ExitCode

	if code == -10 {
		return lipgloss.NewStyle().
			Blink(true).
			Foreground(lipgloss.Color("#ffff55")).
			Render(c.spinner())
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

func (c *CmdCard) spinner() string {
	Map := []rune("")
	// Map := []rune("⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏")
	// Map := []rune("")
	duration := c.Cmd.GetDuration()
	index := int(duration.Milliseconds()/20) % len(Map)
	return string(Map[index])
}
