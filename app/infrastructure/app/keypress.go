package app

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (a *App) HandleKeyPress(msg tea.KeyMsg) tea.Cmd {
	keys := msg.String()
	switch keys {
	case "alt+shift+e", "ctrl+e":
		return tea.Quit
	}

	if a.state.IsRunning() {
		a.sendKey(msg)
		return nil
	}

	switch keys {
	case "tab":
		a.autocomplete.SetItems(a.state.GetAutoComplete(a.input.GetCurrentLine(), a.input.GetCurrentLinePosition()))
		a.showAutocomplete = true
	}

	return nil
}

func (a *App) sendKey(msg tea.KeyMsg) {
	keyStr := msg.String()

	// Handle special keys mapping to ANSI/ASCII
	switch keyStr {
	// Always allow Ctrl+C to send SIGINT
	case "ctrl+c":
		a.state.Send("\x03")
	case "enter":
		a.state.Send("\r")
	case "backspace":
		a.state.Send("\x7f")
	case "tab":
		a.state.Send("\t")
	case "esc":
		a.state.Send("\x1b")
	case "up":
		a.state.Send("\x1b[A")
	case "down":
		a.state.Send("\x1b[B")
	case "right":
		a.state.Send("\x1b[C")
	case "left":
		a.state.Send("\x1b[D")
	case "space":
		a.state.Send(" ")
	default:
		if strings.HasPrefix(keyStr, "ctrl+") && len(keyStr) == 6 {
			// Handle ctrl+a through ctrl+z
			char := keyStr[5]
			if char >= 'a' && char <= 'z' {
				a.state.Write([]byte{char - 'a' + 1})
			}
		} else {
			// For normal characters and other Ctrl keys
			a.state.Send(string(msg.Key().Code))
		}
	}
}
