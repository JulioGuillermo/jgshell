package syntaxdomain

type Highlighter interface {
	Highlight(text string) string
}
