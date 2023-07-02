package instruct

import (
	"fmt"
	"reflect"
)

// decodeStruct uses the structInfo to decode the input to the struct.
func (d *Decoder[IT, DC]) decodeStruct(si *structInfo, input IT, data interface{}, decodeOptions DecodeOptions[IT, DC]) error {
	dataValue := reflectValueElem(reflect.ValueOf(data))
	if err := si.checkSameType(dataValue.Type()); err != nil {
		return err
	}

	// execute the struct operation (using StructOption or inner struct tags). Only executed if "when" is
	// configured as "before".
	err := d.executeStructOperation(SOOptionWhenBefore, dataValue, si, input, decodeOptions)
	if err != nil {
		return err
	}

	for _, sifield := range si.fields {
		field := dataValue.FieldByIndex(sifield.field.Index)

		dataWasSet := false

		switch sifield.tag.Operation {
		case OperationIgnore: // ignore
			dataWasSet = true
		case OperationRecurse:
			// recurse into inner struct
			if err := d.decodeStruct(sifield, input, field.Addr().Interface(), decodeOptions); err != nil {
				return err
			}
			dataWasSet = true
		default:
			var err error
			// execute operation (query, header, etc.)
			dataWasSet, err = d.executeOperation(field, sifield, input, decodeOptions)
			if err != nil {
				return err
			}
		}

		if !dataWasSet && sifield.tag.Required {
			return RequiredError{
				Operation: sifield.tag.Operation,
				FieldName: sifield.fullFieldName(),
				TagName:   sifield.tag.Name,
			}
		}
	}

	// execute the struct operation (using StructOption or inner struct tags). Only executed if "when" is
	// configured as "after".
	err = d.executeStructOperation(SOOptionWhenAfter, dataValue, si, input, decodeOptions)
	if err != nil {
		return err
	}

	return nil
}

// executeStructOperation execute the struct operation (using StructOption or inner struct tags).
func (d *Decoder[IT, DC]) executeStructOperation(when string, dataValue reflect.Value, si *structInfo,
	input IT, decodeOptions DecodeOptions[IT, DC]) error {
	if si.tag == nil || !si.tag.IsSO || soOptionValue(si.tag.SOWhen) != when {
		return nil
	}

	dataWasSet, err := d.executeOperation(dataValue, si, input, decodeOptions)
	if err != nil {
		return err
	}
	if !dataWasSet && si.tag.Required {
		fn := si.fullFieldName()
		if fn == "" {
			fn = si.typ.String()
		}

		return RequiredError{
			IsStructOption: true,
			Operation:      si.tag.Operation,
			FieldName:      structFieldName(si.typ, si.fullFieldName()),
			TagName:        si.tag.Name,
		}
	}
	return nil
}

// executeOperation executes an operation (query, header, etc) on a struct field.
// If the decode interface return IgnoreDecodeValue, the value is not set to it.
func (d *Decoder[IT, DC]) executeOperation(field reflect.Value, sifield *structInfo, input IT,
	decodeOptions DecodeOptions[IT, DC]) (bool, error) {
	// check if the operation exists
	operation, opok := d.options.DecodeOperations[sifield.tag.Operation]
	if !opok {
		return false, fmt.Errorf("unknown operation '%s' for field %s", sifield.tag.Operation, sifield.field.Name)
	}

	// call the decoder interface.
	dataWasSet, value, err := operation.Decode(decodeOptions.Ctx, input, field, sifield.field.Type, sifield.tag)
	if err != nil {
		return false, err
	}

	if dataWasSet && value != IgnoreDecodeValue {
		if sifield.field.Type == nil {
			// struct option can't be set as a value
			return false, OperationNotSupportedError{
				Operation: sifield.tag.Operation,
				FieldName: structFieldName(sifield.typ, sifield.fullFieldName()),
			}
		}

		if err = d.options.Resolver.Resolve(field, value); err != nil {
			return false, fmt.Errorf("error resolving field '%s': %w", sifield.fullFieldName(), err)
		}

		// // convert the value from string/[]string to the struct field type.
		// switch xvalue := value.(type) {
		// case string:
		// 	if err = resolveValue(d.options.Resolver, field, sifield.field.Type, xvalue); err != nil {
		// 		return false, err
		// 	}
		// case []string:
		// 	if err = resolveValues(d.options.Resolver, field, sifield.field.Type, xvalue); err != nil {
		// 		return false, err
		// 	}
		// default:
		// 	return false, fmt.Errorf("unknown decoded value type '%T' for field '%s'", value, sifield.field.Name)
		// }
	}

	return dataWasSet, nil
}
