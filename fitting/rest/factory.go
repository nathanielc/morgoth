package rest

import (
	"github.com/nvcook42/morgoth/config/types"
	"github.com/nvcook42/morgoth/fitting"
	"errors"
	"fmt"
)

type RESTFactory struct {
}

func (self *RESTFactory) NewConf() types.Configuration {
	return new(RESTConf)
}

func (self *RESTFactory) GetInstance(config types.Configuration) (interface{}, error) {
	restConf, ok := config.(*RESTConf)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Configuration is not RESTConf %v", config))
	}
	return &RESTFitting{port:restConf.Port}, nil
}

func init() {
	factory := new(RESTFactory)
	fitting.Registery.RegisterFactory("rest", factory)
}
