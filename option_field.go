package instruct

import "reflect"

// StructOption can be used as a struct field to give options to the struct itself.
type StructOption struct{}

// StructOptionMapTag corresponds to StructOption in a MapTags.
const StructOptionMapTag = "_"

func isOptionField(f reflect.StructField) bool {
	return !f.IsExported() && reflect.TypeOf(StructOption{}) == f.Type
}
