package fitting

import (
	app "github.com/nathanielc/morgoth/app/types"
	"github.com/nathanielc/morgoth/registery"
)

type Fitting interface {
	Name() string
	Start(app.App)
	Stop()
}

var (
	Registery *registery.Registery
)

func init() {
	Registery = registery.New()
}
