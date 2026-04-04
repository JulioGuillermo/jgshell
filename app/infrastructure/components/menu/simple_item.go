package menu

type SimpleItem string

func (i SimpleItem) FilterValue() string { return string(i) }
