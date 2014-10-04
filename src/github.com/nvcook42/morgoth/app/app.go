package app

import (
	//"errors"
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/config"
	"github.com/nvcook42/morgoth/engine"
	//"gopkg.in/yaml.v2"
	_ "github.com/nvcook42/morgoth/engine/list"
)

type App struct {
	config *config.Config
	Engine engine.Engine
}

func New(config *config.Config) (*App, error) {
	app := App{config:config}
	eng, err := app.config.DataEngine.GetEngine()
	if err != nil {
		return nil, err
	}
	app.Engine = eng

	return &app, nil
}

func (self *App) Run() error {
	log.Info("Setup data engine")
	log.Infof("Engine: %v", self.Engine)
	return nil
}
