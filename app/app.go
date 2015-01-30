package app

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nvcook42/morgoth/config"
	"github.com/nvcook42/morgoth/detector/metadata"
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
	manager    mtypes.Manager
	engine     engine.Engine
	config     *config.Config
	fittings   []fitting.Fitting
	schedule   schedule.Schedule
	metastores map[string]metadata.MetadataStore
}

func New(config *config.Config) *App {
	app := App{
		config: config,
	}
	app.schedule = config.GetSchedule()
	app.metastores = make(map[string]metadata.MetadataStore)
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

func (self *App) GetMetadataStore(detectorID string) (metadata.MetadataStore, error) {
	ms, ok := self.metastores[detectorID]
	if !ok {

		dir := self.config.Morgoth.MetaDir
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			glog.Infof("MetaDir does not exist: creating dir '%s'", dir)
			err = os.Mkdir(dir, 0755)
			if err != nil {
				return nil, err
			}
		}
		newMS, err := metadata.New(dir, detectorID)
		if err != nil {
			return nil, err
		}
		self.metastores[detectorID] = newMS
		ms = newMS
	}
	return ms, nil
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
	for _, f := range self.fittings {
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
