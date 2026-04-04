package statusinfrastructure

import (
	"fmt"
	"regexp"
	"strings"

	statusdomain "github.com/julioguillermo/jgshell/status/domain"
)

func (s *StatusLoader) parseGit(result *statusdomain.Status, output string) {
	result.Git = nil

	re := regexp.MustCompile(`(?s)=== GIT START ===(.*?)=== GIT END ===`)
	match := re.FindStringSubmatch(output)
	if match == nil {
		return
	}
	gitContent := match[1]
	if strings.Contains(gitContent, "NO GIT") {
		return
	}

	result.Git = &statusdomain.Git{}

	lines := strings.SplitSeq(gitContent, "\n")
	for line := range lines {
		if strings.HasPrefix(line, "##") {
			s.parseGitBranchLine(result, line)
		} else {
			s.parseGitStatsLine(result, line)
		}
	}
}

func (s *StatusLoader) parseGitBranchLine(result *statusdomain.Status, line string) {
	if line == "" {
		return
	}

	reAhead := regexp.MustCompile(`\[ahead (\d+)\]`)
	reBehind := regexp.MustCompile(`\[behind (\d+)\]`)

	if m := reAhead.FindStringSubmatch(line); m != nil {
		fmt.Sscanf(m[1], "%d", &result.Git.Ahead)
	}
	if m := reBehind.FindStringSubmatch(line); m != nil {
		fmt.Sscanf(m[1], "%d", &result.Git.Behind)
	}

	re := regexp.MustCompile(`## (?P<local>[^\.\s]+)(?:\.{3}(?P<remote>[^\s\[]+))?`)
	match := re.FindStringSubmatch(line)
	if match == nil {
		return
	}

	for i, name := range re.SubexpNames() {
		if i >= len(match) {
			break
		}

		switch name {
		case "local":
			result.Git.BranchLocal = match[i]
		case "remote":
			result.Git.BranchRemote = match[i]
		}
	}
}

func (s *StatusLoader) parseGitStatsLine(result *statusdomain.Status, line string) {
	if len(line) < 2 {
		return
	}

	prefix := line[:2]

	switch prefix {
	case "??":
		result.Git.Untracked++
	case " M", "AM", "MM":
		result.Git.Modified++
	case "M ", "A ", "D ":
		result.Git.Staged++
	case " D":
		result.Git.Deleted++
	case "UU", "AA", "DD", "AU", "UA", "UD", "DU":
		result.Git.Conflicts++
	}
}
