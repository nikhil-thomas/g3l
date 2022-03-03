package lgwt_reflection

import (
	"reflect"
)

func walk(x interface{}, f func(input string)) {
	val := reflect.ValueOf(x)
	//fmt.Printf("%#v", val)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String {
			f(field.String())
		}
		if field.Kind() == reflect.Struct {
			walk(field.Interface(), f)
		}
	}
}

//What if the value of the struct passed in is a pointer?
