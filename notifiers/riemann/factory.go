package riemann

import (
	"errors"
	"fmt"
	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/config"
)

type RiemannFactory struct {
}

func (self *RiemannFactory) NewConf() config.Configuration {
	return new(RiemannConf)
}

func (self *RiemannFactory) GetInstance(config config.Configuration) (interface{}, error) {
	riemannConf, ok := config.(*RiemannConf)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Configuration is not RiemannConf %v", config))
	}

	return New(riemannConf.Host, riemannConf.Port)
}

func init() {
	factory := new(RiemannFactory)
	morgoth.NotifierRegistery.RegisterFactory("riemann", factory)
}
