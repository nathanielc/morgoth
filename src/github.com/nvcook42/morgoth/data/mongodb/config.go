package mongodb

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/defaults"
	"gopkg.in/validator.v2"
)

type MongoDBConf struct {
	Host      string `validate:"min=1"           default:"localhost"`
	Port      uint   `validate:"min=1,max=65535" default:"27017"`
	Database  string `validate:"min=1"`
	IsSharded bool   `default:"false"`
}

func (self *MongoDBConf) Validate() error {
	return validator.Validate(self)
}

func (self *MongoDBConf) Default() {
	err := self.Validate()
	if err != nil {
		errs := err.(validator.ErrorMap)
		for fieldName := range errs {
			switch fieldName {
			case "Host", "Port":
				log.Infof("Using default for MongoDBConf.%s", fieldName)
				defaults.SetDefault(self, fieldName)
			}
		}

	}
}
