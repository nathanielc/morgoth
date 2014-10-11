package schedule

import (
	"errors"
	log "github.com/cihub/seelog"
	"time"
)

type ScheduleFunc func(time.Duration)

type Schedule struct {
	Callback ScheduleFunc
	Duration time.Duration
	Delay    time.Duration
	Period   time.Duration
	running  bool
}

func (self *Schedule) Start() error {
	if self.running {
		return errors.New("Schedule already started")
	}
	self.running = true

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Error during schedule Callback %v", r)
			}
		}()
		ticker := time.NewTicker(self.Period)
		for self.running {
			self.Callback(self.Duration)
			<-ticker.C
		}
		ticker.Stop()
	}()

	return nil
}

func (self *Schedule) Stop() {
	self.running = false
}
