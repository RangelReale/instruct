package instruct

import (
	"reflect"
)

// TypeDecoder decodes inputs to structs for a specific type.
type TypeDecoder[IT any, DC DecodeContext, T any] struct {
	si      *structInfo
	decoder *Decoder[IT, DC]
	err     error
}

// NewTypeDecoder creates a TypeDecoder instance for a specific type without any decode operations. At least
// one must be added for decoding to work.
func NewTypeDecoder[IT any, DC DecodeContext, T any](options TypeDefaultOptions[IT, DC]) *TypeDecoder[IT, DC, T] {
	ret := &TypeDecoder[IT, DC, T]{
		decoder: NewDecoder[IT, DC](options.DefaultOptions),
	}

	var data T

	if options.MapTags != nil {
		ret.decoder.options.DefaultMapTagsSet(reflect.TypeOf(data), options.MapTags)
	}

	si, err := ret.decoder.structInfoFromType(reflect.TypeOf(data))
	if err != nil {
		ret.err = err
	} else {
		ret.si = si
	}

	return ret
}

// Decode decodes the input to the generic type.
func (d *TypeDecoder[IT, DC, T]) Decode(input IT, decodeOptions DecodeOptions[IT, DC]) (T, error) {
	var data T

	if d.err != nil {
		return data, d.err
	}

	// creates a new instance of the type.
	data = decodeTypeNew[T]()

	// decodes using the cached struct info.
	err := d.decoder.decodeInputFromStructInfo(input, d.si, &data, decodeOptions)
	return data, err
}

// decodeTypeNew creates a new value of the type, initializing a pointer if needed.
func decodeTypeNew[T any]() T {
	var v T
	if typ := reflect.TypeOf(v); typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem()).Interface().(T)
	}
	return v
}
