package data

import ()

type EngineLoader struct {
}

func (self *EngineLoader) GetEngine(name, conf string) *Engine {
	return self.GetPlugin(name, conf).(*Engine)
}
