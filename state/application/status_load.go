package stateapplication

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/julioguillermo/jgshell/scripts"
	statedomain "github.com/julioguillermo/jgshell/state/domain"
)

var (
	reOS   = regexp.MustCompile(`(?m)OS: \((?P<os>.*)\)`)
	reUser = regexp.MustCompile(`(?m)User: \((?P<user>.*)\)`)
	reDir  = regexp.MustCompile(`(?m)Dir: \((?P<dir>.*)\)`)
)

func (s *Status) Load(cmd statedomain.FastCmd) {
	s.shell = cmd.GetShell()

	script, err := scripts.StatusScript.ReadFile("status/status.sh")
	if err != nil {
		return
	}

	output, _ := cmd.FastCmdClean(string(script))

	s.parseOS(output)
	s.parseUser(output)
	s.parseDir(output)
	s.parseGit(output)
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

func (s *Status) parseGit(output string) {
	s.git = nil

	re := regexp.MustCompile(`(?s)=== GIT START ===(.*?)=== GIT END ===`)
	match := re.FindStringSubmatch(output)
	if match == nil {
		return
	}
	gitContent := match[1]
	if strings.Contains(gitContent, "NO GIT") {
		return
	}

	s.git = &statedomain.GitStats{}

	lines := strings.SplitSeq(gitContent, "\n")
	for line := range lines {
		if strings.HasPrefix(line, "##") {
			s.parseGitBranchLine(line)
		} else {
			s.parseGitStatsLine(line)
		}
	}
}

func (s *Status) parseGitBranchLine(line string) {
	if line == "" {
		return
	}

	reAhead := regexp.MustCompile(`\[ahead (\d+)\]`)
	reBehind := regexp.MustCompile(`\[behind (\d+)\]`)

	if m := reAhead.FindStringSubmatch(line); m != nil {
		fmt.Sscanf(m[1], "%d", &s.git.Ahead)
	}
	if m := reBehind.FindStringSubmatch(line); m != nil {
		fmt.Sscanf(m[1], "%d", &s.git.Behind)
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
			s.git.BranchLocal = match[i]
		case "remote":
			s.git.BranchRemote = match[i]
		}
	}
}

func (s *Status) parseGitStatsLine(line string) {
	if len(line) < 2 {
		return
	}

	prefix := line[:2]

	switch prefix {
	case "??":
		s.git.Untracked++
	case " M", "AM", "MM": // Modificados no stashed
		s.git.Modified++
	case "M ", "A ", "D ": // Staged (Modificados, Agregados, Eliminados)
		s.git.Staged++
	case " D":
		s.git.Deleted++
	case "UU", "AA", "DD", "AU", "UA", "UD", "DU":
		s.git.Conflicts++
	}
}
