package statusinfrastructure

import (
	"embed"
)

//go:embed scripts
var scripts embed.FS

func (s *StatusLoader) getScript(sh string) (string, error) {
	if sh == "powershell" {
		script, err := scripts.ReadFile("scripts/status.ps1")
		if err != nil {
			return "", err
		}
		return string(script), nil
	}

	script, err := scripts.ReadFile("scripts/status.sh")
	if err != nil {
		return "", err
	}
	return string(script), nil
}
