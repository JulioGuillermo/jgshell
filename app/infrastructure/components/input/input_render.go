package input

import (
	"strings"
	"unicode/utf8"

	"charm.land/lipgloss/v2"
)

// View renders the input field with a modern boxed design.
func (i *Input) Render(suggestion string) string {
	// Render content with cursor and highlighting
	var inputContent string

	highlighted := i.highlighter.Highlight(i.textarea.Value()) + suggestion

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
	var lastAnsi string

	// Special case: cursor at the end
	textLen := utf8.RuneCountInString(i.textarea.Value())

	runes := []rune(highlighted)
	for _, r := range runes {
		// Track ANSI escape sequences
		if r == '\x1b' {
			inAnsi = true
			lastAnsi = ""
		}

		if !inAnsi {
			if charCount == pos {
				// Insert styled cursor. We wrap the next character if it exists.
				// Since we are in a highlighted string, we might want to keep the color
				// but change the background.
				out.WriteString(cursorStyle.Render(string(r)))
				out.WriteString(lastAnsi)
				cursorInserted = true
				charCount++
				continue
			}
			charCount++
		} else {
			lastAnsi += string(r)
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
