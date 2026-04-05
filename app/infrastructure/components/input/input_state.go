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
	pos := i.GetCurrentLinePosition()
	lines[cursorLine] = i.GetCompletionLine(text, line, pos)

	i.textarea.SetValue(strings.Join(lines, "\n"))
	i.textarea.SetCursorColumn(pos + len(text))
}

func (i *Input) GetCompletionLine(completion, line string, cursor int) string {
	start := line[:cursor]
	end := line[cursor:]

	dividerEnd := i.GetCompletionEndDivider(start, completion)
	dividerStart := i.GetCompletionStartDivider(end, completion)

	start = strings.TrimSuffix(start, completion[:dividerEnd])
	end = strings.TrimPrefix(end, completion[dividerStart:])

	return start + completion + end
}

func (i *Input) GetCompletionEndDivider(line, completion string) int {
	for i := len(completion); i >= 0; i-- {
		if strings.HasSuffix(line, completion[:i]) {
			return i
		}
	}
	return 0
}

func (i *Input) GetCompletionStartDivider(line, completion string) int {
	for i := 0; i < len(completion); i++ {
		if strings.HasPrefix(line, completion[i:]) {
			return i
		}
	}
	return len(completion)
}
