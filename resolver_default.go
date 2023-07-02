package instruct

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/RangelReale/instruct/coerce"
)

type DefaultResolver struct {
}

func (r DefaultResolver) Resolve(target reflect.Value, value any) error {
	return DefaultResolve(target, value)
}

func DefaultResolve(target reflect.Value, value any) error {
	if !target.CanSet() {
		return errors.New("not can set")
	}

	if value == nil {
		target.SetZero()
		return nil
	}

	sourceValue := reflect.ValueOf(value)

	if tryAssign(sourceValue.Type(), target.Type(), sourceValue, target) {
		return nil
	}

	switch target.Kind() {
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

	return fmt.Errorf("cannot coerce source of type '%s' into target of type '%s'",
		sourceValue.Type().Kind(), target.Type().Kind())
}

func tryAssign(st, tt reflect.Type, sv, tv reflect.Value) bool {
	st = unpointType(st)
	tt = unpointType(tt)
	sv = unpointValue(sv)
	tv = unpointValue(tv)

	if !tv.CanSet() {
		return false
	}

	if tt.AssignableTo(st) {
		// the source can be directly assigned to the target
		tv.Set(sv)
		return true
	}

	if st.ConvertibleTo(tt) {
		tv.Set(sv.Convert(tt))
		return true
	}

	return false
}
