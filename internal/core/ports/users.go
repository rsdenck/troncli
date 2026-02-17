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
}
