package detector

import (
	"github.com/nvcook42/morgoth/registery"
)

type Detector interface {
	Detect()
}

var (
	Registery *registery.Registery
)

func init() {
	Registery = registery.New()
}
