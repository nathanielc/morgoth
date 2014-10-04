
package fitting

import (
	"github.com/nvcook42/morgoth/registery"
)

type Fitting interface {
	Start()
	Stop()
}

var (
	Registery *registery.Registery
)

func init() {
	Registery = registery.New()
}
