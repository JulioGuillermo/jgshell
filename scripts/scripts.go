package scripts

import "embed"

//go:embed shell
var ShellScript embed.FS

//go:embed status
var StatusScript embed.FS

//go:embed autocomplete
var AutoCompleteScript embed.FS
