package statedomain

type Status interface {
	Load(FastCmd)
	User() string
	Dir() string
	OS() string
	Shell() string
}
