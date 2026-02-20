package ports

type User struct {
	Username string
	UID      string
	GID      string
	Info     string
	HomeDir  string
	Shell    string
}

type Group struct {
	Groupname string
	GID       string
	Members   []string
}

type UserManager interface {
	ListUsers() ([]User, error)
	ListGroups() ([]Group, error)

	// Management
	AddUser(username string, options UserOptions) error
	DeleteUser(username string, removeHome bool) error
	ModifyUser(username string, options UserOptions) error

	// Group Management
	AddGroup(groupname string, gid string) error
	DeleteGroup(groupname string) error
}

type UserOptions struct {
	UID     string
	GID     string
	Groups  []string
	Shell   string
	HomeDir string
	Comment string
}
