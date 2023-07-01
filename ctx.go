package instruct

// DecodeContext is the context sent to DecodeOperation.
type DecodeContext interface {
	// ValueUsed signals that the value was used.
	ValueUsed(operation string, name string)
	// GetUsedValues returns the list of used values for the operation.
	GetUsedValues(operation string) map[string]bool
	// FieldNameMapper returns the FieldNameMapper instance to be used for converting the struct field name.
	FieldNameMapper() FieldNameMapper
}

// DefaultDecodeContext implements the standard decode context.
type DefaultDecodeContext struct {
	fieldNameMapper FieldNameMapper
	usedValues      map[string]map[string]bool
}

// NewDefaultDecodeContext creates an instance of DefaultDecodeContext.
func NewDefaultDecodeContext(fnMapper FieldNameMapper) DefaultDecodeContext {
	return DefaultDecodeContext{
		fieldNameMapper: fnMapper,
	}
}

func (d *DefaultDecodeContext) ValueUsed(operation string, name string) {
	if d.usedValues == nil {
		d.usedValues = map[string]map[string]bool{}
	}
	if _, ok := d.usedValues[operation]; !ok {
		d.usedValues[operation] = map[string]bool{}
	}
	d.usedValues[operation][name] = true
}

func (d *DefaultDecodeContext) GetUsedValues(operation string) map[string]bool {
	if d.usedValues == nil {
		return nil
	}
	operationValues, ok := d.usedValues[operation]
	if !ok {
		return nil
	}

	return operationValues
}

func (d *DefaultDecodeContext) FieldNameMapper() FieldNameMapper {
	return d.fieldNameMapper
}
