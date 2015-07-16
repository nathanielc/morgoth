// The defaults package provides a way for default values to be
// defined in the tags of struct fields.
package defaults

import (
	"errors"
	"fmt"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"reflect"
	"strconv"
)

const tag string = "default"

// Sets default values on a struct
//
// NOTE: This method is not for validation
// hence no errors are returned
type Defaulter interface {
	Default()
}

// Set the tag defined default of a field specified by fieldName on
// the struct pointed to by obj, regardless of the current value of the field.
//
// Example:
//   type MyConfigStruct struct {
//       A int    `default:"42"`
//       B string `default:"this is a string literal"`
//   }
//
//   c := MyConfigStruct{}
//   SetDefault(&c, "A")
//   SetDefault(&c, "B")
//
//   fmt.Printf("%v", c) // {42 this is a string literal}
//
// See the test cases for more examples.
//
// NOTE: It is considered an error to call SetDefault if no default has been
// defined. Use HasDefault if neccessary
func SetDefault(obj interface{}, fieldName string) (interface{}, error) {

	objValue := reflect.ValueOf(obj)
	if objValue.Kind() != reflect.Ptr {
		return nil, errors.New("Must pass pointer to obj")
	}

	elem := objValue.Elem()
	elemType := elem.Type()

	if elem.Kind() != reflect.Struct {
		return nil, errors.New("Cannot default fields of non struct")
	}

	field, exists := elemType.FieldByName(fieldName)
	if !exists {
		return nil, errors.New(fmt.Sprintf("Not field %s exists", fieldName))
	}
	defaultStr := field.Tag.Get(tag)
	if len(defaultStr) == 0 {
		return nil, errors.New(fmt.Sprintf("No default value specified for field '%s'", fieldName))
	}
	fieldValue := elem.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return nil, errors.New(fmt.Sprintf("Field '%s' is not valid", fieldName))
	}
	if !fieldValue.CanSet() {
		return nil, errors.New(fmt.Sprintf("Field '%s' is not settable", fieldName))
	}
	switch field.Type.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(defaultStr)
		if err != nil {
			return nil, err
		}
		fieldValue.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(defaultStr, 10, 64)
		if err != nil {
			return nil, err
		}
		fieldValue.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(defaultStr, 10, 64)
		if err != nil {
			return nil, err
		}
		fieldValue.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(defaultStr, 64)
		if err != nil {
			return nil, err
		}
		fieldValue.SetFloat(f)
	case reflect.String:
		fieldValue.SetString(defaultStr)
	}
	return fieldValue.Interface(), nil
}

// Check whether a specified field has a default defined
func HasDefault(obj interface{}, fieldName string) (bool, error) {

	objValue := reflect.Indirect(reflect.ValueOf(obj))

	objType := objValue.Type()

	if objType.Kind() != reflect.Struct {
		return false, errors.New("Cannot default fields of non struct")
	}

	field, exists := objType.FieldByName(fieldName)
	if !exists {
		return false, errors.New(fmt.Sprintf("Not field %s exists", fieldName))
	}
	defaultStr := field.Tag.Get(tag)
	return len(defaultStr) > 0, nil
}

// Set the defaults of all fields of obj that have been defined
func SetAllDefaults(obj interface{}) error {

	objValue := reflect.ValueOf(obj)
	if objValue.Kind() != reflect.Ptr {
		return errors.New("Must pass pointer to obj")
	}

	elem := objValue.Elem()
	elemType := elem.Type()

	num := elemType.NumField()
	for i := 0; i < num; i++ {
		field := elemType.Field(i)
		glog.Infof("Defaulting %s", field.Name)
		if def, err := HasDefault(obj, field.Name); def && err == nil {
			_, err := SetDefault(obj, field.Name)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}
