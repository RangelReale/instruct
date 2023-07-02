package instruct

import (
	"errors"
	"reflect"
)

// Decoder decodes inputs to structs.
type Decoder[IT any, DC DecodeContext] struct {
	options DefaultOptions[IT, DC]
}

// NewDecoder creates a Decoder instance without any decode operations. At least one must be added for
// decoding to work.
func NewDecoder[IT any, DC DecodeContext](options DefaultOptions[IT, DC]) *Decoder[IT, DC] {
	ret := &Decoder[IT, DC]{
		options: options,
	}
	return ret
}

// Decode decodes the input to the struct passed in "data".
func (d *Decoder[IT, DC]) Decode(input IT, data any, decodeOptions DecodeOptions[IT, DC]) error {
	if isZero(decodeOptions.Ctx) {
		return errors.New("decode context cannot be nil")
	}

	if decodeOptions.UseDecodeMapTagsAsDefault && decodeOptions.MapTags != nil {
		// helper option to automatically set map tags as default, meant for free-standind "Decode" functions only.
		d.options.defaultMapTags.Set(reflectElem(reflect.TypeOf(data)), decodeOptions.MapTags)
		decodeOptions.MapTags = nil
	}

	return d.decodeInput(input, data, decodeOptions)
}
