package log

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	"github.com/nathanielc/morgoth/config"
)

type LogConf struct {
	File string `validate:"min=1"`
}

func (self *LogConf) Validate() error {
	return validator.Validate(self)
}

//Sets any invalid fields to their default value
func (self *LogConf) Default() {
	config.PerformDefault(self)
}
