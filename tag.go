package instruct

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Tag contains the options parsed from the struct tags or MapTags
type Tag struct {
	Operation string     // decode operation
	Name      string     // data name (for example, header or query param name)
	Required  bool       // whether this field is required to be set
	Options   TagOptions // options
	IsSO      bool
	SOWhen    string // struct options: when to parse (before or after the fields)
	SORecurse bool   // struct options: whether to recurse into inner struct
}

type TagOptions struct {
	options map[string]string
}

func NewTagOptions() TagOptions {
	return TagOptions{
		options: map[string]string{},
	}
}

func (t *TagOptions) Exists(name string) bool {
	_, ok := t.options[name]
	return ok
}

func (t *TagOptions) Get(name string) (string, bool) {
	v, ok := t.options[name]
	return v, ok
}

func (t *TagOptions) Value(name string, defaultValue string) string {
	if tv, ok := t.options[name]; ok {
		return tv
	}
	return defaultValue
}

func (t *TagOptions) BoolValue(name string, defaultValue bool) (bool, error) {
	if tv, ok := t.options[name]; ok {
		b, err := strconv.ParseBool(tv)
		if err != nil {
			return false, err
		}
		return b, nil
	}
	return defaultValue, nil
}

// parseStructTagStructField parses a Tag from a struct tag
func parseStructTagStructField[IT any, DC DecodeContext](ctx *buildContext, field reflect.StructField, level level,
	options *DefaultOptions[IT, DC]) (*Tag, error) {
	tags, ok := field.Tag.Lookup(options.TagName)
	if ok {
		return parseTags(field.Name, tags, options)
	}

	if field.Anonymous && isStruct(field.Type) {
		return &Tag{
			Operation: OperationRecurse,
			Required:  options.DefaultRequired,
			Name:      options.FieldNameMapper(OperationRecurse, field.Name),
		}, nil
	}

	return nil, nil
}

// parseStructTagStructField parses a Tag from MapTags
func parseStructTagMapTags[IT any, DC DecodeContext](ctx *buildContext, field reflect.StructField, level level, mapTags MapTags,
	options *DefaultOptions[IT, DC]) (*Tag, error) {
	if mapTags == nil {
		return nil, nil
	}

	if ft, ok := mapTags.findPath(level.Path()); ok {
		ctx.ValueUsed(mapTagValueKey, level.StringPath())

		switch xft := ft.(type) {
		case string:
			return parseTags(field.Name, xft, options)
		case MapTags, map[string]any:
			if isStruct(field.Type) {
				return &Tag{
					Operation: OperationRecurse,
					Required:  options.DefaultRequired,
					Name:      options.FieldNameMapper(OperationRecurse, field.Name),
				}, nil
			}
		}

		return nil, fmt.Errorf("unknown map tags item type (only 'string', 'MapTags' and 'map[string]any' are allowed): %T", ft)
	}

	return nil, nil
}

// parseStructTagStructOption finds a struct option field from either the MapTags or the struct fields.
func parseStructTagStructOption[IT any, DC DecodeContext](ctx *buildContext, t reflect.Type, lvl level,
	mapTags MapTags, options *DefaultOptions[IT, DC]) (*Tag, error) {
	if !ctx.skipMapTags {
		tag, err := structInfoFindOptionsFieldMapTags(ctx, t, lvl, mapTags, options)
		if err != nil {
			return nil, err
		}
		if tag != nil {
			tag.IsSO = true
			return tag, nil
		}
	}

	if !ctx.skipStructField {
		tag, err := structInfoFindOptionsFieldStructField(ctx, t, lvl, mapTags, options)
		if err != nil {
			return nil, err
		}
		if tag != nil {
			tag.IsSO = true
			return tag, nil
		}
	}

	return nil, nil
}

// parseStructTagStructField parses a Tag from a MapTags or a struct tag.
func parseStructTag[IT any, DC DecodeContext](ctx *buildContext, structType *reflect.Type, field *reflect.StructField,
	lvl level, mapTags MapTags,
	options *DefaultOptions[IT, DC]) (*Tag, error) {
	if field != nil && !ctx.skipMapTags {
		tag, err := parseStructTagMapTags(ctx, *field, lvl, mapTags, options)
		if err != nil {
			return nil, err
		}
		if tag != nil {
			return tag, nil
		}
	}
	if field != nil && !ctx.skipStructField {
		tag, err := parseStructTagStructField(ctx, *field, lvl, options)
		if err != nil {
			return nil, err
		}
		if tag != nil {
			return tag, nil
		}
	}
	if structType != nil && !ctx.skipStructOption {
		tag, err := parseStructTagStructOption(ctx, *structType, lvl, mapTags, options)
		if err != nil {
			return nil, err
		}
		if tag != nil {
			return tag, nil
		}
	}

	return nil, nil
}

// parseTags parses a Tag from a textual description, in the form "operation,field1=value1,field2=value2".
func parseTags[IT any, DC DecodeContext](fieldName string, tagValue string, options *DefaultOptions[IT, DC]) (*Tag, error) {
	ret := &Tag{
		Name:     "",
		Required: options.DefaultRequired,
		Options:  NewTagOptions(),
	}

	s := tagValue
	for s != "" {
		var value string
		value, s, _ = strings.Cut(s, ",")
		if ret.Operation == "" {
			if value == "" {
				return nil, errors.New("operation cannot be blank")
			}
			ret.Operation = value
		} else if value != "" {
			oname, oval, ofound := strings.Cut(value, "=")
			if !ofound {
				return nil, fmt.Errorf("unnamed tag option: %s", value)
			}
			if oname == "name" {
				ret.Name = oval
			} else if oname == "required" {
				b, err := strconv.ParseBool(oval)
				if err != nil {
					return nil, fmt.Errorf("error parsing 'required' boolean option: %w", err)
				}
				ret.Required = b
			} else if strings.HasPrefix(oname, "so_") {
				if oname == "so_when" {
					if oval != SOOptionWhenBefore && oval != SOOptionWhenAfter {
						return nil, fmt.Errorf("invalid 'when' option value: %s", oval)
					}
					ret.SOWhen = oval
				} else if oname == "so_recurse" {
					b, err := strconv.ParseBool(oval)
					if err != nil {
						return nil, fmt.Errorf("error parsing 'so_recurse' boolean option: %w", err)
					}
					ret.SORecurse = b
				} else {
					return nil, fmt.Errorf("unknown struct option name: %s", oname)
				}
			} else {
				ret.Options.options[oname] = oval
			}
		}
	}

	if ret.Name == "" {
		ret.Name = options.FieldNameMapper(ret.Operation, fieldName)
	}
	return ret, nil
}
