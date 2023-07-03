package instruct

import (
	"reflect"
	"time"

	"github.com/RangelReale/instruct/coerce"
)

type DefaultResolverValueResolverTime struct {
	layout string
}

func NewDefaultResolverValueResolverTypeTime(layout string) *DefaultResolverValueResolverTime {
	return &DefaultResolverValueResolverTime{
		layout: layout,
	}
}

func (d *DefaultResolverValueResolverTime) ResolveCustomTypeValue(target reflect.Value, value any) error {
	if target.CanInterface() {
		switch target.Interface().(type) {
		case time.Time:
			c, err := coerce.Time(value, d.layout)
			target.Set(reflect.ValueOf(c))
			return err

			// switch v := value.(type) {
			// case time.Time:
			// 	target.Set(reflect.ValueOf(v))
			// 	return nil
			// case string:
			// 	t, err := time.Parse(time.RFC3339, v)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	target.Set(reflect.ValueOf(t))
			// 	return nil
			// }
		}
	}
	return ErrCoerceUnknown
}

func (d *DefaultResolverValueResolverTime) ResolveCustomTypeValueReflect(target reflect.Value,
	sourceValue reflect.Value, value any) error {
	return ErrCoerceUnknown
}
