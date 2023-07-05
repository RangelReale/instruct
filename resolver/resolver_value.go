package resolver

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/RangelReale/instruct/coerce"
	"github.com/RangelReale/instruct/types"
)

// ValueResolver resolves simple types for a Resolver.
// It should NOT handle slices, pointers, or maps.
type ValueResolver interface {
	// ResolveValue resolve the value to the proper type and return the value.
	ResolveValue(target reflect.Value, value any) error
}

// TypeValueResolver is a custom type handler for a ValueResolver.
// It should NOT process value using reflection (for performance reasons).
type TypeValueResolver interface {
	ResolveTypeValue(target reflect.Value, value any) error
}

// TypeValueResolverReflect is a custom type handler for a ValueResolver.
// It SHOULD process value using reflection.
type TypeValueResolverReflect interface {
	ResolveTypeValueReflect(target reflect.Value, sourceValue reflect.Value, value any) error
}

type DefaultValueResolver struct {
	CustomTypes        []TypeValueResolver
	CustomTypesReflect []TypeValueResolverReflect
}

func (r DefaultValueResolver) ResolveValue(target reflect.Value, value any) error {
	if !target.CanSet() {
		return fmt.Errorf("cannot set '%s' ", target.Type().Kind())
	}

	// resolve custom types without reflection, like time.Time
	if target.CanInterface() {
		for _, customType := range r.CustomTypes {
			err := customType.ResolveTypeValue(target, value)
			if err == nil {
				return nil
			}
			if errors.Is(err, types.ErrCoerceUnknown) {
				continue
			}
			return err
		}
	}

	// resolve primitive types without reflection
	switch target.Type().Kind() {
	case reflect.Bool:
		c, err := coerce.Bool(value)
		target.SetBool(c)
		return err
	case reflect.Float32:
		c, err := coerce.Float32(value)
		target.SetFloat(float64(c))
		return err
	case reflect.Float64:
		c, err := coerce.Float64(value)
		target.SetFloat(c)
		return err
	case reflect.Int:
		c, err := coerce.Int(value)
		target.SetInt(int64(c))
		return err
	case reflect.Int8:
		c, err := coerce.Int8(value)
		target.SetInt(int64(c))
		return err
	case reflect.Int16:
		c, err := coerce.Int16(value)
		target.SetInt(int64(c))
		return err
	case reflect.Int32:
		c, err := coerce.Int32(value)
		target.SetInt(int64(c))
		return err
	case reflect.Int64:
		c, err := coerce.Int64(value)
		target.SetInt(int64(c))
		return err
	case reflect.Uint:
		c, err := coerce.Uint(value)
		target.SetUint(uint64(c))
		return err
	case reflect.Uint8:
		c, err := coerce.Uint8(value)
		target.SetUint(uint64(c))
		return err
	case reflect.Uint16:
		c, err := coerce.Uint16(value)
		target.SetUint(uint64(c))
		return err
	case reflect.Uint32:
		c, err := coerce.Uint32(value)
		target.SetUint(uint64(c))
		return err
	case reflect.Uint64:
		c, err := coerce.Uint64(value)
		target.SetUint(uint64(c))
		return err
	case reflect.String:
		c, err := coerce.String(value)
		target.SetString(c)
		return err
	}

	// resolve using reflection
	sourceValue := reflect.ValueOf(value)

	// resolve custom types using reflection.
	for _, customType := range r.CustomTypesReflect {
		err := customType.ResolveTypeValueReflect(target, sourceValue, value)
		if err == nil {
			return nil
		}
		if errors.Is(err, types.ErrCoerceUnknown) {
			continue
		}
		return err
	}

	// check if types are directly assignable

	if target.Type().AssignableTo(sourceValue.Type()) {
		// the source can be directly assigned to the target
		target.Set(sourceValue)
		return nil
	}

	if target.Type().Kind() != reflect.Slice && sourceValue.Type().ConvertibleTo(target.Type()) {
		// the value is convertible to the target type
		// (slices are handled manually)
		target.Set(sourceValue.Convert(target.Type()))
		return nil
	}

	return fmt.Errorf("%w: cannot coerce source of type '%T' into target of type '%s'",
		types.ErrCoerceUnknown, value, target.Type().Kind())
}
