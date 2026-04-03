package menu

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	ItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8800ff"))
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff0088"))
)

type itemDelegate struct {
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	total := len(m.Items())
	total = len(fmt.Sprint(total))

	i, ok := listItem.(Item)
	if !ok {
		return
	}

	fStr := fmt.Sprintf("%%0%dd. %%s", total)
	str := fmt.Sprintf(fStr, index+1, i)

	output := ""
	if index == m.Index() {
		output = SelectedItemStyle.Render("> " + str)
	} else {
		output = ItemStyle.Render("  " + str)
	}

	fmt.Fprint(w, output)
}
