package log

import (
	"errors"
	"fmt"
	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/config"
)

type LogFactory struct {
}

func (self *LogFactory) NewConf() config.Configuration {
	return new(LogConf)
}

func (self *LogFactory) GetInstance(config config.Configuration) (interface{}, error) {
	logConf, ok := config.(*LogConf)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Configuration is not LogConf %v", config))
	}

	return New(logConf.File)
}

func init() {
	factory := new(LogFactory)
	morgoth.NotifierRegistery.RegisterFactory("log", factory)
}
