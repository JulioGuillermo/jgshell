package app

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Millisecond*20, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
