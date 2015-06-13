package config

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	config "github.com/nathanielc/morgoth/config/types"
)

type MorgothConf struct {
	MetaDir string `yaml:"meta_dir" validate:"nonzero" default:"meta"`
}

func (self MorgothConf) Validate() error {
	return validator.Validate(self)
}

func (self *MorgothConf) Default() {
	config.PerformDefault(self)
}
