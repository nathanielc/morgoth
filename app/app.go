package app

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/config"
	"github.com/nvcook42/morgoth/engine"
	_ "github.com/nvcook42/morgoth/plugins"
)

type App struct {
	config *config.Config
	Engine engine.Engine
}

func New(config *config.Config) *App {
	app := App{config: config}
	return &app
}

func (self *App) Run() error {

	log.Info("Setup engine")
	eng, err := self.config.EngineConf.GetEngine()
	if err != nil {
		return err
	}
	self.Engine = eng
	err = self.Engine.Initialize()
	if err != nil {
		return err
	}

	log.Info("Setup metrics manager")

	log.Info("Starting all fittings")

	log.Info("All fittings have finished. Exiting")

	return nil
}
