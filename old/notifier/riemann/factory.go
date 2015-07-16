package riemann

import (
	"errors"
	"fmt"
	"github.com/nathanielc/morgoth/config/types"
	"github.com/nathanielc/morgoth/notifier"
)

type RiemannFactory struct {
}

func (self *RiemannFactory) NewConf() types.Configuration {
	return new(RiemannConf)
}

func (self *RiemannFactory) GetInstance(config types.Configuration) (interface{}, error) {
	riemannConf, ok := config.(*RiemannConf)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Configuration is not RiemannConf %v", config))
	}

	return New(riemannConf.Host, riemannConf.Port)
}

func init() {
	factory := new(RiemannFactory)
	notifier.Registery.RegisterFactory("riemann", factory)
}
