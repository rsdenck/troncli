package ports

import "context"

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
	ExecutePlugin(ctx context.Context, name string, args ...string) error
}
