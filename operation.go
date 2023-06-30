package instruct

import (
	"reflect"
)

// Default operations.
const (
	OperationIgnore  string = "-"
	OperationRecurse        = "recurse"
)

const (
	mapTagValueKey string = "map-tag"
)

// StructOption "when" values.
const (
	SOOptionWhenBefore  = "before"
	SOOptionWhenAfter   = "after"
	SOOptionWhenDefault = SOOptionWhenAfter
)

func soOptionValue(soOption string) string {
	if soOption == "" {
		soOption = SOOptionWhenDefault
	}
	return soOption
}

// DecodeOperation is the interface for the input-to-struct decoders.
type DecodeOperation[IT any, DC DecodeContext] interface {
	Decode(ctx DC, input IT, field reflect.Value, typ reflect.Type, tag *Tag) (found bool, value any, err error)
}

// DecodeOperationValidate allows a DecodeOperation to do a final validation.
type DecodeOperationValidate[IT any, DC DecodeContext] interface {
	Validate(ctx DC, input IT) error
}

// DecodeOperationFunc wraps a DecodeOperation as a function.
type DecodeOperationFunc[IT any, DC DecodeContext] func(ctx DC, input IT, field reflect.Value, typ reflect.Type, tag *Tag) (bool, any, error)

func (f DecodeOperationFunc[IT, DC]) Decode(ctx DC, input IT, field reflect.Value, typ reflect.Type, tag *Tag) (bool, any, error) {
	return f(ctx, input, field, typ, tag)
}

// IgnoreDecodeValue can be returned from [DecodeOperation.Decode] to signal that the value should not be set on the
// struct field. This is used for example in the "body" decoder.
var IgnoreDecodeValue = struct{}{}
