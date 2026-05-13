package utility

import (
	"reflect"
)

// Interface2Interface Mortal 2026-04-25 15:04:22 CST
// Example:
//
//	type test struct {
//		 Name string `json:"name"`
//	}
//
// c := test{}
//
//	a := Interface2String(map[string]interface{}{
//		 name: "tttttttttttt",
//	}, &c)
func Interface2Interface(value interface{}, res interface{}) error {
	v := Interface2String(value)
	err := JsonDecode(v, &res)
	if err != nil {
		return err
	}
	return nil
}

// Tern Mortal 2026-04-26 13:49:33 CST 三元表达式
func Tern[T any, K any](boolVal bool, a T, b K) any {
	var result any
	if boolVal {
		result = a
	} else {
		result = b
	}
	val := reflect.ValueOf(result)
	if val.Kind() == reflect.Func {
		funcType := val.Type()
		numIn := funcType.NumIn()
		if numIn == 0 {
			results := val.Call(nil)
			if len(results) > 0 {
				return results[0].Interface()
			}
			return nil
		}
		return result
	}

	return result
}
