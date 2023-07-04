package options

import "github.com/RangelReale/instruct"

// ExtractOptions extracts only options of a specific type.
func ExtractOptions[IT any, DC instruct.DecodeContext, T Option[IT, DC]](options []Option[IT, DC]) []T {
	var ret []T
	for _, opt := range options {
		if o, ok := opt.(T); ok {
			ret = append(ret, o)
		}
	}
	return ret
}

// ConcatOptionsBefore returns an array with "options" before "source".
func ConcatOptionsBefore[IT any, DC instruct.DecodeContext, T Option[IT, DC]](source []T, options ...T) []T {
	return append(append([]T{}, options...), source...)
}
