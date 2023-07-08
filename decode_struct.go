package instruct

import (
	"fmt"
	"reflect"

	"github.com/RangelReale/instruct/types"
)

// decodeStruct uses the structInfo to decode the input to the struct.
func (d *Decoder[IT, DC]) decodeStruct(si *structInfo, input IT, dataValue reflect.Value, decodeOptions DecodeOptions[IT, DC]) error {
	dataValue = reflectValueElem(dataValue)
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
		fieldValue := dataValue.FieldByIndex(sifield.field.Index)

		dataWasSet := false

		switch sifield.tag.Operation {
		case OperationIgnore: // ignore
			dataWasSet = true
		case OperationRecurse:
			// recurse into inner struct
			reflectEnsurePointerValue(&fieldValue)
			if err := d.decodeStruct(sifield, input, fieldValue, decodeOptions); err != nil {
				return err
			}
			dataWasSet = true
		default:
			var err error
			// execute operation (query, header, etc.)
			dataWasSet, err = d.executeOperation(fieldValue, sifield, input, decodeOptions)
			if err != nil {
				return err
			}
		}

		if !dataWasSet && sifield.tag.Required {
			return types.RequiredError{
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

		return types.RequiredError{
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
			return false, types.OperationNotSupportedError{
				Operation: sifield.tag.Operation,
				FieldName: structFieldName(sifield.typ, sifield.fullFieldName()),
			}
		}

		if err = d.options.Resolver.Resolve(field, value); err != nil {
			return false, fmt.Errorf("error resolving field '%s': %w", sifield.fullFieldName(), err)
		}
	}

	return dataWasSet, nil
}
