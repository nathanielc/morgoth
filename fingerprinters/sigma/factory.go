package sigma

import (
	"errors"
	"fmt"
	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/config"
)

type SigmaFactory struct {
}

func (self *SigmaFactory) NewConf() config.Configuration {
	return new(SigmaConf)
}

func (self *SigmaFactory) GetInstance(config config.Configuration) (interface{}, error) {
	conf, ok := config.(*SigmaConf)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Configuration is not SigmaConf%v", config))
	}
	engine := &Sigma{
		deviations: conf.Deviations,
	}
	return engine, nil
}

func init() {
	factory := new(SigmaFactory)
	morgoth.FingerprinterRegistery.RegisterFactory("sigma", factory)
}
