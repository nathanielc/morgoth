package morgoth

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nathanielc/morgoth/schedule"
	"time"
)

type ScheduledQuery struct {
	query      AlertQuery
	schedule   *schedule.Schedule
	queryQueue chan AlertQuery
	tags       map[string]string
}

func NewScheduledQuery(query AlertQuery, delay, period time.Duration) *ScheduledQuery {

	schedule := &schedule.Schedule{
		Delay:  delay,
		Period: period,
	}
	sq := &ScheduledQuery{
		query:    query,
		schedule: schedule,
	}
	schedule.Callback = sq.callback

	return sq

}

func (self *ScheduledQuery) Start(queryQueue chan AlertQuery) {
	self.queryQueue = queryQueue
	self.schedule.Start()
}

func (self *ScheduledQuery) callback(start, stop time.Time) {
	glog.V(1).Info("Scheduling query:", self.query)
	self.queryQueue <- self.query
}
