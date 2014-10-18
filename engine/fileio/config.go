package fileio

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/defaults"
	"gopkg.in/validator.v2"
)

type FileIOConf struct {
	Dir string `validate:"min=1" default:"./fileiodb/"`
}

func (self *FileIOConf) Validate() error {
	return validator.Validate(self)
}

//Sets any invalid fields to their default value
func (self *FileIOConf) Default() {
	err := self.Validate()
	if err != nil {
		errs := err.(validator.ErrorMap)
		for fieldName := range errs {
			if ok, _ := defaults.HasDefault(self, fieldName); ok {
				log.Infof("Using default for fileioConf.%s", fieldName)
				defaults.SetDefault(self, fieldName)
			}
		}
	}
}
