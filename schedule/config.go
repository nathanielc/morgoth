package schedule

import (
	"errors"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	config "github.com/nathanielc/morgoth/config/types"
	"strconv"
	"time"
)

type RotationConf struct {
	Period string `validate:"timestr"`
}

type ScheduleConf struct {
	Rotations []RotationConf
	Delay     string `validate:"timestr" default:"1m"`
}

//Sets any invalid fields to their defualt value
func (self *ScheduleConf) Default() {
	config.PerformDefault(self)
	if self.Rotations == nil {
		self.Rotations = []RotationConf{}
		glog.Warning("No rotations were configured. Using one 5m, 10s rotation")
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
		s.Rotations[i].Period = p
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
