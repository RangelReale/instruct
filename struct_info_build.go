package instruct

import (
	"fmt"
	"reflect"
)

func buildStructInfo[IT any, DC DecodeContext](t reflect.Type, mapTags MapTags, options DefaultOptions[IT, DC]) (*structInfo, error) {
	ctx := &buildContext{}

	t = reflectElem(t)

	// parse struct option if available
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
		// check unused defaultMapTags keys
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
				typ:   reflectElem(field.Type),
				field: field,
				path:  curlevel.Path(),
			}
		}

		// parse struct tag or equivalent map tag.
		if tag, err := parseStructTag(ctx, field, curlevel, mapTags, &options); err != nil {
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
