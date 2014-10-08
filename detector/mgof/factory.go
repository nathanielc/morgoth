package mgof

import (
	"github.com/nvcook42/morgoth/config/types"
	"github.com/nvcook42/morgoth/detector"
)

type MGOFFactory struct {
}

func (self *MGOFFactory) NewConf() types.Configuration {
	return new(MGOFConf)
}

func (self *MGOFFactory) GetInstance(config types.Configuration) (interface{}, error) {
	return new(MGOF), nil
}

func init() {
	factory := new(MGOFFactory)
	detector.Registery.RegisterFactory("mgof", factory)
}
