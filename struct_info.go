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

func buildStructInfo[IT any, DC DecodeContext](t reflect.Type, mapTags MapTags, options DefaultOptions[IT, DC]) (*structInfo, error) {
	ctx := &buildContext{}

	t = reflectElem(t)

	// parse struct tags if available
	tag, err := parseStructTagStructOption(ctx, t, level{}, mapTags, &options)
	if err != nil {
		return nil, err
	}

	// build struct info for root struct.
	si, err := buildStructInfoItem(ctx, &structInfo{
		typ: t,
		tag: tag,
	}, level{}, mapTags, options)
	if err != nil {
		return nil, err
	}

	if mapTags != nil {
		// check unused MapTags keys
		err = mapTags.checkUnusedFields(ctx.GetUsedValues(mapTagValueKey))
		if err != nil {
			return nil, err
		}
	}

	return si, err
}

// buildStructInfoItem builds a structInfo for the fields of the passed struct.
// This function is used both to create a new structInfo and override one with new MapTags. In the former case,
// it returns copies of the fields and don't change the original.
func buildStructInfoItem[IT any, DC DecodeContext](ctx *buildContext, si *structInfo, lvl level, mapTags MapTags, options DefaultOptions[IT, DC]) (*structInfo, error) {
	siBuild := buildCloneStructInfo(ctx, si, false)

	// // try to find option field first
	// if siTag, err := structInfoFindOptionsField(ctx, siBuild.typ, lvl, mapTags, options); err == nil && siTag != nil {
	// 	siBuild.tag = siTag
	// 	siBuild.hasStructOption = true
	// } else if err != nil {
	// 	return nil, err
	// }

	if siBuild.tag != nil && siBuild.tag.IsSO {
		// check struct option
		if siBuild.tag.Operation == OperationIgnore {
			return nil, fmt.Errorf("cannot ignore struct option for field '%s'", lvl.StringPath())
		}

		// if not recursing, skip checking fields
		if !siBuild.tag.SORecurse {
			return siBuild, nil
		}
	}

	for i := 0; i < siBuild.typ.NumField(); i++ {
		field := siBuild.typ.Field(i)
		if !field.IsExported() || isOptionField(field) {
			continue
		}

		curlevel := lvl.Append(field.Name)
		sifield := si.findField(field)
		if sifield != nil {
			sifield = buildCloneStructInfo(ctx, sifield, true)
		} else {
			sifield = &structInfo{
				typ:   field.Type,
				field: field,
				path:  curlevel.Path(),
			}
		}

		// parse struct tag or equivalent map tag.
		if tag, err := parseStructTag(ctx, &si.typ, &field, curlevel, mapTags, &options); err != nil {
			return nil, fmt.Errorf("error on field '%s': %w", curlevel.StringPath(), err)
		} else if tag != nil {
			sifield.tag = tag
		}

		if sifield.tag == nil {
			return nil, fmt.Errorf("field '%s' configuration not found", curlevel.StringPath())
		}

		if sifield.tag.Operation == OperationRecurse {
			// recurse into inner struct
			if !isStruct(field.Type) {
				return nil, fmt.Errorf("field '%s' must be a struct to use recurse but is '%s'", field.Name, field.Type.String())
			}
			var err error
			sifield, err = buildStructInfoItem(ctx, sifield, lvl.AppendIfTrue(!field.Anonymous, field.Name), mapTags, options)
			if err != nil {
				return nil, err
			}
		}

		siBuild.fields = append(siBuild.fields, sifield)
	}

	return siBuild, nil
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

// buildCloneStructInfo clones a structInfo if ctx.clone is true, otherwise just return the original one.
func buildCloneStructInfo(ctx *buildContext, si *structInfo, withFields bool) *structInfo {
	if !ctx.clone {
		return si
	}
	ret := &structInfo{
		typ:   si.typ,
		field: si.field,
		tag:   si.tag,
		path:  si.path,
	}
	if withFields {
		ret.fields = si.fields
	}
	return ret
}
