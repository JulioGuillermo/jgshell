package input

import (
	"strings"
	"unicode/utf8"
)

func (i *Input) Position() int {
	// textarea uses (row, col) internally for cursor, but we need absolute position for completion
	// Let's calculate it from the value and cursor position
	val := i.textarea.Value()
	cursorLine := i.textarea.Line()
	cursorCol := i.textarea.LineInfo().CharOffset

	lines := strings.Split(val, "\n")
	pos := 0
	for l := range cursorLine {
		pos += utf8.RuneCountInString(lines[l]) + 1 // +1 for newline
	}
	pos += cursorCol
	return pos
}

func (i *Input) GetCurrentLine() string {
	val := i.textarea.Value()
	cursorLine := i.textarea.Line()
	lines := strings.Split(val, "\n")
	return lines[cursorLine]
}

func (i *Input) GetCurrentLinePosition() int {
	return i.textarea.LineInfo().CharOffset
}
