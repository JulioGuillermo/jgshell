package input

import (
	"strings"
	"unicode/utf8"

	"github.com/julioguillermo/jgshell/tools"
)

func (i *Input) Position() int {
	// textarea uses (row, col) internally for cursor, but we need absolute position for completion
	// Let's calculate it from the value and cursor position
	val := i.textarea.Value()
	cursorLine := i.textarea.Line()
	cursorCol := i.GetCurrentLinePosition()

	lines := strings.Split(val, "\n")
	pos := 0
	for l := range cursorLine {
		pos += utf8.RuneCountInString(lines[l]) + 1 // +1 for newline
	}
	pos += cursorCol
	return pos
}

func (i *Input) GetCurrentLineRow() int {
	return i.textarea.Line()
}

func (i *Input) GetLinesCount() int {
	return i.textarea.LineCount()
}

func (i *Input) GetCurrentLine() string {
	val := i.textarea.Value()
	cursorLine := i.textarea.Line()
	lines := strings.Split(val, "\n")
	return lines[cursorLine]
}

func (i *Input) GetCurrentLinePosition() int {
	lineInfo := i.textarea.LineInfo()
	return lineInfo.StartColumn + lineInfo.CharOffset
}

func (i *Input) InsertAutocomplete(text string) {
	val := i.textarea.Value()
	cursorLine := i.textarea.Line()
	lines := strings.Split(val, "\n")
	line := lines[cursorLine]

	start := i.GetCurrentLinePosition()
	end := start
	for start-1 >= 0 && start-1 < len(line) && tools.IsAlphaNumeric(line[start-1]) {
		start--
	}
	for end >= 0 && end < len(line) && tools.IsAlphaNumeric(line[end]) {
		end++
	}

	start = max(0, min(start, len(line)))
	end = max(0, min(end, len(line)))

	lines[cursorLine] = line[:start] + text + line[end:]

	i.textarea.SetValue(strings.Join(lines, "\n"))
	i.textarea.SetCursorColumn(start + len(text))
}
