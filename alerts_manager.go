package morgoth

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
)

type AlertsManager struct {
	scheduledAlertQueries []*scheduledAlertQuery
	engine                Engine
	queryQueue            chan alertQuery
}

func NewAlertsManager(
	engine Engine,
	scheduledAlertQueries []*scheduledAlertQuery,
) *AlertsManager {
	return &AlertsManager{
		scheduledAlertQueries: scheduledAlertQueries,
		engine:                engine,
		queryQueue:            make(chan alertQuery, queryBufferSize),
	}
}

func (self *AlertsManager) Start() {
	for _, sq := range self.scheduledAlertQueries {
		sq.Start(self.queryQueue)
	}

	go self.processQueries()
}

func (self *AlertsManager) processQueries() {
	for query := range self.queryQueue {
		glog.V(2).Info("Executing query:", query)
		windows, err := self.engine.GetWindows(query.query)
		if err != nil {
			glog.Errorf("Failed to execute query: '%s' %s", query, err)
			continue
		}
		glog.Info(windows)

		for _, w := range windows {
			sum := 0.0
			for _, point := range w.Data {
				sum += point
			}
			glog.V(1).Info("Alert Query total ", sum)
			if sum > query.threshold {
				for _, n := range query.notifiers {
					n.Notify(query.message, w)
				}
			}
		}
	}
}

type scheduledAlertQuery struct {
	sq        *ScheduledQuery
	threshold float64
	notifiers []Notifier
	message   string
	queue     chan alertQuery
}

func (self *scheduledAlertQuery) Start(queue chan alertQuery) {
	self.queue = queue
	self.sq.Callback = self.callback
	self.sq.Start()
}

func (self *scheduledAlertQuery) callback(query Query) {
	aq := alertQuery{
		query:     query,
		threshold: self.threshold,
		notifiers: self.notifiers,
		message:   self.message,
	}

	self.queue <- aq
}

type alertQuery struct {
	query     Query
	threshold float64
	notifiers []Notifier
	message   string
}
