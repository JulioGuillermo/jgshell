package menu

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type Autocomplete struct {
	list  list.Model
	items []Item
}

func NewAutocomplete() *Autocomplete {
	l := list.New([]list.Item{}, itemDelegate{}, 20, 5)
	// l.Title = "What do you want for dinner?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return &Autocomplete{
		list: l,
	}
}

func (a *Autocomplete) SetItems(items []string) {
	a.items = make([]Item, len(items))
	for i, item := range items {
		a.items[i] = Item(item)
	}
	listItems := make([]list.Item, len(a.items))
	for i, item := range a.items {
		listItems[i] = item
	}
	a.list.SetItems(listItems)
}

func (a *Autocomplete) Init() tea.Cmd {
	return nil
}

func (a *Autocomplete) Update(msg tea.Msg, width int) (*Autocomplete, tea.Cmd) {
	l, c := a.list.Update(msg)
	a.list = l
	a.list.SetWidth(width)
	return a, c
}

func (a *Autocomplete) View(width, height int) string {
	if len(a.items) == 0 {
		return ""
	}

	return a.list.View()
}
