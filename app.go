package morgoth

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nathanielc/morgoth/schedule"
	"os"
	"os/signal"
	"sync"
)

type App struct {
	mapper   *Mapper
	engine   Engine
	config   *config.Config
	schedule *schedule.Schedule
}

func New(config *config.Config) *App {
	app := App{
		config: config,
	}
	app.schedule = config.GetSchedule()
	return &app
}

func (self *App) Run() error {
	glog.Info("Setup signal handler")
	go self.signalHandler()

	glog.Info("Setup engine")
	eng, err := self.config.EngineConf.GetEngine()
	if err != nil {
		return err
	}
	self.engine = eng
	err = self.engine.Initialize()
	if err != nil {
		return err
	}

	glog.Info("Setup metrics manager")
	supervisors := self.config.GetSupervisors(self)
	glog.V(2).Infof("Supervisors: %v", supervisors)
	self.manager = metric.NewManager(supervisors)

	glog.Info("Setup metric schedules")
	err = self.engine.ConfigureSchedule(&self.schedule)
	if err != nil {
		glog.Errorf("Error configuring schedules %s", err.Error())
	}

	self.fittings = self.config.GetFittings()
	glog.Infof("Starting all fittings: %v", self.fittings)
	var wg sync.WaitGroup
	for i, f := range self.fittings {
		if f == nil {
			glog.Errorf("Fitting #%d is nil", i+1)
			continue
		}
		wg.Add(1)
		go func(fitting fitting.Fitting, wg *sync.WaitGroup) {
			defer wg.Done()
			glog.Infof("Starting fitting %v", fitting.Name())
			fitting.Start(self)
		}(f, &wg)

	}

	glog.Info("Starting metric manager")
	self.schedule.Callback = self.manager.Detect
	self.schedule.Start()

	glog.Info("Waiting for fittings to terminate")
	wg.Wait()

	glog.Info("All fittings have finished. Exiting")

	return nil
}

func (self *App) shutdown() {
	glog.V(2).Info("Closing all metastores...")
	for _, db := range self.metastores {
		db.Close()
	}
	glog.V(2).Info("Stopping all fittings...")
	for _, fitting := range self.fittings {
		fitting.Stop()
	}
	glog.Info("App shutdown complete")
}

func (self *App) signalHandler() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for _ = range signals {
		glog.Info("Received interrupt, shuting down...")
		self.shutdown()
	}
}
