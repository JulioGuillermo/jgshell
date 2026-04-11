package statusinfrastructure

import (
	"embed"
)

//go:embed scripts
var scripts embed.FS

func (s *StatusLoader) getScript(sh string) (string, error) {
	switch sh {
	case "powershell":
		return s.loadScript("status.ps1")
	case "fish":
		return s.loadScript("status.fish")
	default:
		return s.loadScript("status.sh")
	}
}

func (s *StatusLoader) loadScript(script string) (string, error) {
	output, err := scripts.ReadFile("scripts/" + script)
	return string(output), err
}
