package instruct

import (
	"reflect"
	"strings"
)

// reflectElem returns the first non-pointer type from the [reflect.Type].
func reflectElem(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// reflectValueElem returns the first non-pointer type from the [reflect.Value].
func reflectValueElem(t reflect.Value) reflect.Value {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func structFieldName(structTyp reflect.Type, fieldName string) string {
	if fieldName == "" {
		return structTyp.String()
	}
	return fieldName
}

func isZero[T any](v T) bool {
	return reflect.ValueOf(&v).Elem().IsZero()
}

// intSliceEquals returns true if the arrays are equal (in the same order).
func intSliceEquals(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// isStruct checks if the [reflect.Type] is a struct
func isStruct(t reflect.Type) bool {
	return reflectElem(t).Kind() == reflect.Struct
}

// level is a helper for building paths.
// It is immutable, all methods return a new copy.
type level struct {
	names []string
}

func (s level) Append(name string) level {
	return level{
		names: append(append([]string{}, s.names...), name),
	}
}

func (s level) AppendIfTrue(cond bool, name string) level {
	if !cond {
		return s
	}
	return s.Append(name)
}

func (s level) StringPath() string {
	if len(s.names) == 0 {
		return ""
	}
	return strings.Join(s.names, ".")
}

func (s level) StringPathWithName(name string) string {
	if len(s.names) == 0 {
		return name
	}
	return strings.Join(s.names, ".") + "." + name
}

func (s level) Path() []string {
	return s.names
}

func (s level) PathWithName(name string) []string {
	return append(append([]string{}, s.names...), name)
}
