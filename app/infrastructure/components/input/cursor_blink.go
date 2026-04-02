package input

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

type CursorBlink time.Time

func doBlink() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return CursorBlink(t)
	})
}
