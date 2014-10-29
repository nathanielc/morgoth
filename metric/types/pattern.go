package types

import (
	"errors"
	"regexp"
)

type Pattern string

func (self Pattern) Validate() error {
	if len(self) == 0 {
		return errors.New("Pattern cannot be empty")
	}
	_, err := regexp.Compile(string(self))
	return err
}
