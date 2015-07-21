package morgoth

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nathanielc/morgoth/schedule"
	"time"
)

type ScheduledQueryBuilder struct {
	builder    QueryBuilder
	schedule   *schedule.Schedule
	queryQueue chan Query
	tags       map[string]string
}

func NewScheduledQueryBuilder(builder QueryBuilder, delay, period time.Duration, tags map[string]string) *ScheduledQueryBuilder {

	schedule := &schedule.Schedule{
		Delay:  delay,
		Period: period,
	}
	sq := &ScheduledQueryBuilder{
		builder:  builder,
		schedule: schedule,
		tags:     tags,
	}
	schedule.Callback = sq.callback

	return sq

}

func (self *ScheduledQueryBuilder) Start(queryQueue chan Query) {
	self.queryQueue = queryQueue
	self.schedule.Start()
}

func (self *ScheduledQueryBuilder) callback(start, stop time.Time) {
	query := self.builder.GetForTimeRange(start, stop)
	query.tags = self.tags
	glog.V(1).Info("Scheduling query:", query)
	self.queryQueue <- query
}
