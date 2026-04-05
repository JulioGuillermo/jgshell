package input

func (i *Input) HistoryClear() {
	i.historyIdx = -1
	i.original = ""
}

func (i *Input) HistoryUp() {
	if i.GetCurrentLineRow() > 0 {
		return
	}
	if i.historyIdx == -1 {
		i.original = i.Value()
		i.historyItems = i.ctl.Filter(i.original)
		i.historyIdx = len(i.historyItems) - 1
	} else if i.historyIdx > 0 {
		i.historyIdx--
	}
	if i.historyIdx >= 0 {
		i.SetValue(i.historyItems[i.historyIdx])
	}
}

func (i *Input) HistoryDown() {
	if i.GetCurrentLineRow() < i.GetLinesCount()-1 {
		return
	}
	if i.historyIdx == -1 {
		return
	}
	i.historyIdx++
	if i.historyIdx >= len(i.historyItems) {
		i.historyIdx = -1
	}
	if i.historyIdx == -1 {
		i.SetValue(i.original)
	} else {
		i.SetValue(i.historyItems[i.historyIdx])
	}
}
