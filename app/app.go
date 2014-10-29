package app

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/config"
	"github.com/nvcook42/morgoth/engine"
	"github.com/nvcook42/morgoth/fitting"
	"github.com/nvcook42/morgoth/metric"
	mtypes "github.com/nvcook42/morgoth/metric/types"
	_ "github.com/nvcook42/morgoth/plugins"
	"github.com/nvcook42/morgoth/schedule"
	"os"
	"os/signal"
	"sync"
)

type App struct {
	manager  mtypes.Manager
	engine   engine.Engine
	config   *config.Config
	fittings []fitting.Fitting
	schedule schedule.Schedule
}

func New(config *config.Config) *App {
	app := App{
		config: config,
	}
	app.schedule = config.GetSchedule()
	return &app
}

func (self *App) GetReader() engine.Reader {
	return self.engine.GetReader()
}

func (self *App) GetSchedule() schedule.Schedule {
	return self.schedule
}

//
// Return proxy to writer so we can intercept the requests and
// inform the metric manager of new metrics
//
func (self *App) GetWriter() engine.Writer {
	//Ensure 1:1 mapping from proxy to engine writer
	proxy := &writerProxy{
		self.engine.GetWriter(),
		self.manager,
	}
	return proxy
}

func (self *App) Run() error {
	log.Info("Setup signal handler")
	go self.signalHandler()

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
	supervisors := self.config.GetSupervisors(self)
	log.Debugf("Supervisors: %v", supervisors)
	self.manager = metric.NewManager(self, supervisors)

	log.Info("Setup metric schedules")
	err = self.engine.ConfigureSchedule(self.schedule)
	if err != nil {
		log.Errorf("Error configuring schedules %s", err.Error())
	}

	self.fittings = self.config.GetFittings()
	log.Infof("Starting all fittings: %v", self.fittings)
	var wg sync.WaitGroup
	for _, f := range self.fittings {
		wg.Add(1)
		go func(fitting fitting.Fitting, wg *sync.WaitGroup) {
			defer wg.Done()
			log.Infof("Starting fitting %v", fitting.Name())
			fitting.Start(self)
		}(f, &wg)

	}

	log.Info("Waiting for fittings to terminate")
	wg.Wait()

	log.Info("All fittings have finished. Exiting")

	return nil
}

func (self *App) stopFittings() {
	for _, fitting := range self.fittings {
		fitting.Stop()
	}
}

func (self *App) signalHandler() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for _ = range signals {
		log.Info("Received interrupt, shuting down...")
		self.stopFittings()
	}
}
