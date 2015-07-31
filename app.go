package morgoth

import (
	"fmt"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

type App struct {
	apiServer     *APIServer
	alertsManager *AlertsManager
	mapper        *Mapper
	manager       *Manager
	engine        Engine
	config        *Config
	shutdownChan  chan bool
	Stats         AppStats
}

type AppStats struct {
	MapperStats *MapperStats
}

func NewApp(config *Config) *App {
	app := App{
		config:       config,
		shutdownChan: make(chan bool),
	}
	return &app
}

func (self *App) Run() (err error) {
	glog.Info("Starting app")
	glog.Info("Setup signal handler")
	go self.signalHandler()

	glog.Info("Setup Engine")
	self.engine, err = self.config.EngineConf.GetEngine()
	if err != nil {
		return
	}
	err = self.engine.Initialize()
	if err != nil {
		return
	}

	glog.Info("Setup Mapper")
	detectorMatchers := make([]*DetectorMatcher, len(self.config.Mappings))
	for i, mapping := range self.config.Mappings {
		namePattern, err := regexp.Compile(mapping.Name)
		if err != nil {
			return fmt.Errorf("Error parsing regexp: %s: %s", mapping.Name, err)
		}
		tagPatterns := make(map[string]*regexp.Regexp, len(mapping.Tags))
		for tag, pattern := range mapping.Tags {
			r, err := regexp.Compile(pattern)
			if err != nil {
				return fmt.Errorf("Error parsing regexp: %s: %s", pattern, err)
			}
			tagPatterns[tag] = r
		}
		fingerprinters := make([]Fingerprinter, len(mapping.Detector.Fingerprinters))
		for i, fp := range mapping.Detector.Fingerprinters {
			f, err := fp.GetFingerprinter()
			if err != nil {
				return fmt.Errorf("Error creating Fingerprinter: %s", err)
			}

			fingerprinters[i] = f
		}
		detectorBuilder := NewDetectorBuilder(
			mapping.Detector.NormalCount,
			mapping.Detector.Consensus,
			mapping.Detector.MinSupport,
			mapping.Detector.ErrorTolerance,
			fingerprinters,
		)
		glog.V(1).Infof("Created detector builder: %v", mapping.Detector)
		matcher := NewDetectorMatcher(namePattern, tagPatterns, detectorBuilder)
		glog.V(1).Infof("Created detector matcher: %v", matcher)
		detectorMatchers[i] = matcher

	}
	detectorMappers := make([]*DetectorMapper, 0)
	self.mapper = NewMapper(detectorMappers, detectorMatchers)
	self.Stats.MapperStats = &self.mapper.Stats

	glog.Info("Setup Manager")
	scheduledDataQueries := make([]*scheduledDataQuery, len(self.config.Schedules))
	for i, sc := range self.config.Schedules {
		builder, err := self.engine.NewQueryBuilder(sc.Query, sc.GroupByInterval)
		if err != nil {
			return fmt.Errorf("Invalid query string: '%s', %s", sc.Query, err)
		}
		q := &scheduledDataQuery{
			sq: NewScheduledQuery(
				builder,
				sc.Delay,
				sc.Period,
			),
			tags: sc.Tags,
		}
		scheduledDataQueries[i] = q
	}

	self.manager = NewManager(self.mapper, self.engine, scheduledDataQueries)

	glog.Info("Starting Manager...")
	self.manager.Start()

	glog.Infof("Starting APIServer on :%d", self.config.APIPort)
	self.apiServer = NewAPIServer(self, self.config.APIPort)
	self.apiServer.Start()

	glog.Infof("Starting Alert Manger...")
	scheduledAlertQueries := make([]*scheduledAlertQuery, len(self.config.Alerts))
	for i, ac := range self.config.Alerts {
		builder, err := self.engine.NewQueryBuilder(ac.Query, ac.GroupByInterval)
		if err != nil {
			return fmt.Errorf("Invalid query string: '%s', %s", ac.Query, err)
		}

		notifiers := make([]Notifier, len(ac.Notifiers))
		for i, nc := range ac.Notifiers {
			n, err := nc.GetNotifier()
			if err != nil {
				return fmt.Errorf("Invalid Notifier for query: '%s' Err: %s", ac.Query, err)
			}
			notifiers[i] = n
		}

		q := &scheduledAlertQuery{
			sq: NewScheduledQuery(
				builder,
				ac.Delay,
				ac.Period,
			),
			threshold: ac.Threshold,
			notifiers: notifiers,
			message:   ac.Message,
		}

		scheduledAlertQueries[i] = q
	}

	self.alertsManager = NewAlertsManager(self.engine, scheduledAlertQueries)
	self.alertsManager.Start()

	// Wait for shutdown
	<-self.shutdownChan

	// Begin shutdown
	self.manager.Stop()
	self.apiServer.Stop()

	glog.Info("App shutdown")
	return
}

func (self *App) shutdownHandler() {
	self.shutdownChan <- true
	glog.Info("App shutdown handler complete")
}

func (self *App) signalHandler() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	for signal := range signals {
		glog.Infof("Received %s, shuting down...", signal)
		self.shutdownHandler()
	}
}
