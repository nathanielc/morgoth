package morgoth

import (
	"github.com/nathanielc/morgoth/schedule"
	"time"
)

type Manager struct {
	scheduledQueries []*ScheduledQuery
	mapper           *Mapper
}

func (self *Manager) Start() {
	for _, sq := range self.scheduledQueries {
		sq.Start()
	}
}

func (self *Manager) ProcessWindows(windows []*Window) {

	var detector *detection.Detection
	for _, w := range windows {
		detector = mapper.Map(w)
		if detector == nil {
			continue
		}

		if detector.IsAnomalous(w) {
			self.RecordAnomalous(w)
		}
	}
}

func (self *Manager) RecordAnomalous(w *Window) {
	//TODO
}

type ScheduledQuery struct {
	schedule *schedule.Schedule
	query    Query
	engine   Engine
	manager  *Manager
}

func (self *ScheduledQuery) Start() {
	self.schedule.Start()
}

func (self *ScheduledQuery) callback(rot schedule.Rotation, start, stop time.Time) {
	q := query.QueryForTimeRange(start, stop)
	windows := ExecuteQuery(q)
	self.manager.ProcessWindows(windows)
}
