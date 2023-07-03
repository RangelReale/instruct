package coerce

import (
	"fmt"
	"reflect"
	"time"
)

// Time coerces v to time.Time.
func Time(v interface{}, layout string) (time.Time, error) {
	for {
		switch sw := v.(type) {
		case time.Time:
			return sw, nil
		case nil:
			return time.Time{}, nil
		case string:
			t, err := time.Parse(layout, sw)
			if err != nil {
				return time.Time{}, err
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
					return time.Time{}, nil
				}
			}
			v = rv.Interface()
			continue

		}
		//
		return time.Time{}, fmt.Errorf("%w; coerce %v to time.Time", ErrUnsupported, v)
	}
}
