package instruct

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// structInfo caches struct field configurations.
type structInfo struct {
	typ    reflect.Type        // non-pointer final type (from "field.Type()")
	field  reflect.StructField // field, empty on root struct
	tag    *Tag                // tag
	path   []string            // complete field path including itself, using the unmodified struct field name
	fields []*structInfo       // child fields
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

func (s *structInfo) dump(w io.Writer) error {
	return s.dumpIndent("", w)
}

func (s *structInfo) dumpIndent(indent string, w io.Writer) error {
	var ferr error
	var err error
	_, err = fmt.Fprintf(w, "%s- ", indent)
	ferr = errors.Join(ferr, err)
	name := "{ROOT}"
	if s.field.Type != nil {
		name = s.field.Name
	}
	_, err = fmt.Fprintf(w, "%s [/%s] (type: %s, kind: %s) ", name, strings.Join(s.path, "/"), s.typ.Name(), s.typ.Kind().String())
	ferr = errors.Join(ferr, err)
	if s.field.Type != nil {
		_, err = fmt.Fprintf(w, "[field type: %s, kind: %s] ", s.field.Type.String(), s.field.Type.Kind().String())
		ferr = errors.Join(ferr, err)
		// _, err = fmt.Fprintf(w, "[path: '%s']", strings.Join(s.path, "/"))
		// ferr = errors.Join(ferr, err)
	}
	_, err = fmt.Fprintf(w, "\n")
	ferr = errors.Join(ferr, err)
	for _, field := range s.fields {
		ferr = errors.Join(ferr, field.dumpIndent(indent+"\t", w))
	}
	return err
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
