package app

import (
	//"errors"
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/config"
	//"gopkg.in/yaml.v2"
)

type App struct {
	config *config.Config
}

func New(config *config.Config) (*App, error) {
	app := App{config}

	return &app, nil
}

func (self *App) Run() error {
	log.Info("Setup data engine")
	//dataEngine := self.config.Get("data_engine", nil).(yaml.MapSlice)
	//if dataEngine == nil {
	//	return errors.New("No data_engine section found in config")

	//}
	//engineType := dataEngine[0].Key.(string)
	//engineConf := dataEngine[0].Value.(yaml.MapSlice)
	//log.Printf("Configuring %v with conf: %v", engineType, engineConf)
	return nil
}
