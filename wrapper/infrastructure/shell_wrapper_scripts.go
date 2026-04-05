package wrapperinfrastructure

import "embed"

//go:embed scripts
var Scripts embed.FS

type WrapperConfig struct {
	Script string
	Loader string
}

var WrapperConfigs = map[string]WrapperConfig{
	"powershell": {
		Script: "scripts/wrapper_powershell.ps1",
		Loader: `$setupScript = @'
%s
'@
# Ejecuta el contenido del string en la sesión actual
Invoke-Expression $setupScript

Write-Output "Wrapper script executed successfully"
`,
	},
	"bash": {
		Script: "scripts/wrapper_bash.sh",
		Loader: "",
	},
	"zsh": {
		Script: "scripts/wrapper_zsh.zsh",
		Loader: "",
	},
	"fish": {
		Script: "scripts/wrapper_fish.fish",
		Loader: "",
	},
	"nu": {
		Script: "scripts/wrapper_nushell.nu",
		Loader: "",
	},
	"sh": {
		Script: "scripts/wrapper_sh.sh",
		Loader: "",
	},
}

func (s *ShellWrapper) getConfig(sh string) (*WrapperConfig, error) {
	config, ok := WrapperConfigs[sh]
	if !ok {
		config = WrapperConfigs["bash"]
	}

	if config.Loader == "" {
		config.Loader = `eval "$(cat << 'EOF'
%s
EOF
)"
printf "Wrapper script executed successfully\n"
`
	}

	script, err := Scripts.ReadFile(config.Script)
	if err != nil {
		return nil, err
	}
	config.Script = string(script)

	return &config, nil
}
