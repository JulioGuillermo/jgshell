package persistencedomain

type PersistenceCtl interface {
	Push(cmds string) error
	Get() []string
	Filter(start string) []string
	FilterLast(start string) string
}
