package engine

import (
	"github.com/nvcook42/morgoth/registery"
)

type Engine interface {
	GetReader() Reader
	GetWriter() Writer
}

var (
	Registery *registery.Registery
)

func init() {
	Registery = registery.New()
}
