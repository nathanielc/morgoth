package fitting

import (
	"github.com/nvcook42/morgoth/registery"
	app "github.com/nvcook42/morgoth/app/types"
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
