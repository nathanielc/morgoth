package plugin

type Plugin func(string) interface{}

type Loader interface {
	GetPlugin(name, conf string) interface{}
}
type PluginLoader struct {
	plugins map[string]Plugin
}

func (self *PluginLoader) RegisterPlugin(name string, plugin Plugin) {
	self.plugins[name] = plugin
}

func (self *PluginLoader) GetPlugin(name, conf string) interface{} {
	return self.plugins[name](conf)
}
