package input

import (
	"strings"
	"unicode/utf8"

	"charm.land/lipgloss/v2"
)

func (i *Input) Position() int {
	// textarea uses (row, col) internally for cursor, but we need absolute position for completion
	// Let's calculate it from the value and cursor position
	val := i.textarea.Value()
	cursorLine := i.textarea.Line()
	cursorCol := i.textarea.LineInfo().CharOffset

	lines := strings.Split(val, "\n")
	pos := 0
	for l := 0; l < cursorLine; l++ {
		pos += utf8.RuneCountInString(lines[l]) + 1 // +1 for newline
	}
	pos += cursorCol
	return pos
}

// View renders the input field with a modern boxed design.
func (i *Input) Render() string {
	// Render content with cursor and highlighting
	var inputContent string

	highlighted := i.highlighter.Highlight(i.textarea.Value())

	if i.showCursor {
		pos := i.Position()
		inputContent = i.renderWithCursor(highlighted, pos)
	} else {
		inputContent = highlighted
	}

	return inputContent
}

func (i *Input) renderWithCursor(highlighted string, pos int) string {
	cursorStyle := lipgloss.NewStyle().Reverse(true)

	// Simple ANSI-aware character counting
	var out strings.Builder
	var charCount int
	var inAnsi bool
	var cursorInserted bool

	// Special case: cursor at the end
	textLen := utf8.RuneCountInString(i.textarea.Value())

	runes := []rune(highlighted)
	for _, r := range runes {
		// Track ANSI escape sequences
		if r == '\x1b' {
			inAnsi = true
		}

		if !inAnsi {
			if charCount == pos {
				// Insert styled cursor. We wrap the next character if it exists.
				// Since we are in a highlighted string, we might want to keep the color
				// but change the background.
				out.WriteString(cursorStyle.Render(string(r)))
				cursorInserted = true
				charCount++
				continue
			}
			charCount++
		}

		out.WriteRune(r)

		if inAnsi && r == 'm' {
			inAnsi = false
		}
	}

	// If cursor was at the end, append it
	if !cursorInserted && pos >= textLen {
		out.WriteString(cursorStyle.Render(" "))
	}

	return out.String()
}
