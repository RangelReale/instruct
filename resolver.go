package instruct

import "reflect"

// Resolver converts strings to the type of the struct field.
type Resolver interface {
	Resolve(target reflect.Value, value any) error
}
