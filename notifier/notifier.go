package notifier

import (
	"github.com/nvcook42/morgoth/registery"
)

type Notifier interface {
	Notify()
}

var (
	Registery *registery.Registery
)

func init() {
	Registery = registery.New()
}
