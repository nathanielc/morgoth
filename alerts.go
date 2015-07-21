package morgoth

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
)

type AlertsManager struct {
	scheduledQueries []*ScheduledQuery
	engine           Engine
	queryQueue       chan AlertQuery
}

func NewAlertsManager(
	engine Engine,
	scheduledQueries []*ScheduledQuery,
) *AlertsManager {
	return &AlertsManager{
		scheduledQueries: scheduledQueries,
		engine:           engine,
		queryQueue:       make(chan AlertQuery, 100),
	}
}

func (self *AlertsManager) Start() {
	for _, sq := range self.scheduledQueries {
		sq.Start(self.queryQueue)
	}

	go self.processQueries()
}

func (self *AlertsManager) processQueries() {
	for query := range self.queryQueue {
		glog.V(2).Info("Executing query:", query)
		windows, err := self.engine.GetWindows(query.Query)
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
			if sum > query.Threshold {
				for _, n := range query.Notifiers {
					n.Notify(query.Message, w)
				}
			}
		}
	}
}
