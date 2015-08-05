package jsdiv

import (
	"errors"
	"fmt"
	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/config"
)

type JSDivFactory struct {
}

func (self *JSDivFactory) NewConf() config.Configuration {
	return new(JSDivConf)
}

func (self *JSDivFactory) GetInstance(config config.Configuration) (interface{}, error) {
	conf, ok := config.(*JSDivConf)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Configuration is not JSDivConf%v", config))
	}
	engine := &JSDiv{
		min:    conf.Min,
		max:    conf.Max,
		nBins:  conf.NBins,
		pValue: conf.PValue,
	}
	return engine, nil
}

func init() {
	factory := new(JSDivFactory)
	morgoth.FingerprinterRegistery.RegisterFactory("jsdiv", factory)
}
