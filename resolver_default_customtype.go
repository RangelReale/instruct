package instruct

import (
	"reflect"
	"time"

	"github.com/RangelReale/instruct/coerce"
)

// DefaultResolverValueResolverTime resolves time.Time values.
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
		}
	}
	return ErrCoerceUnknown
}

func (d *DefaultResolverValueResolverTime) ResolveCustomTypeValueReflect(target reflect.Value,
	sourceValue reflect.Value, value any) error {
	return ErrCoerceUnknown
}

// DefaultResolverValueResolverTimeDuration resolves time.Duration values.
type DefaultResolverValueResolverTimeDuration struct {
}

func (d *DefaultResolverValueResolverTimeDuration) ResolveCustomTypeValue(target reflect.Value, value any) error {
	if target.CanInterface() {
		switch target.Interface().(type) {
		case time.Duration:
			c, err := coerce.TimeDuration(value)
			target.Set(reflect.ValueOf(c))
			return err
		}
	}
	return ErrCoerceUnknown
}

func (d *DefaultResolverValueResolverTimeDuration) ResolveCustomTypeValueReflect(target reflect.Value,
	sourceValue reflect.Value, value any) error {
	return ErrCoerceUnknown
}
