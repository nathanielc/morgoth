package config

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	config "github.com/nvcook42/morgoth/config/types"
)

type MorgothConf struct {
	MetaDir string      `yaml:"meta_dir" validate:"nonzero" default:"meta"`
}

func (self MorgothConf) Validate() error {
	return validator.Validate(self)
}

func (self *MorgothConf) Default() {
	config.PerformDefault(self)
}
