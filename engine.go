package morgoth

import (
	_ "github.com/nathanielc/morgoth/engines/list"
	"github.com/nathanielc/morgoth/registery"
)

type Engine interface {
	Initialize() error
	GetWindows(query Query) ([]*Window, error)
}

var Registery *registery.Registery

func init() {
	Registery = registery.New()
}
