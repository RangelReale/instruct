package instruct

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"

	"github.com/RangelReale/instruct/coerce"
)

var (
	textUnmarshalerType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()
)

type DefaultResolverValueResolver interface {
	// ResolveValue resolve the value to the proper type and return the value.
	// This method assumes target is never a reflect.Ptr, this should be handled before calling it.
	ResolveValue(target reflect.Value, value any) error
}

type DefaultResolverValueResolverCustomType interface {
	ResolveCustomTypeValue(target reflect.Value, value any) error
}

type DefaultResolverValueResolverCustomTypeReflect interface {
	ResolveCustomTypeValueReflect(target reflect.Value, sourceValue reflect.Value, value any) error
}

type DefaultResolver struct {
	valueResolver DefaultResolverValueResolver
}

func NewDefaultResolver(valueResolver DefaultResolverValueResolver) *DefaultResolver {
	if valueResolver == nil {
		valueResolver = &DefaultResolverValue{}
	}
	return &DefaultResolver{valueResolver: valueResolver}
}

func (r DefaultResolver) Resolve(target reflect.Value, value any) error {
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

type DefaultResolverValue struct {
	CustomTypes        []DefaultResolverValueResolverCustomType
	CustomTypesReflect []DefaultResolverValueResolverCustomTypeReflect
}

func (r DefaultResolverValue) ResolveValue(target reflect.Value, value any) error {
	if !target.CanSet() {
		return fmt.Errorf("cannot set '%s' ", target.Type().Kind())
	}

	if target.CanInterface() {
		for _, customType := range r.CustomTypes {
			err := customType.ResolveCustomTypeValue(target, value)
			if err == nil {
				return nil
			}
			if errors.Is(err, ErrCoerceUnknown) {
				continue
			}
			return err
		}
	}

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

	sourceValue := reflect.ValueOf(value)

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

	for _, customType := range r.CustomTypesReflect {
		err := customType.ResolveCustomTypeValueReflect(target, sourceValue, value)
		if err == nil {
			return nil
		}
		if errors.Is(err, ErrCoerceUnknown) {
			continue
		}
		return err
	}

	return fmt.Errorf("%w: cannot coerce source of type '%T' into target of type '%s'",
		ErrCoerceUnknown, value, target.Type().Kind())
}
