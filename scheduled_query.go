package morgoth

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nathanielc/morgoth/schedule"
	"time"
)

type QueryCallback func(q Query)

type ScheduledQuery struct {
	builder  QueryBuilder
	schedule *schedule.Schedule
	Callback QueryCallback
}

func NewScheduledQuery(builder QueryBuilder, delay, period time.Duration) *ScheduledQuery {

	schedule := &schedule.Schedule{
		Delay:  delay,
		Period: period,
	}
	sq := &ScheduledQuery{
		builder:  builder,
		schedule: schedule,
	}
	schedule.Callback = sq.callback

	return sq

}

func (self *ScheduledQuery) Start() {
	self.schedule.Start()
}

func (self *ScheduledQuery) callback(start, stop time.Time) {
	query := self.builder.GetForTimeRange(start, stop)
	if self.callback != nil {
		glog.V(1).Info("Scheduling query:", query)
		self.Callback(query)
	}
}
