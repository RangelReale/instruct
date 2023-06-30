package instruct

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Resolver converts strings to the type of the struct field.
type Resolver interface {
	Resolve(t interface{}, v string) (interface{}, error)
}

// DefaultResolver calls DefaultResolve.
type DefaultResolver struct {
}

func (d DefaultResolver) Resolve(t interface{}, v string) (interface{}, error) {
	return DefaultResolve(t, v)
}

// resolveValues iterates over string values to resolve a slice value on the field
func resolveValues(resolver Resolver, field reflect.Value, typ reflect.Type, values []string) error {
	r := reflect.MakeSlice(typ, len(values), len(values))
	for i, value := range values {
		if err := resolveValue(resolver, r.Index(i), typ, value); err != nil {
			return err
		}
	}
	field.Set(reflect.ValueOf(r.Interface()))
	return nil
}

// resolveValue resolves and sets the string value to appropriate type on the field
func resolveValue(resolver Resolver, field reflect.Value, typ reflect.Type, value string) error {
	if field.Kind() == reflect.Pointer {
		v, err := resolver.Resolve(reflect.New(typ.Elem()).Elem().Interface(), value)
		if err != nil {
			return err
		}

		field.Set(reflect.New(typ.Elem()))
		field.Elem().Set(reflect.ValueOf(v))
		return nil
	}
	v, err := resolver.Resolve(field.Interface(), value)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(v))
	return nil
}

// DefaultResolve resolve the string value to the proper type and return the value.
func DefaultResolve(t interface{}, v string) (interface{}, error) {
	switch t.(type) {
	case string:
		return v, nil
	case bool:
		return strconv.ParseBool(v)
	case time.Time:
		return time.Parse(time.RFC3339, v)
	case time.Duration:
		return time.ParseDuration(v)
	case int:
		i, err := strconv.ParseInt(v, 10, 32)
		return int(i), err
	case int64:
		return strconv.ParseInt(v, 10, 64)
	case int32:
		i, err := strconv.ParseInt(v, 10, 32)
		return int32(i), err
	case int16:
		i, err := strconv.ParseInt(v, 10, 16)
		return int16(i), err
	case int8:
		i, err := strconv.ParseInt(v, 10, 8)
		return int8(i), err
	case float64:
		return strconv.ParseFloat(v, 64)
	case float32:
		i, err := strconv.ParseFloat(v, 32)
		return float32(i), err
	case uint:
		i, err := strconv.ParseUint(v, 10, 32)
		return uint(i), err
	case uint64:
		return strconv.ParseUint(v, 10, 64)
	case uint32:
		i, err := strconv.ParseUint(v, 10, 32)
		return uint32(i), err
	case uint16:
		i, err := strconv.ParseUint(v, 10, 16)
		return uint16(i), err
	case uint8:
		i, err := strconv.ParseUint(v, 10, 8)
		return uint8(i), err
	case complex128:
		return strconv.ParseComplex(v, 128)
	case complex64:
		i, err := strconv.ParseComplex(v, 64)
		return complex64(i), err
	default:
		return nil, fmt.Errorf("unsupported type: %v", reflect.TypeOf(t))
	}
}
