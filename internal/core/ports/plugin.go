package ports

type Plugin struct {
	Name        string
	Description string
	Version     string
	Path        string
}

type PluginManager interface {
	ListPlugins() ([]Plugin, error)
	InstallPlugin(nameOrUrl string) error
	RemovePlugin(name string) error
	ExecutePlugin(name string, args []string) error
}
