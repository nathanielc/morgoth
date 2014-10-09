package app

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/config"
	"github.com/nvcook42/morgoth/engine"
	mtypes "github.com/nvcook42/morgoth/metric/types"
	"github.com/nvcook42/morgoth/metric"
	_ "github.com/nvcook42/morgoth/plugins"
)

type App struct {
	config *config.Config
	engine engine.Engine
	manager mtypes.Manager
}

func New(config *config.Config) *App {
	app := App{config: config}
	return &app
}

func (self *App) GetReader() engine.Reader {
	return self.engine.GetReader()
}

func (self *App) GetWriter() engine.Writer {
	return self.engine.GetWriter()
}

func (self *App) GetManager() mtypes.Manager {
	return self.manager
}

func (self *App) Run() error {

	log.Info("Setup engine")
	eng, err := self.config.EngineConf.GetEngine()
	if err != nil {
		return err
	}
	self.engine = eng
	err = self.engine.Initialize()
	if err != nil {
		return err
	}

	log.Info("Setup metrics manager")
	self.manager = metric.NewManager(self.config.Metrics, self)

	log.Info("Starting all fittings")

	log.Info("All fittings have finished. Exiting")

	return nil
}
