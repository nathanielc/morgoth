package kstest

import (
	"errors"
	"fmt"
	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/config"
)

type KSTestFactory struct {
}

func (self *KSTestFactory) NewConf() config.Configuration {
	return new(KSTestConf)
}

func (self *KSTestFactory) GetInstance(config config.Configuration) (interface{}, error) {
	conf, ok := config.(*KSTestConf)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Configuration is not KSTestConf%v", config))
	}
	engine := &KSTest{
		confidence: conf.Confidence,
	}
	return engine, nil
}

func init() {
	factory := new(KSTestFactory)
	morgoth.FingerprinterRegistery.RegisterFactory("kstest", factory)
}
