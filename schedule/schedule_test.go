package schedule_test

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nathanielc/morgoth/schedule"
	"testing"
	"time"
)

func TestScheduleShouldStartAndStop(t *testing.T) {
	defer glog.Flush()
	assert := assert.New(t)

	unit := 10 * time.Millisecond
	calledCount := 0
	testF := func(start time.Time, dur time.Time) {
		remainder := start.Nanosecond() % int(unit)
		assert.Equal(0, remainder, "Start not truncated to 10ms")
		calledCount++
	}
	s := schedule.Schedule{
		Callback: testF,
		Period:   unit,
	}

	assert.NotNil(s)

	sleepCount := 3
	start := time.Now()
	err := s.Start()
	assert.Nil(err)
	time.Sleep(unit * time.Duration(sleepCount))
	s.Stop()
	elapsed := time.Since(start)
	assert.True(calledCount >= sleepCount && calledCount <= int(elapsed/unit))
}

func TestScheduleShouldNotDoubleStart(t *testing.T) {
	defer glog.Flush()
	assert := assert.New(t)

	s := schedule.Schedule{
		Callback: func(start time.Time, dur time.Time) {},
		Period:   time.Millisecond,
	}

	assert.NotNil(s)

	err := s.Start()
	assert.Nil(err)

	err = s.Start()
	assert.NotNil(err)

	s.Stop()
}
