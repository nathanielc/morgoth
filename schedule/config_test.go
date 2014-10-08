package schedule_test

import (
	"github.com/nvcook42/morgoth/schedule"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestScheduleConfShouldDefault(t *testing.T) {
	assert := assert.New(t)

	sc := schedule.ScheduleConf{
		Duration: 0,
		Period:   0,
		Delay:    1,
	}

	sc.Default()

	assert.Equal(60, sc.Duration)
	assert.Equal(60, sc.Period)
	assert.Equal(1, sc.Delay)
}

func TestScheduleConfShouldValidate(t *testing.T) {
	assert := assert.New(t)

	sc := schedule.ScheduleConf{
		Duration: 3,
		Period:   2,
		Delay:    1,
	}

	err := sc.Validate()
	assert.Nil(err)

}

func TestScheduleConfShouldFailValidate(t *testing.T) {
	assert := assert.New(t)

	sc := schedule.ScheduleConf{
		Duration: 0,
		Period:   0,
		Delay:    1,
	}

	err := sc.Validate()
	assert.NotNil(err)

}

func TestScheduleConfShouldParse(t *testing.T) {
	assert := assert.New(t)

	var data string = `---
duration: 45
period: 12
delay: 60
`

	sc := schedule.ScheduleConf{}

	err := yaml.Unmarshal([]byte(data), &sc)
	assert.Nil(err)

	assert.Equal(45, sc.Duration)
	assert.Equal(12, sc.Period)
	assert.Equal(60, sc.Delay)

}
