package schedule

import (
	"errors"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"time"
)

const (
	day = time.Duration(24 * time.Hour)
)

type ScheduleFunc func(start, stop time.Time)

type Schedule struct {
	Callback ScheduleFunc
	Delay    time.Duration
	Period   time.Duration
	running  bool
}

func (self *Schedule) Start() error {
	if self.running {
		return errors.New("Schedule already started")
	}
	self.running = true

	stop := time.Now()
	if self.Period > day {
		stop = stop.Truncate(day)
	} else {
		stop = stop.Truncate(self.Period)
	}
	stop = stop.Add(self.Period)
	go func(stop time.Time, period time.Duration) {
		now := time.Now()
		glog.V(2).Info("Starting schedule", stop.Add(self.Delay), stop.Add(self.Delay).Sub(now))
		time.Sleep(stop.Add(self.Delay).Sub(now))
		ticker := time.NewTicker(period)
		for self.running {
			self.Callback(stop.Add(-period), stop)
			stop = stop.Add(period)
			<-ticker.C
		}
		ticker.Stop()
	}(stop, self.Period)
	return nil
}

func (self *Schedule) Stop() {
	self.running = false
}
