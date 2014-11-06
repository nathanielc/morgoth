package schedule

import (
	"errors"
	"github.com/golang/glog"
	"github.com/nvcook42/morgoth/defaults"
	"gopkg.in/validator.v2"
	"strconv"
	"time"
)

type RotationConf struct {
	Period     string `validate:"timestr"`
	Resolution string `validate:"timestr"`
}

type ScheduleConf struct {
	Rotations []RotationConf
	Delay     string `validate:"timestr" default:"1m"`
}

//Sets any invalid fields to their defualt value
func (self *ScheduleConf) Default() {
	err := self.Validate()
	if err != nil {
		errs := err.(validator.ErrorMap)
		for fieldName := range errs {
			if ok, _ := defaults.HasDefault(self, fieldName); ok {
				glog.Infof("Using default for Schedule.%s", fieldName)
				defaults.SetDefault(self, fieldName)
			}
		}
	}
	if self.Rotations == nil {
		self.Rotations = []RotationConf{
			RotationConf{"5m", "10s"},
			RotationConf{"15m", "30s"},
			RotationConf{"1h", "1m"},
			RotationConf{"6h", "6m"},
			RotationConf{"1d", "24m"},
			RotationConf{"10d", "4h"},
		}
	}

}

func (self ScheduleConf) Validate() error {
	validator.SetValidationFunc("timestr", validateTimeStr)
	for _, rotation := range self.Rotations {
		err := validator.Validate(rotation)
		if err != nil {
			return err
		}
	}
	return validator.Validate(self)
}

func validateTimeStr(v interface{}, param string) error {
	timestr, ok := v.(string)
	if !ok {
		return validator.ErrUnsupported
	}
	err := validator.Valid(timestr, "regexp=^[-+]?\\d+[smhd]$")
	if err != nil {
		return err
	}
	return nil
}

func (self *ScheduleConf) GetSchedule() Schedule {
	s := Schedule{}

	delay, err := StrToDuration(self.Delay)
	if err == nil {
		s.Delay = delay
	}

	s.Rotations = make([]Rotation, len(self.Rotations))
	for i, rotation := range self.Rotations {
		p, err := StrToDuration(rotation.Period)
		if err != nil {
			continue
		}
		r, err := StrToDuration(rotation.Resolution)
		if err != nil {
			continue
		}
		s.Rotations[i].Period = p
		s.Rotations[i].Resolution = r
	}

	return s
}

func StrToDuration(str string) (time.Duration, error) {
	if len(str) < 2 {
		return 0, errors.New("Invalid time string " + str)
	}
	scaleStr := str[len(str)-1]
	scale := int64(0)
	switch scaleStr {
	case 's':
		scale = 1
	case 'm':
		scale = 60
	case 'h':
		scale = 3600
	case 'd':
		scale = 86400
	default:
		return 0, errors.New("Invalid time unit " + string(scaleStr))
	}

	valueStr := str[0 : len(str)-1]
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, errors.New("Invalid time string " + err.Error())
	}

	return time.Duration(time.Duration(value*scale) * time.Second), nil
}
