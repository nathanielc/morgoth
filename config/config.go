package config

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/gopkg.in/validator.v2"
	"github.com/nathanielc/morgoth/defaults"
	"reflect"
	"strings"
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
				if strings.Contains(fieldName, ".") {
					//If field Name contains a '.' than it is a subfield of a current field
					continue
				}
				glog.V(4).Info("Invalid field searching for default ", fieldName)
				if ok, err := defaults.HasDefault(conf, fieldName); ok {
					value, err := defaults.SetDefault(conf, fieldName)
					if err != nil {
						glog.Errorf("Failed to set default of %s on %s", fieldName, name)
					} else {
						glog.Infof("Defaulted %v.%s to '%v'", name, fieldName, value)
					}
				} else {
					glog.V(4).Info("No default found: ", err)
				}
			}
		} else {
			glog.Errorf("Failed to validate %s: %s", name, err)
		}
	}
}
