package engine

import (
	"github.com/Sirupsen/logrus"
)

// All plugins/extensions have to implement this
type Plugin interface {
	// Name of the plugin
	Name() string

	// Called when the plugin should be ran (create ui elements and such)
	// Also given a logger that should be used for this plugin
	Run(*logrus.Entry)
}

// Plugins can implement this if they need to cleanup before application exit
type PluginStopper interface {
	Stop()
}

var (
	Plugins = make(map[string]Plugin)
	started bool
)

func RegisterPlugin(p Plugin) {
	if started {
		panic("Tried adding plugin after starting")
	}
	Plugins[p.Name()] = p
}

func FindPlugin(name string) Plugin {
	return Plugins[name]
}
