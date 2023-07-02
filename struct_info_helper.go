package instruct

import (
	"fmt"
	"reflect"
	"sync"
)

type buildContext struct {
	usedValues       map[string]map[string]bool
	clone            bool
	skipMapTags      bool
	skipStructField  bool
	skipStructOption bool
}

func (d *buildContext) ValueUsed(operation string, name string) {
	if d.usedValues == nil {
		d.usedValues = map[string]map[string]bool{}
	}
	if _, ok := d.usedValues[operation]; !ok {
		d.usedValues[operation] = map[string]bool{}
	}
	d.usedValues[operation][name] = true
}

func (d *buildContext) GetUsedValues(operation string) map[string]bool {
	if d.usedValues == nil {
		return nil
	}
	operationValues, ok := d.usedValues[operation]
	if !ok {
		return nil
	}

	return operationValues
}

// structInfoFindOptionsFieldStructField finds a struct option field inside the struct fields.
func structInfoFindOptionsFieldStructField[IT any, DC DecodeContext](ctx *buildContext, t reflect.Type, lvl level,
	mapTags MapTags, options *DefaultOptions[IT, DC]) (*Tag, error) {
	var tag *Tag

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if isOptionField(field) {
			if tag != nil {
				return nil, fmt.Errorf("only one StructOption is allowed per struct for field '%s'", lvl.StringPath())
			}
			var err error
			tag, err = parseStructTag(ctx, nil, &field, lvl, mapTags, options)
			if err != nil {
				return nil, fmt.Errorf("error on field '%s': %w", lvl.StringPath(), err)
			}
		}
	}

	return tag, nil
}

// structInfoFindOptionsFieldMapTags finds a struct option field from the MapTags.
func structInfoFindOptionsFieldMapTags[IT any, DC DecodeContext](ctx *buildContext, t reflect.Type, lvl level,
	mapTags MapTags, options *DefaultOptions[IT, DC]) (*Tag, error) {
	if mapTags != nil {
		if _, ok := mapTags.findStringPath(lvl.Append(StructOptionMapTag).Path()); ok {
			return parseStructTag(ctx, nil, &reflect.StructField{Name: StructOptionMapTag},
				lvl.Append(StructOptionMapTag), mapTags, options)
		}
	}
	return nil, nil
}

// structInfoProvider abstracts a posssible cache of structInfo
type structInfoProvider[IT any, DC DecodeContext] interface {
	provide(t reflect.Type, mapTags MapTags, options DefaultOptions[IT, DC]) (*structInfo, error)
}

// DefaultStructInfoProvider is a structInfoProvider that never caches.
type DefaultStructInfoProvider[IT any, DC DecodeContext] struct {
}

func (d DefaultStructInfoProvider[IT, DC]) provide(t reflect.Type, mapTags MapTags, options DefaultOptions[IT, DC]) (*structInfo, error) {
	return buildStructInfo(t, mapTags, options)
}

// CachedStructInfoProvider is a structInfoProvider that always caches.
type CachedStructInfoProvider[IT any, DC DecodeContext] struct {
	cache sync.Map
}

func (d *CachedStructInfoProvider[IT, DC]) provide(t reflect.Type, mapTags MapTags, options DefaultOptions[IT, DC]) (*structInfo, error) {
	csi, ok := d.cache.Load(t)
	if ok {
		return csi.(*structInfo), nil
	}

	si, err := buildStructInfo(t, mapTags, options)
	if err != nil {
		return nil, err
	}

	d.cache.Store(t, si)

	return si, nil
}
