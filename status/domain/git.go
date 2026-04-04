package statusdomain

type Git struct {
	BranchLocal  string
	BranchRemote string
	Ahead        int
	Behind       int
	Untracked    int
	Modified     int
	Staged       int
	Deleted      int
	Conflicts    int
}
