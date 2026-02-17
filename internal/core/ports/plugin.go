package ports

type Plugin struct {
	Name        string
	Description string
	Version     string
	Path        string
}

type PluginManager interface {
	ListPlugins() ([]Plugin, error)
	InstallPlugin(urlOrPath string) error
	RemovePlugin(name string) error
}
