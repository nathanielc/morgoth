package schedule

import (
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/defaults"
	"gopkg.in/validator.v2"
	"time"
)

type ScheduleConf struct {
	Duration uint `yaml:"duration" validate:"min=1" default:"60"`
	Period   uint `yaml:"period"   validate:"min=1" default:"60"`
	Delay    uint `yaml:"delay"    validate:"min=0" default:"60"`
}

//Sets any invalid fields to their defualt value
func (self *ScheduleConf) Default() {
	err := self.Validate()
	if err != nil {
		errs := err.(validator.ErrorMap)
		for fieldName := range errs {
			log.Infof("Using default for Schedule.%s", fieldName)
			defaults.SetDefault(self, fieldName)
		}
	}
}

func (self ScheduleConf) Validate() error {
	return validator.Validate(self)
}

func (self *ScheduleConf) GetSchedule() Schedule {
	s := Schedule{
		Period:   time.Duration(self.Period),
		Duration: time.Duration(self.Duration),
		Delay:    time.Duration(self.Delay),
	}

	return s
}
