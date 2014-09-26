package validator

import (
	"errors"
	"fmt"
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

	num := value.NumField()
	for i := 0; i < num; i++ {
		field := value.Field(i)
		inter := field.Interface()
		v, ok := inter.(Validator)
		if ok {
			valid := v.Validate()
			if valid != nil {
				t := reflect.TypeOf(obj)
				ft := reflect.TypeOf(inter)
				return errors.New(fmt.Sprintf("Field %v of %v is not valid: %s", ft, t, valid.Error()))
			}
		} else {
			t := reflect.TypeOf(obj)
			ft := reflect.TypeOf(inter)
			return errors.New(fmt.Sprintf("Field %v of %v is not a validator", ft, t))
		}
	}

	return nil
}
