package input

import (
	"strings"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	statedomain "github.com/julioguillermo/jgshell/state/domain"
)

type Input struct {
	state    statedomain.State
	textarea textarea.Model
	onSend   func(string)
}

func New(state statedomain.State, onSend func(string)) *Input {
	ta := textarea.New()
	ta.ShowLineNumbers = false
	ta.DynamicHeight = true

	style := ta.Styles()
	style.Focused.CursorLine = style.Focused.CursorLine.UnsetBackground().
		BorderLeft(false)
	style.Focused.Base = style.Focused.Base.BorderLeft(false)
	style.Focused.Text = style.Focused.Text.BorderLeft(false)
	style.Focused.CursorLineNumber = style.Focused.CursorLineNumber.BorderLeft(false)
	ta.SetStyles(style)

	return &Input{
		state:    state,
		onSend:   onSend,
		textarea: ta,
	}
}

func (i *Input) Init() tea.Cmd {
	return textarea.Blink
}

func (i *Input) Update(msg tea.Msg) (*Input, tea.Cmd) {
	if !i.textarea.Focused() {
		i.textarea.Focus()
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i.onSend(i.Value())
			i.textarea.SetValue("")
			return i, nil
		case "shift+enter", "alt+enter":
			i.textarea.InsertString("\n")
		}
	}

	i.textarea, cmd = i.textarea.Update(msg)

	return i, cmd
}

func (i *Input) View(width, height int) string {
	i.textarea.SetWidth(width)
	return lipgloss.NewStyle().
		Width(width).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00FF88")).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		Render(i.textarea.View())
}

func (i *Input) Value() string {
	value := i.textarea.Value()
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "\n", "\\\n")
	return value
}
