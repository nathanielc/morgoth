package schedule_test

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/schedule"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestScheduleShouldStartAndStop(t *testing.T) {
	defer log.Flush()
	assert := assert.New(t)

	unit := 10 * time.Millisecond
	calledCount := 0
	testF := func(start time.Time, dur time.Time) {
		remainder := start.Nanosecond() % int(unit)
		assert.Equal(0, remainder, "Start not truncated to 10ms")
		calledCount++
	}
	s := schedule.Schedule{
		Callback:  testF,
		Rotations: []schedule.Rotation{schedule.Rotation{Period: unit}},
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
	defer log.Flush()
	assert := assert.New(t)

	s := schedule.Schedule{
		Callback:  func(start time.Time, dur time.Time) {},
		Rotations: []schedule.Rotation{schedule.Rotation{Period: time.Millisecond}},
	}

	assert.NotNil(s)

	err := s.Start()
	assert.Nil(err)

	err = s.Start()
	assert.NotNil(err)

	s.Stop()
}

func TestRotationShouldConvertToString(t *testing.T) {
	defer log.Flush()
	assert := assert.New(t)

	r := schedule.Rotation{
		Period:     time.Minute * 6,
		Resolution: time.Second * 7,
	}

	str := r.String()
	assert.Equal("rot.7.360", str)
}
