package app

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (a *App) sendPaste(msg tea.PasteMsg) {
	if !a.ctl.IsRunning() {
		return
	}
	a.ctl.Run(msg.Content)
}

func (a *App) sendKey(msg tea.KeyMsg) {
	if !a.ctl.IsRunning() {
		return
	}

	keyStr := msg.String()

	// Handle special keys mapping to ANSI/ASCII
	switch keyStr {
	// Always allow Ctrl+C to send SIGINT
	case "ctrl+c":
		a.ctl.Run("\x03")
	case "enter":
		a.ctl.Run("\r")
	case "backspace":
		a.ctl.Run("\x7f")
	case "tab":
		a.ctl.Run("\t")
	case "esc":
		a.ctl.Run("\x1b")
	case "up":
		a.ctl.Run("\x1b[A")
	case "down":
		a.ctl.Run("\x1b[B")
	case "right":
		a.ctl.Run("\x1b[C")
	case "left":
		a.ctl.Run("\x1b[D")
	case "space":
		a.ctl.Run(" ")
	default:
		if strings.HasPrefix(keyStr, "ctrl+") && len(keyStr) == 6 {
			// Handle ctrl+a through ctrl+z
			char := keyStr[5]
			if char >= 'a' && char <= 'z' {
				a.ctl.Run(string(char - 'a' + 1))
			}
		} else {
			// For normal characters and other Ctrl keys
			a.ctl.Run(string([]rune(msg.Key().Text)))
		}
	}
}
