package schedule_test

import (
	"github.com/nvcook42/morgoth/schedule"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestScheduleShouldStartAndStop(t *testing.T) {
	assert := assert.New(t)

	calledCount := 0
	testF := func() {
		calledCount++
	}
	unit := time.Millisecond
	s := schedule.Schedule{
		Callback: testF,
		Period: unit,
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
	assert := assert.New(t)


	s := schedule.Schedule{
		Callback: func() {},
		Period: time.Millisecond,
	}

	assert.NotNil(s)

	err := s.Start()
	assert.Nil(err)

	err = s.Start()
	assert.NotNil(err)

	s.Stop()
}
