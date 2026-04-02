package tools

import "strings"

func cleanLineRet(line string) string {
	components := strings.Split(line, "\r")
	final := ""
	for _, c := range components {
		if len(final) <= len(c) {
			final = c
			continue
		}
		final = c + final[len(c):]
	}
	return final
}

func cleanLineRunes(line string) string {
	runes := []rune(line)
	clean := make([]rune, 0, len(runes))
	for i := range len(runes) {
		if runes[i] == '\b' {
			clean = clean[:len(clean)-1]
			continue
		}
		if runes[i] == '\t' {
			clean = append(clean, []rune("    ")...)
			continue
		}
		clean = append(clean, runes[i])
	}
	return string(clean)
}

func CleanText(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		line = cleanLineRet(line)
		line = cleanLineRunes(line)
		lines[i] = line
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}
