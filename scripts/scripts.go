package scripts

import "embed"

//go:embed wrapper
var WrapperScript embed.FS

//go:embed autocomplete
var AutoCompleteScript embed.FS
