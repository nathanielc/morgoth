package schedule

import (
	"errors"
	"time"
)

type ScheduleFunc func()

type Schedule struct {
	f       ScheduleFunc
	period  time.Duration
	running bool
}

func NewSchedule(f ScheduleFunc, period time.Duration) *Schedule {
	s := new(Schedule)
	s.f = f
	s.period = period
	s.running = false
	return s
}

func (self *Schedule) Start() error {
	if self.running {
		return errors.New("Schedule already started")
	}
	self.running = true

	go func() {
		ticker := time.NewTicker(self.period)
		for self.running {
			self.f()
			<-ticker.C
		}
		ticker.Stop()
	}()

	return nil
}

func (self *Schedule) Stop() {
	self.running = false
}
