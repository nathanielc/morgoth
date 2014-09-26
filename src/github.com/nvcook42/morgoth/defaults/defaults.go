package defaults

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

const tag string = "default"

type Defaulter interface {
	// Sets default values on a struct
	// NOTE: This method is not for validation
	//       hence no errors are returned
	Default()
}

// Set the tag defined default of a field specified by fieldName on
// the struct pointed to by obj, regardless of the current value of the field.
func SetDefault(obj interface{}, fieldName string) error {

	objValue := reflect.ValueOf(obj)
	if objValue.Kind() != reflect.Ptr {
		return errors.New("Must pass pointer to obj")
	}

	elem := objValue.Elem()
	elemType := elem.Type()

	if elem.Kind() != reflect.Struct {
		return errors.New("Cannot default fields of non struct")
	}

	field, exists := elemType.FieldByName(fieldName)
	if !exists {
		return errors.New(fmt.Sprintf("Not field %s exists", fieldName))
	}
	defaultStr := field.Tag.Get(tag)
	if len(defaultStr) == 0 {
		return errors.New(fmt.Sprintf("No default value specified for field '%s'", fieldName))
	}
	fieldValue := elem.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return errors.New(fmt.Sprintf("Field '%s' is not valid", fieldName))
	}
	if !fieldValue.CanSet() {
		return errors.New(fmt.Sprintf("Field '%s' is not settable", fieldName))
	}
	switch field.Type.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(defaultStr)
		if err != nil {
			return err
		}
		fieldValue.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(defaultStr, 10, 64)
		if err != nil {
			return err
		}
		fieldValue.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(defaultStr, 10, 64)
		if err != nil {
			return err
		}
		fieldValue.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(defaultStr, 64)
		if err != nil {
			return err
		}
		fieldValue.SetFloat(f)
	case reflect.String:
		fieldValue.SetString(defaultStr)
	}
	return nil
}
