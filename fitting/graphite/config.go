package graphite

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/defaults"
	"gopkg.in/validator.v2"
)

type GraphiteConf struct {
	Port uint `validate:"min=1,max=65535" default:"2003"`
}

func (self *GraphiteConf) Validate() error {
	return validator.Validate(self)
}

//Sets any invalid fields to their default value
func (self *GraphiteConf) Default() {
	err := self.Validate()
	if err != nil {
		errs := err.(validator.ErrorMap)
		for fieldName := range errs {
			if ok, _ := defaults.HasDefault(self, fieldName); ok {
				log.Infof("Using default for GraphiteConf.%s", fieldName)
				defaults.SetDefault(self, fieldName)
			}
		}
	}
}
