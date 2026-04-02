package stateapplication

import (
	"embed"
	"regexp"

	statedomain "github.com/julioguillermo/jgshell/state/domain"
)

//go:embed status.sh
var StatusScript embed.FS

var (
	reOS   = regexp.MustCompile(`(?m)OS: \((?P<os>.*)\)`)
	reUser = regexp.MustCompile(`(?m)User: \((?P<user>.*)\)`)
	reDir  = regexp.MustCompile(`(?m)Dir: \((?P<dir>.*)\)`)
)

func (s *Status) Load(cmd statedomain.FastCmd) {
	script, err := StatusScript.ReadFile("status.sh")
	if err != nil {
		return
	}

	output, _ := cmd.FastCmd(string(script))

	s.parseOS(output)
	s.parseUser(output)
	s.parseDir(output)
}

func (s *Status) parseOS(output string) {
	s.os = parse(reOS, output)
}

func (s *Status) parseUser(output string) {
	s.user = parse(reUser, output)
}

func (s *Status) parseDir(output string) {
	s.dir = parse(reDir, output)
}

func parse(re *regexp.Regexp, output string) string {
	match := re.FindStringSubmatch(output)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}
