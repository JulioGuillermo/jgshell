package statusinfrastructure

import (
	"regexp"

	statusdomain "github.com/julioguillermo/jgshell/status/domain"
)

var (
	reOS   = regexp.MustCompile(`(?m)OS: \((?P<os>.*)\)`)
	reUser = regexp.MustCompile(`(?m)User: \((?P<user>.*)\)`)
	reDir  = regexp.MustCompile(`(?m)Dir: \((?P<dir>.*)\)`)
)

func (s *StatusLoader) parseOS(result *statusdomain.Status, output string) {
	result.OS = parse(reOS, output)
}

func (s *StatusLoader) parseUser(result *statusdomain.Status, output string) {
	result.User = parse(reUser, output)
}

func (s *StatusLoader) parseDir(result *statusdomain.Status, output string) {
	result.Dir = parse(reDir, output)
}

func parse(re *regexp.Regexp, output string) string {
	match := re.FindStringSubmatch(output)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}
