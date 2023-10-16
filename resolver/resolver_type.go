package resolver

import (
	"reflect"
	"time"

	"github.com/rrgmc/instruct/coerce"
	"github.com/rrgmc/instruct/types"
)

// ValueResolverTime resolves time.Time values.
type ValueResolverTime struct {
	layout string
}

func NewValueResolverTime(layout string) *ValueResolverTime {
	return &ValueResolverTime{
		layout: layout,
	}
}

func (d *ValueResolverTime) ResolveTypeValue(target reflect.Value, value any) error {
	if target.CanInterface() {
		switch target.Interface().(type) {
		case time.Time:
			c, err := coerce.Time(value, d.layout)
			target.Set(reflect.ValueOf(c))
			return err
		}
	}
	return types.ErrCoerceUnknown
}

// ValueResolverTimeDuration resolves time.Duration values.
type ValueResolverTimeDuration struct {
}

func NewValueResolverTimeDuration() *ValueResolverTimeDuration {
	return &ValueResolverTimeDuration{}
}

func (d *ValueResolverTimeDuration) ResolveTypeValue(target reflect.Value, value any) error {
	if target.CanInterface() {
		switch target.Interface().(type) {
		case time.Duration:
			c, err := coerce.TimeDuration(value)
			target.Set(reflect.ValueOf(c))
			return err
		}
	}
	return types.ErrCoerceUnknown
}
