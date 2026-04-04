package autocompletedomain

type Autocomplete interface {
	GetAutocomplete(line string, cursor int) ([]string, error)
}
