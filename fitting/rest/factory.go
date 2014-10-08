package rest

import (
	"github.com/nvcook42/morgoth/config/types"
	"github.com/nvcook42/morgoth/fitting"
)

type RESTFactory struct {
}

func (self *RESTFactory) NewConf() types.Configuration {
	return new(RESTConf)
}

func (self *RESTFactory) GetInstance(config types.Configuration) (interface{}, error) {
	return new(RESTFitting), nil
}

func init() {
	factory := new(RESTFactory)
	fitting.Registery.RegisterFactory("rest", factory)
}
