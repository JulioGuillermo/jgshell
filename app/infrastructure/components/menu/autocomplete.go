package menu

import (
	"slices"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/paginator"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

type Autocomplete struct {
	list     list.Model
	items    []Item
	OnSelect func(string)
	OnClose  func()
}

func NewAutocomplete() *Autocomplete {
	l := list.New([]list.Item{}, itemDelegate{}, 20, 10)
	// l.Title = "What do you want for dinner?"
	l.SetFilteringEnabled(true)
	// l.SetFilteringEnabled(false)
	l.SetShowFilter(true)
	l.SetShowTitle(false)
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.DisableQuitKeybindings()

	return &Autocomplete{
		list: l,
	}
}

func (a *Autocomplete) SetItems(items []string) {
	slices.Sort(items)
	a.items = make([]Item, len(items))
	for i, item := range items {
		a.items[i] = Item(item)
	}
	listItems := make([]list.Item, len(a.items))
	for i, item := range a.items {
		listItems[i] = item
	}
	a.list.SetItems(listItems)
	a.list.SetFilterText("")
}

func (a *Autocomplete) Init() tea.Cmd {
	return nil
}

func (a *Autocomplete) Update(msg tea.Msg, width int) (*Autocomplete, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		a.OnKey(msg.String())
	}

	l, c := a.list.Update(msg)
	a.list = l
	a.list.SetWidth(width)

	return a, c
}

func (a *Autocomplete) OnKey(key string) {
	switch key {
	case "esc", "space", "backspace":
		if a.list.SettingFilter() {
			return
		}
		if a.OnClose != nil {
			a.OnClose()
		}
	case "enter":
		a.onSelect()
	case "tab":
		a.list.CursorDown()
	case "shift+tab":
		a.list.CursorUp()
	}
}

func (a *Autocomplete) onSelect() {
	if len(a.items) == 0 {
		return
	}
	if a.list.SettingFilter() {
		return
	}
	idx := a.list.Index()
	if a.list.FilterValue() != "" {
		idx = a.list.GlobalIndex()
	}
	selected := a.items[idx]
	if a.OnSelect != nil {
		a.OnSelect(string(selected))
	}
	if a.OnClose != nil {
		a.OnClose()
	}
}

func (a *Autocomplete) View(width, height int) string {
	if len(a.items) == 0 {
		return ""
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		a.list.View(),
		a.paginationView(width),
	)
}

func (a *Autocomplete) paginationView(width int) string {
	if a.list.Paginator.TotalPages < 2 {
		return ""
	}

	s := a.list.Paginator.View()

	if ansi.StringWidth(s) > width {
		a.list.Paginator.Type = paginator.Arabic
		s = a.list.Styles.ArabicPagination.Render(a.list.Paginator.View())
	}

	style := a.list.Styles.PaginationStyle

	return style.Render(s)
}
