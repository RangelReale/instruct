package instruct

import (
	"reflect"
	"strings"

	"github.com/RangelReale/instruct/resolver"
)

// FieldNameMapper maps a struct field name to the target field name.
// The default one uses [strings.ToLower].
type FieldNameMapper func(operation string, name string) string

type DefaultOptions[IT any, DC DecodeContext] struct {
	TagName            string                             // struct tag name. Default "inreq".
	DefaultRequired    bool                               // whether the default for fields should be "required" or "not required"
	DecodeOperations   map[string]DecodeOperation[IT, DC] // list of decode operations
	defaultMapTags     *mapTagsList                       // list of DEFAULT map tags
	FieldNameMapper    FieldNameMapper                    // field name mapper. Default one uses [strings.ToLower].
	structInfoProvider structInfoProvider[IT, DC]         // allows caching of structInfo
	Resolver           Resolver                           // interface used to convert strings to the struct field type.
}

func (o *DefaultOptions[IT, DC]) DefaultMapTagsSet(t reflect.Type, m MapTags) {
	t = reflectElem(t)
	o.defaultMapTags.Set(t, m)
	o.structInfoProvider.remove(t)
}

func (o *DefaultOptions[IT, DC]) StructInfoCache(cache bool) {
	if cache {
		o.structInfoProvider = &CachedStructInfoProvider[IT, DC]{}
	} else {
		o.structInfoProvider = &DefaultStructInfoProvider[IT, DC]{}
	}
}

type DecodeOptions[IT any, DC DecodeContext] struct {
	Ctx                       DC      // decode context to be sent to DecodeOperation.
	MapTags                   MapTags // decode call-specific MapTags. They may override existing ones.
	UseDecodeMapTagsAsDefault bool    // internal flag to allow Decode functions without an instance to set MapTags as a default one.
}

// NewDefaultOptions returns a DefaultOptions with the default values.
func NewDefaultOptions[IT any, DC DecodeContext]() DefaultOptions[IT, DC] {
	return DefaultOptions[IT, DC]{
		TagName:            "instruct",
		DefaultRequired:    true,
		DecodeOperations:   map[string]DecodeOperation[IT, DC]{},
		defaultMapTags:     &mapTagsList{},
		FieldNameMapper:    DefaultFieldNameMapper,
		structInfoProvider: DefaultStructInfoProvider[IT, DC]{},
		Resolver:           resolver.NewResolver(nil),
	}
}

// NewDecodeOptions returns a DecodeOptions with the default values.
func NewDecodeOptions[IT any, DC DecodeContext]() DecodeOptions[IT, DC] {
	return DecodeOptions[IT, DC]{}
}

// helpers

// DefaultFieldNameMapper converts names to lowercase using [strings.ToLower].
func DefaultFieldNameMapper(operation string, name string) string {
	return strings.ToLower(name)
}
