package mgof

import (
	"errors"
	"fmt"
	"github.com/nathanielc/morgoth/config/types"
	"github.com/nathanielc/morgoth/detector"
)

type MGOFFactory struct {
}

func (self *MGOFFactory) NewConf() types.Configuration {
	return new(MGOFConf)
}

func (self *MGOFFactory) GetInstance(config types.Configuration) (interface{}, error) {
	conf, ok := config.(*MGOFConf)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Configuration is not MGOFConf%v", config))
	}
	mgof := &MGOF{
		config: conf,
	}
	return mgof, nil
}

func init() {
	factory := new(MGOFFactory)
	detector.Registery.RegisterFactory("mgof", factory)
}
