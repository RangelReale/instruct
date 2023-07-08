package instruct

import (
	"fmt"
	"reflect"
	"strings"
)

// structInfo caches struct field configurations.
type structInfo struct {
	typ    reflect.Type
	field  reflect.StructField
	tag    *Tag
	path   []string
	fields []*structInfo
}

func (s *structInfo) fullFieldName() string {
	return strings.Join(s.path, ".")
}

func (s *structInfo) checkSameType(t reflect.Type) error {
	t = reflectElem(t)
	if t != s.typ {
		return fmt.Errorf("invalid data type, expected %s got %s", s.typ.String(), t.String())
	}
	return nil
}

func (s *structInfo) findField(f reflect.StructField) *structInfo {
	for _, field := range s.fields {
		if intSliceEquals(field.field.Index, f.Index) {
			return field
		}
	}
	return nil
}

func (s *structInfo) fieldByName(name string) *structInfo {
	for _, field := range s.fields {
		if field.field.Name == name {
			return field
		}
	}
	return nil
}

// structInfoWithMapTags overrides a structInfo with a MapTags. This creates a clone of all the objects and don't
// change the original in any way.
func structInfoWithMapTags[IT any, DC DecodeContext](si *structInfo, mapTags MapTags, options DefaultOptions[IT, DC]) (*structInfo, error) {
	ctx := &buildContext{
		clone:            true,
		skipStructField:  true,
		skipStructOption: true,
	}

	newsi, err := buildStructInfoItem(ctx, si, level{}, mapTags, options)
	if err != nil {
		return nil, err
	}

	if mapTags != nil {
		err = mapTags.checkUnusedFields(ctx.GetUsedValues(mapTagValueKey))
		if err != nil {
			return nil, err
		}
	}

	return newsi, nil
}
