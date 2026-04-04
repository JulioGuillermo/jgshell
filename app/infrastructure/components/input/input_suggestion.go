package input

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func (i *Input) GetSuggestion() string {
	value := i.textarea.Value()
	if value == "" {
		return ""
	}
	return i.ctl.FilterLast(value)
}

func (i *Input) GetRenderSuggestion() string {
	if i.suggestion == "" {
		return ""
	}

	value := i.textarea.Value()
	suggestion := strings.TrimPrefix(i.suggestion, value)
	suggestion = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#555555")).
		Render(suggestion)
	return suggestion
}

func (i *Input) ApplySuggestion() {
	if i.suggestion == "" {
		return
	}

	value := i.textarea.Value()
	if strings.HasPrefix(value, i.suggestion) {
		return
	}

	i.SetValue(i.suggestion)
	i.suggestion = ""
}
