package mgof

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/defaults"
	"gopkg.in/validator.v2"
)

type MGOFConf struct {
	CHI float64 `default:"0.5"`
}

func (self *MGOFConf) Default() {
	err := self.Validate()
	if err != nil {
		errs := err.(validator.ErrorMap)
		for fieldName := range errs {
			log.Infof("Using default for MGOFConf.%s", fieldName)
			defaults.SetDefault(self, fieldName)
		}
	}

}

func (self *MGOFConf) Validate() error {
	return validator.Validate(self)
}
