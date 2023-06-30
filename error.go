package instruct

import (
	"fmt"
	"reflect"
)

// An ValuesNotUsedError is returned when some values were not used.
type ValuesNotUsedError struct {
	Operation string
}

func (e ValuesNotUsedError) Error() string {
	return fmt.Sprintf("some values were not used on operation '%s'", e.Operation)
}

// An InvalidDecodeError describes an invalid argument passed to Decode.
// (The argument to Decode must be a non-nil pointer.)
type InvalidDecodeError struct {
	Type reflect.Type
}

func (e *InvalidDecodeError) Error() string {
	if e.Type == nil {
		return "error: Decode(nil)"
	}

	if e.Type.Kind() != reflect.Pointer {
		return "error: Decode(non-pointer " + e.Type.String() + ")"
	}
	return "error: Decode(nil " + e.Type.String() + ")"
}
