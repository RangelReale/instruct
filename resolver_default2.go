package instruct

import (
	"fmt"
	"reflect"

	"github.com/RangelReale/instruct/coerce"
)

type DefaultResolver struct {
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
			sourceElemValue := sourceValue.Index(i)
			var targetValue reflect.Value
			if elemType.Kind() == reflect.Ptr {
				targetValue = reflect.New(elemType.Elem())
			} else {
				targetValue = reflect.New(elemType)
			}
			if err := r.Resolve(targetValue, sourceElemValue.Interface()); err != nil {
				return err
			}
			if elemType.Kind() == reflect.Ptr {
				targetSliceValue = reflect.Append(targetSliceValue, targetValue)
			} else {
				targetSliceValue = reflect.Append(targetSliceValue, reflect.Indirect(targetValue))
			}
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

	return DefaultResolveValue(target, value)
}

// DefaultResolveValue resolve the value to the proper type and return the value.
func DefaultResolveValue(target reflect.Value, value any) error {
	if !target.CanSet() {
		return fmt.Errorf("cannot set '%s' ", target.Type().Kind())
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
	case reflect.Interface:
		target.Set(reflect.ValueOf(value))
		return nil
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

	return fmt.Errorf("cannot coerce source of type '%T' into target of type '%s'",
		value, target.Type().Kind())
}
