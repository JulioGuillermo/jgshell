package persistencedomain

type Persistence interface {
	SaveHistory(history []string) error
	LoadHistory() ([]string, error)
}
