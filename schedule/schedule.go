package schedule

import (
	"errors"
	"fmt"
	"time"
)

const (
	day = time.Duration(24 * time.Hour)
)

type ScheduleFunc func(time.Time, time.Time)

type Rotation struct {
	Period     time.Duration
	Resolution time.Duration
}

func (self *Rotation) String() string {
	return fmt.Sprintf("rot.%d.%d", self.Resolution/time.Second, self.Period/time.Second)
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
		start := time.Now()
		if period > day {
			start = start.Truncate(day)
		} else {
			start = start.Truncate(period)
		}
		start = start.Add(period)
		go func(start time.Time, period time.Duration) {
			now := time.Now()
			time.Sleep(start.Add(self.Delay).Sub(now))
			ticker := time.NewTicker(period)
			for self.running {
				self.Callback(start, start.Add(period))
				start = start.Add(period)
				<-ticker.C
			}
			ticker.Stop()
		}(start, period)
	}
	return nil
}

func (self *Schedule) Stop() {
	self.running = false
}
