package instruct

import (
	"fmt"
	"reflect"
	"sync"
)

// MapTags is an alternative to struct tags, and can be used to override them.
type MapTags map[string]any

func (m MapTags) findPath(path []string) (any, bool) {
	return findPath(m, path)
}

func (m MapTags) findStringPath(path []string) (string, bool) {
	v, ok := m.findPath(path)
	if !ok {
		return "", false
	}

	if s, ok := v.(string); ok {
		return s, true
	}

	return "", false
}

func (m MapTags) checkUnusedFields(usedFields map[string]bool) error {
	mpKeys := map[string]bool{}
	buildMapTagList(m, level{}, mpKeys)

	for mkey, _ := range mpKeys {
		_, found := usedFields[mkey]
		if !found {
			return fmt.Errorf("map tags field '%s' was declared but not used", mkey)
		}
	}

	return nil
}

// mapTagsList is a thread-safe per-type list of MapTags.
type mapTagsList struct {
	list sync.Map
}

func (l *mapTagsList) Get(t reflect.Type) MapTags {
	if m, ok := l.list.Load(t); ok {
		return m.(MapTags)
	}
	return nil
}

func (l *mapTagsList) Exists(t reflect.Type) bool {
	_, ok := l.list.Load(t)
	return ok
}

func (l *mapTagsList) Set(t reflect.Type, m MapTags) {
	l.list.Store(reflectElem(t), m)
}

// getMapTags checks if the type is a MapTags-compatible map and returns it.
func getMapTags(v any) (MapTags, bool) {
	if mv, ok := v.(MapTags); ok {
		return mv, true
	}
	if mv, ok := v.(map[string]any); ok {
		return mv, true
	}
	return nil, false
}

// buildMapTagList builds a map with all MapTags keys to be used for comparison.
func buildMapTagList(m MapTags, lvl level, x map[string]bool) {
	for mkey, mval := range m {
		x[lvl.StringPathWithName(mkey)] = true
		if mv, ok := getMapTags(mval); ok {
			buildMapTagList(mv, lvl.Append(mkey), x)
		}
	}
}

// findPath recurse into MapTags to find the path.
func findPath(m MapTags, path []string) (any, bool) {
	if len(path) == 0 {
		return nil, false
	}
	v, ok := m[path[0]]
	if !ok {
		return nil, false
	}

	if len(path) > 1 {
		if mv, ok := getMapTags(v); ok {
			return findPath(mv, path[1:])
		}
		return nil, false
	}

	return v, true
}
