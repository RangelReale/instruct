package resolver

import (
	"fmt"
	"reflect"
)

// Resolver is the default Resolver.
type Resolver struct {
	valueResolver ValueResolver
}

// NewResolver creates a new default Resolver.
// If not ValueResolver was set, a default one that suports only primitive types is used.
func NewResolver(options ...Option) *Resolver {
	ret := &Resolver{}
	for _, opt := range options {
		opt(ret)
	}
	if ret.valueResolver == nil {
		ret.valueResolver = NewDefaultValueResolver()
	}
	return ret
}

type Option func(resolver *Resolver)

// WithValueResolver sets a custom ValueResolver to be used instead of the default.
func WithValueResolver(valueResolver ValueResolver) Option {
	return func(r *Resolver) {
		r.valueResolver = valueResolver
	}
}

func (r Resolver) Resolve(target reflect.Value, value any) error {
	if target.Kind() == reflect.Slice {
		if !target.CanSet() {
			return fmt.Errorf("cannot set '%s' ", target.Type().Kind())
		}

		sourceValue := reflect.ValueOf(value)

		if sourceValue.Type().Kind() != reflect.Slice {
			return fmt.Errorf("expected an array to coerce an array into")
		}
		elemType := target.Type().Elem()
		targetSliceValue := reflect.MakeSlice(reflect.SliceOf(elemType), 0, 0)
		for i := 0; i < sourceValue.Len(); i++ {
			targetValue := reflect.New(elemType)
			if err := r.Resolve(targetValue.Elem(), sourceValue.Index(i).Interface()); err != nil {
				return err
			}
			targetSliceValue = reflect.Append(targetSliceValue, targetValue.Elem())
		}
		target.Set(targetSliceValue)
		return nil
	} else if target.Kind() == reflect.Pointer {
		ptrValue := reflect.New(target.Type().Elem())
		if err := r.Resolve(ptrValue.Elem(), value); err != nil {
			return err
		}
		target.Set(ptrValue)
		return nil
	}

	return r.valueResolver.ResolveValue(target, value)
}
