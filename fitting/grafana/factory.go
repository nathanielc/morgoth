package grafana

import (
	"errors"
	"fmt"
	"github.com/nvcook42/morgoth/config/types"
	"github.com/nvcook42/morgoth/fitting"
)

type GrafanaFactory struct {
}

func (self *GrafanaFactory) NewConf() types.Configuration {
	return new(GrafanaConf)
}

func (self *GrafanaFactory) GetInstance(config types.Configuration) (interface{}, error) {
	grafanaConf, ok := config.(*GrafanaConf)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Configuration is not GrafanaConf %v", config))
	}
	return &GrafanaFitting{conf: *grafanaConf}, nil
}

func init() {
	factory := new(GrafanaFactory)
	fitting.Registery.RegisterFactory("grafana", factory)
}
