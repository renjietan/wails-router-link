package utility

import (
	"reflect"
	"strconv"
)

func IsEmptyValue(obj interface{}) bool {
	if obj == nil {
		return true
	}

	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0
	default:
		return reflect.DeepEqual(obj, reflect.Zero(reflect.TypeOf(obj)).Interface())
	}
}

func BoolValue(str string) bool {
	value, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}
	return value
}

func FloatValue(str string) float64 {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return value
}

func IntValue(str string, defaultValue int) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return value
}
