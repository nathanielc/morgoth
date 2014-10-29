package types

import (
	"errors"
	"regexp"
	"github.com/nvcook42/morgoth/schedule"
)

type Pattern string

func (self Pattern) Validate() error {
	if len(self) == 0 {
		return errors.New("Pattern cannot be empty")
	}
	_, err := regexp.Compile(string(self))
	return err
}

func (self *Pattern) GetString(rotation *schedule.Rotation) string {
	if rotation != nil {
		return rotation.String() + "." + string(*self)
	}

	return string(*self)
}
