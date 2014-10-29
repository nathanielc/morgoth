package schedule_test

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/schedule"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
	"time"
)

func TestScheduleConfShouldDefault(t *testing.T) {
	defer log.Flush()
	assert := assert.New(t)

	sc := schedule.ScheduleConf{}

	sc.Default()

	assert.Equal([]schedule.RotationConf{
		schedule.RotationConf{"5m", "10s"},
		schedule.RotationConf{"15m", "30s"},
		schedule.RotationConf{"1h", "1m"},
		schedule.RotationConf{"6h", "6m"},
		schedule.RotationConf{"1d", "24m"},
		schedule.RotationConf{"10d", "4h"},
	}, sc.Rotations)
	assert.Equal("1m", sc.Delay)
}

func TestScheduleConfShouldValidate(t *testing.T) {
	defer log.Flush()
	assert := assert.New(t)

	sc := schedule.ScheduleConf{
		Rotations: []schedule.RotationConf{schedule.RotationConf{"1m", "1s"}},
		Delay:     "4m",
	}

	err := sc.Validate()
	assert.Nil(err)

}

func TestScheduleConfShouldFailValidateRotations(t *testing.T) {
	defer log.Flush()
	assert := assert.New(t)

	sc := schedule.ScheduleConf{
		Rotations: []schedule.RotationConf{schedule.RotationConf{"1", "1s"}},
		Delay:     "1m",
	}

	err := sc.Validate()
	assert.NotNil(err)

	sc = schedule.ScheduleConf{
		Rotations: []schedule.RotationConf{schedule.RotationConf{"1m", "1"}},
		Delay:     "1m",
	}

	err = sc.Validate()
	assert.NotNil(err)

}

func TestScheduleConfShouldFailValidateDelay(t *testing.T) {
	defer log.Flush()
	assert := assert.New(t)

	sc := schedule.ScheduleConf{
		Rotations: []schedule.RotationConf{schedule.RotationConf{"1m", "1s"}},
		Delay:     "",
	}

	err := sc.Validate()
	assert.NotNil(err)

}

func TestScheduleConfShouldParse(t *testing.T) {
	defer log.Flush()
	assert := assert.New(t)

	var data string = `---
rotations:
  - {period: 1m, resolution: 1s}
  - {period: 1h, resolution: 1m}
  - {period: 1d, resolution: 1h}
delay: 60s
`

	sc := schedule.ScheduleConf{}

	err := yaml.Unmarshal([]byte(data), &sc)
	assert.Nil(err)

	assert.Equal([]schedule.RotationConf{
		schedule.RotationConf{"1m", "1s"},
		schedule.RotationConf{"1h", "1m"},
		schedule.RotationConf{"1d", "1h"},
	}, sc.Rotations)
	assert.Equal("60s", sc.Delay)
}

func TestScheduleConfShouldGetSchedule(t *testing.T) {
	defer log.Flush()
	assert := assert.New(t)

	sc := schedule.ScheduleConf{
		Rotations: []schedule.RotationConf{schedule.RotationConf{"1m", "1s"}},
		Delay:     "60m",
	}

	err := sc.Validate()
	assert.Nil(err)

	s := sc.GetSchedule()
	if !assert.Equal(1, len(s.Rotations)) {
		assert.Fail("Schedule rotations should have exactly one element")
	}
	assert.Equal(time.Minute, s.Rotations[0].Period)
	assert.Equal(time.Second, s.Rotations[0].Resolution)

	assert.Equal(time.Hour, s.Delay)
}

func TestStrToDuration(t *testing.T) {
	defer log.Flush()
	assert := assert.New(t)

	tm, err := schedule.StrToDuration("1s")
	assert.Nil(err)
	assert.Equal(time.Second, tm)

	tm, err = schedule.StrToDuration("2s")
	assert.Nil(err)
	assert.Equal(time.Second*2, tm)

	tm, err = schedule.StrToDuration("1m")
	assert.Nil(err)
	assert.Equal(time.Minute, tm)

	tm, err = schedule.StrToDuration("5m")
	assert.Nil(err)
	assert.Equal(time.Minute*5, tm)

	tm, err = schedule.StrToDuration("1h")
	assert.Nil(err)
	assert.Equal(time.Hour, tm)

	tm, err = schedule.StrToDuration("10h")
	assert.Nil(err)
	assert.Equal(time.Hour*10, tm)

	tm, err = schedule.StrToDuration("1d")
	assert.Nil(err)
	assert.Equal(time.Hour*24, tm)

	tm, err = schedule.StrToDuration("3d")
	assert.Nil(err)
	assert.Equal(time.Hour*24*3, tm)
}
