package instruct

import (
	"errors"
	"fmt"
	"reflect"
)

func (d *Decoder[IT, DC]) decodeInput(input IT, data any, decodeOptions DecodeOptions[IT, DC]) error {
	return d.decodeInputFromType(input, reflect.TypeOf(data), data, decodeOptions)
}

func (d *Decoder[IT, DC]) decodeInputFromType(input IT, typ reflect.Type, data any, decodeOptions DecodeOptions[IT, DC]) error {
	si, err := d.structInfoFromType(typ)
	if err != nil {
		return err
	}
	return d.decodeInputFromStructInfo(input, si, data, decodeOptions)
}

func (d *Decoder[IT, DC]) decodeInputFromStructInfo(input IT, si *structInfo, data any, decodeOptions DecodeOptions[IT, DC]) error {
	if isZero(decodeOptions.Ctx) {
		return errors.New("decode context cannot be nil")
	}

	rv := reflect.ValueOf(data)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidDecodeError{reflect.TypeOf(data)}
	}

	if decodeOptions.MapTags != nil {
		var err error

		// if there is a decode-specific map tag, create a new struct info based on the default one.
		si, err = structInfoWithMapTags(si, decodeOptions.MapTags, d.options)
		if err != nil {
			return err
		}
	}

	err := d.decodeStruct(si, input, data, decodeOptions)
	if err != nil {
		return err
	}

	// validate decode operations
	for _, operation := range d.options.DecodeOperations {
		if v, ok := operation.(DecodeOperationValidate[IT, DC]); ok {
			err = v.Validate(decodeOptions.Ctx, input)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// structInfoFromType builds an structInfo from a [reflect.Type], using the decoder options.
func (d *Decoder[IT, DC]) structInfoFromType(typ reflect.Type) (*structInfo, error) {
	if typ == nil {
		return nil, fmt.Errorf("cannot decode to nil")
	}
	typ = reflectElem(typ)
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("can only decode to struct, received: %v", typ.Kind())
	}

	return d.options.structInfoProvider.provide(typ, d.options.defaultMapTags.Get(typ), d.options)
}
