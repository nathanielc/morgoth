package fitting

import (
	app "github.com/nvcook42/morgoth/app/types"
	"github.com/nvcook42/morgoth/registery"
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
