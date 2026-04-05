package app

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/atotto/clipboard"
	"github.com/julioguillermo/jgshell/app/infrastructure/components/cmdcard"
	executordomain "github.com/julioguillermo/jgshell/executor/domain"
)

func (a *App) handleClick(msg tea.MouseClickMsg, height int) {
	if msg.Button != tea.MouseLeft {
		return
	}

	if msg.Y > height {
		return
	}

	cmd, _ := a.getClickedCmd(msg)
	if cmd == nil {
		return
	}

	content := fmt.Sprintf("$ %s\n%s", cmd.Cmd, cmd.CleanOuput())

	clipboard.WriteAll(content)
}

func (a *App) getClickedCmd(msg tea.MouseClickMsg) (*executordomain.Cmd, int) {
	viewportY := msg.Y + a.cmdViewPort.ViewportYOffset()
	offset := 0

	cmds := a.ctl.GetHistory()
	for i, cmd := range cmds {
		card := cmdcard.NewCmdCard(cmd, a.highlighter).View(a.width)
		height := lipgloss.Height(card)
		if viewportY > offset && viewportY < offset+height {
			return cmd, i
		}
		offset += height
	}
	return nil, -1
}
