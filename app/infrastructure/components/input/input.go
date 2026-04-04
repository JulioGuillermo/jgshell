package input

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	controllerdomain "github.com/julioguillermo/jgshell/controller/domain"
	syntaxdomain "github.com/julioguillermo/jgshell/syntax/domain"
)

type Input struct {
	ctl         controllerdomain.ShellController
	textarea    textarea.Model
	onSend      func(string)
	highlighter syntaxdomain.Highlighter
	showCursor  bool
	lastInput   time.Time
}

func New(ctl controllerdomain.ShellController, onSend func(string), highlighter syntaxdomain.Highlighter) *Input {
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
		ctl:         ctl,
		onSend:      onSend,
		textarea:    ta,
		highlighter: highlighter,
	}
}

func (i *Input) Init() tea.Cmd {
	return doBlink()
}

func (i *Input) Update(msg tea.Msg) (*Input, tea.Cmd) {
	if !i.textarea.Focused() {
		i.textarea.Focus()
	}

	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case CursorBlink:
		if time.Since(i.lastInput) > time.Millisecond*500 {
			i.showCursor = !i.showCursor
			i.lastInput = time.Now()
		}
		cmds = append(cmds, doBlink())
	case tea.KeyMsg:
		i.lastInput = time.Now()
		i.showCursor = true

		switch msg.String() {
		case "enter":
			if i.Value() == "" {
				return i, nil
			}
			i.onSend(strings.ReplaceAll(i.Value(), "\\\n", "\n"))
			i.textarea.SetValue("")
			return i, nil
		case "shift+enter", "alt+enter":
			i.textarea.InsertString("\n")
		}
	}

	ta, cmd := i.textarea.Update(msg)
	i.textarea = ta
	cmds = append(cmds, cmd)

	return i, tea.Batch(cmds...)
}

func (i *Input) View(width, height int) string {
	i.textarea.SetWidth(width)
	return lipgloss.NewStyle().
		Width(width).
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00FF88")).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		Render(i.Render())
}

func (i *Input) Value() string {
	value := i.textarea.Value()
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "\n", "\\\n")
	return value
}

func (i *Input) SetValue(value string) {
	i.textarea.SetValue(value)
}

func (i *Input) Insert(value string) {
	i.textarea.InsertString(value)
}
