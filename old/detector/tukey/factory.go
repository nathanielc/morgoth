package tukey

import (
	"errors"
	"fmt"
	"github.com/nathanielc/morgoth/config/types"
	"github.com/nathanielc/morgoth/detector"
)

type TukeyFactory struct {
}

func (self *TukeyFactory) NewConf() types.Configuration {
	return new(TukeyConf)
}

func (self *TukeyFactory) GetInstance(config types.Configuration) (interface{}, error) {
	conf, ok := config.(*TukeyConf)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Configuration is not TukeyConf%v", config))
	}
	tukey := &Tukey{
		threshold: conf.Threshold,
	}
	return tukey, nil
}

func init() {
	factory := new(TukeyFactory)
	detector.Registery.RegisterFactory("tukey", factory)
}
