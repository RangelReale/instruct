package resolver

import (
	"encoding"
	"reflect"

	"github.com/RangelReale/instruct/types"
)

var (
	textUnmarshalerType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()
)

// ValueResolverReflectTextUnmarshaler checks if target implements [encoding.TextUnmarshaler]
// and use it to resolve from string.
type ValueResolverReflectTextUnmarshaler struct {
}

func (d *ValueResolverReflectTextUnmarshaler) ResolveTypeValueReflect(target reflect.Value,
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

	return types.ErrCoerceUnknown
}
