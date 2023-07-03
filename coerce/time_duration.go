package coerce

import (
	"fmt"
	"reflect"
	"time"
)

// TimeDuration coerces v to time.Duration.
func TimeDuration(v interface{}) (time.Duration, error) {
	for {
		switch sw := v.(type) {
		case time.Duration:
			return sw, nil
		case nil:
			return time.Duration(0), nil
		case string:
			t, err := time.ParseDuration(sw)
			if err != nil {
				return time.Duration(0), err
			}
			return t, nil
		}
		//
		// Beyond this point we need reflection.
		T := reflect.TypeOf(v)
		//
		// - T.Kind() is a primitive
		//		convert to actual primitive and try again
		// - T.Kind() is a pointer
		//		dereference pointer and try again
		switch T.Kind() {
		case reflect.String:
			v = reflect.ValueOf(v).Convert(TypeString).Interface().(string)
			continue
		case reflect.Ptr:
			rv := reflect.ValueOf(v)
			for ; rv.Kind() == reflect.Ptr; rv = rv.Elem() {
				if rv.IsNil() {
					return time.Duration(0), nil
				}
			}
			v = rv.Interface()
			continue

		}
		//
		return time.Duration(0), fmt.Errorf("%w; coerce %v to time.Duration", ErrUnsupported, v)
	}
}
