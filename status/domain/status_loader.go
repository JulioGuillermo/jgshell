package statusdomain

type StatusLoader interface {
	Load() (*Status, error)
}
