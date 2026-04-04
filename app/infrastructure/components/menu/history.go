package menu

import (
	"slices"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/paginator"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

type History struct {
	list     list.Model
	items    []SimpleItem
	OnSelect func(string)
	OnClose  func()
}

func NewHistory() *History {
	l := list.New([]list.Item{}, SimpleItemDelegate{}, 20, 10)
	l.Title = "History"
	l.Styles.Title = l.Styles.Title.UnsetBackground().
		Foreground(lipgloss.Color("#00ff88")).
		Padding(0)
	l.SetFilteringEnabled(true)
	// l.SetFilteringEnabled(false)
	l.SetShowFilter(true)
	l.SetShowTitle(true)
	l.SetShowPagination(false)
	l.SetShowHelp(true)
	l.SetShowStatusBar(false)
	l.DisableQuitKeybindings()

	return &History{
		list: l,
	}
}

func (h *History) SetItems(items []string) {
	slices.Reverse(items)
	h.items = make([]SimpleItem, len(items))
	for i, item := range items {
		h.items[i] = SimpleItem(item)
	}
	listItems := make([]list.Item, len(h.items))
	for i, item := range h.items {
		listItems[i] = item
	}
	h.list.SetItems(listItems)
	h.list.SetFilterText("")
}

func (h *History) Init() tea.Cmd {
	return nil
}

func (h *History) Update(msg tea.Msg, width int, height int) (*History, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		h.OnKey(msg.String())
	}

	l, c := h.list.Update(msg)
	h.list = l
	h.list.SetWidth(width)
	h.list.SetHeight(max(min(height-5, 20), 3))

	return h, c
}

func (h *History) OnKey(key string) {
	switch key {
	case "esc", "space", "backspace":
		if h.list.SettingFilter() {
			return
		}
		if h.OnClose != nil {
			h.OnClose()
		}
	case "enter":
		h.onSelect()
	case "tab":
		h.list.CursorDown()
	case "shift+tab":
		h.list.CursorUp()
	}
}

func (h *History) onSelect() {
	if len(h.items) == 0 {
		return
	}
	if h.list.SettingFilter() {
		return
	}
	idx := h.list.Index()
	if h.list.FilterValue() != "" {
		idx = h.list.GlobalIndex()
	}
	selected := h.items[idx]
	if h.OnSelect != nil {
		h.OnSelect(string(selected))
	}
	if h.OnClose != nil {
		h.OnClose()
	}
}

func (h *History) View(width, height int) string {
	if len(h.items) == 0 {
		return ""
	}
	elements := []string{h.list.View()}
	if pagination := h.paginationView(width - 2); pagination != "" {
		elements = append(elements, pagination)
	}

	return lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00ff88")).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				elements...,
			),
		)
}

func (h *History) paginationView(width int) string {
	if h.list.Paginator.TotalPages < 2 {
		return ""
	}

	s := h.list.Paginator.View()

	if ansi.StringWidth(s) > width {
		h.list.Paginator.Type = paginator.Arabic
		s = h.list.Styles.ArabicPagination.Render(h.list.Paginator.View())
	}

	style := h.list.Styles.PaginationStyle

	return style.Render(s)
}
