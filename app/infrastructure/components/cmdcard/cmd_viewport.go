package cmdcard

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	statedomain "github.com/julioguillermo/jgshell/state/domain"
	syntaxdomain "github.com/julioguillermo/jgshell/syntax/domain"
)

type CmdViewPort struct {
	highlighter syntaxdomain.Highlighter
	viewport    viewport.Model
	bottom      bool
}

func NewCmdViewPort(width, height int, highlighter syntaxdomain.Highlighter) *CmdViewPort {
	return &CmdViewPort{
		viewport:    viewport.New(),
		highlighter: highlighter,
	}
}

func (m *CmdViewPort) Update(cmds []statedomain.Cmd, width int, msg tea.Msg) (*CmdViewPort, tea.Cmd) {
	cards := make([]string, 0, len(cmds))
	for _, c := range cmds {
		card := CmdCard(c, width, m.highlighter)
		cards = append(cards, card)
	}
	m.viewport.SetContent(lipgloss.JoinVertical(
		lipgloss.Left,
		cards...,
	))

	vp, cmd := m.viewport.Update(msg)
	m.viewport = vp

	switch msg.(type) {
	case tea.MouseMsg:
		m.bottom = m.viewport.AtBottom()
	}

	return m, cmd
}

func (m CmdViewPort) View() string {
	if m.bottom {
		m.viewport.GotoBottom()
	}

	return m.viewport.View()
}

func (m *CmdViewPort) Resize(width, height int) {
	m.viewport.SetWidth(width)
	m.viewport.SetHeight(height)
}

func (m *CmdViewPort) GoToBottom() {
	m.bottom = true
}
