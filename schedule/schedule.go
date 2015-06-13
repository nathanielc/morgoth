package schedule

import (
	"errors"
	"fmt"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"time"
)

const (
	day = time.Duration(24 * time.Hour)
)

type ScheduleFunc func(rotation Rotation, start, stop time.Time)

type Rotation struct {
	Period time.Duration
}

func (self *Rotation) GetPrefix() string {
	return fmt.Sprintf("rot.%d.", self.Period/time.Second)
}

type Schedule struct {
	Callback  ScheduleFunc
	Delay     time.Duration
	Rotations []Rotation
	running   bool
}

func (self *Schedule) Start() error {
	if self.running {
		return errors.New("Schedule already started")
	}
	self.running = true

	for _, rotation := range self.Rotations {
		period := rotation.Period
		stop := time.Now()
		if period > day {
			stop = stop.Truncate(day)
		} else {
			stop = stop.Truncate(period)
		}
		stop = stop.Add(period)
		go func(rotation Rotation, stop time.Time, period time.Duration) {
			now := time.Now()
			glog.V(2).Info("Starting schedule", rotation, stop.Add(self.Delay), stop.Add(self.Delay).Sub(now))
			time.Sleep(stop.Add(self.Delay).Sub(now))
			ticker := time.NewTicker(period)
			for self.running {
				self.Callback(rotation, stop.Add(-period), stop)
				stop = stop.Add(period)
				<-ticker.C
			}
			ticker.Stop()
		}(rotation, stop, period)
	}
	return nil
}

func (self *Schedule) Stop() {
	self.running = false
}
