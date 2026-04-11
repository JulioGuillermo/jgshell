package routerdomain

type Element interface {
	String() string
	FinalString() string
	IsEnded() bool
}
