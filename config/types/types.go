// Common interfaces needed for configuration
package types

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	"github.com/nvcook42/morgoth/defaults"
	"reflect"
)

type Validator interface {
	Validate() error
}

type Configuration interface {
	defaults.Defaulter
	Validator
}

//Sets any invalid fields to their default value
func PerformDefault(conf Configuration) {
	err := conf.Validate()
	name := reflect.ValueOf(conf).Type()
	if err != nil {
		if errs, ok := err.(validator.ErrorMap); ok {
			for fieldName := range errs {
				if ok, _ := defaults.HasDefault(conf, fieldName); ok {
					value, err := defaults.SetDefault(conf, fieldName)
					if err != nil {
						glog.Errorf("Failed to set default of %s on %s", fieldName, name)
					} else {
						glog.Infof("Defaulted %v.%s to '%v'", name, fieldName, value)
					}
				}
			}
		} else {
			glog.Errorf("Failed to validate %s: %s", name, err)
		}
	}
}
