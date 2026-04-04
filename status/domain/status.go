package statusdomain

type Status struct {
	OS    string
	Shell string
	User  string
	Dir   string
	Git   *Git
}
