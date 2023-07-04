package options

// ExtractOptions extracts only options of a specific type.
func ExtractOptions[RET Option, IN ~[]Option](options IN) []RET {
	var ret []RET
	for _, opt := range options {
		if o, ok := opt.(RET); ok {
			ret = append(ret, o)
		}
	}
	return ret
}

// ConcatOptionsBefore returns an array with "options" before "source".
func ConcatOptionsBefore[T Option](source []T, options ...T) []T {
	return append(append([]T{}, options...), source...)
}
