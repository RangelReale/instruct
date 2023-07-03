package instruct

import (
	"encoding"
	"reflect"
	"time"

	"github.com/RangelReale/instruct/coerce"
)

var (
	textUnmarshalerType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()
)

// DefaultResolverValueResolverTime resolves time.Time values.
type DefaultResolverValueResolverTime struct {
	layout string
}

func NewDefaultResolverValueResolverTime(layout string) *DefaultResolverValueResolverTime {
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

// DefaultResolverValueResolverReflectTextUnmarshaler checks if target implements [encoding.TextUnmarshaler]
// and use it to resolve from string.
type DefaultResolverValueResolverReflectTextUnmarshaler struct {
}

func (d *DefaultResolverValueResolverReflectTextUnmarshaler) ResolveCustomTypeValueReflect(target reflect.Value,
	sourceValue reflect.Value, value any) error {
	switch sourceValue.Type().Kind() {
	case reflect.String:
		if reflect.PointerTo(target.Type()).Implements(textUnmarshalerType) {
			xtarget := reflect.New(target.Type())
			um := xtarget.Interface().(encoding.TextUnmarshaler)
			if err := um.UnmarshalText([]byte(value.(string))); err != nil {
				return err
			}
			target.Set(xtarget.Elem())
			return nil
		}
	}

	return ErrCoerceUnknown
}
