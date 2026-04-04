package toolsapplication

import "strings"

type OutputCleaner struct{}

func NewOutputCleaner() *OutputCleaner {
	return &OutputCleaner{}
}

func (o *OutputCleaner) Clean(output string) string {
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		line = o.cleanLineRet(line)
		line = o.cleanLineRunes(line)
		lines[i] = line
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func (o *OutputCleaner) cleanLineRet(line string) string {
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

func (o *OutputCleaner) cleanLineRunes(line string) string {
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
