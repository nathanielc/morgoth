package kstest

import (
	"errors"
	"fmt"
	"github.com/nvcook42/morgoth/config/types"
	"github.com/nvcook42/morgoth/detector"
)

type KSTestFactory struct {
}

func (self *KSTestFactory) NewConf() types.Configuration {
	return new(KSTestConf)
}

func (self *KSTestFactory) GetInstance(config types.Configuration) (interface{}, error) {
	conf, ok := config.(*KSTestConf)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Configuration is not KSTestConf%v", config))
	}
	kstest := &KSTest{
		config: conf,
	}
	return kstest, nil
}

func init() {
	factory := new(KSTestFactory)
	detector.Registery.RegisterFactory("kstest", factory)
}
