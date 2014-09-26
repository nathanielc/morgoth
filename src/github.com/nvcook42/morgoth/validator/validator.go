package validator

import (
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"reflect"
)

type Validator interface {
	Validate() error
}

//Validate that at least one field of obj is valid
func ValidateOne(obj interface{}) error {
	value := reflect.Indirect(reflect.ValueOf(obj))
	if value.Kind() != reflect.Struct {
		return errors.New("Cannot validate non struct")
	}

	num := value.NumField()
	for i := 0; i < num; i++ {
		field := value.Field(i)
		v, ok := field.Interface().(Validator)
		if ok {
			valid := v.Validate()
			if valid == nil {
				return nil
			}
		}
	}

	return errors.New("All fields are invalid")

}

//Validate that all fields of obj are valid
func ValidateAll(obj interface{}) error {
	value := reflect.Indirect(reflect.ValueOf(obj))
	if value.Kind() != reflect.Struct {
		return errors.New("Cannot validate non struct")
	}

	allValid := true
	num := value.NumField()
	for i := 0; i < num && allValid; i++ {
		field := value.Field(i)
		inter := field.Interface()
		v, ok := inter.(Validator)
		if ok {
			valid := v.Validate()
			allValid = valid == nil
		} else {
			t := reflect.TypeOf(obj)
			ft := reflect.TypeOf(inter)
			log.Warnf("Field %v of %v is not a validator", ft, t)
			allValid = false
		}
	}

	if allValid {
		return nil
	}

	t := reflect.TypeOf(obj)
	return errors.New(fmt.Sprintf("Not all fields of %v are valid", t))
}
