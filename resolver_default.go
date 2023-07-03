package instruct

// import (
// 	"fmt"
// 	"reflect"
//
// 	"github.com/RangelReale/instruct/coerce"
// )
//
// type DefaultResolver struct {
// }
//
// func (r DefaultResolver) Resolve(target reflect.Value, value any) error {
// 	return DefaultResolve(target, value)
// }
//
// func DefaultResolve(target reflect.Value, value any) error {
// 	source := reflect.ValueOf(value)
//
// 	targetType := unpointType(target.Type())
// 	sourceType := unpointType(source.Type())
// 	targetValue := unpointValue(target)
// 	sourceValue := unpointValue(source)
//
// 	if tryAssign(sourceType, targetType, sourceValue, targetValue) {
// 		return nil
// 	}
//
// 	if !targetValue.CanSet() {
// 		return fmt.Errorf("cannot set field value")
// 	}
//
// 	switch targetType.Kind() {
// 	case reflect.Bool:
// 		c, err := coerce.Bool(value)
// 		targetValue.SetBool(c)
// 		return err
// 	case reflect.Float32:
// 		c, err := coerce.Float32(value)
// 		targetValue.SetFloat(float64(c))
// 		return err
// 	case reflect.Float64:
// 		c, err := coerce.Float64(value)
// 		targetValue.SetFloat(c)
// 		return err
// 	case reflect.Int:
// 		c, err := coerce.Int(value)
// 		targetValue.SetInt(int64(c))
// 		return err
// 	case reflect.Int8:
// 		c, err := coerce.Int8(value)
// 		targetValue.SetInt(int64(c))
// 		return err
// 	case reflect.Int16:
// 		c, err := coerce.Int16(value)
// 		targetValue.SetInt(int64(c))
// 		return err
// 	case reflect.Int32:
// 		c, err := coerce.Int32(value)
// 		targetValue.SetInt(int64(c))
// 		return err
// 	case reflect.Int64:
// 		c, err := coerce.Int64(value)
// 		targetValue.SetInt(int64(c))
// 		return err
// 	case reflect.Uint:
// 		c, err := coerce.Uint(value)
// 		targetValue.SetUint(uint64(c))
// 		return err
// 	case reflect.Uint8:
// 		c, err := coerce.Uint8(value)
// 		targetValue.SetUint(uint64(c))
// 		return err
// 	case reflect.Uint16:
// 		c, err := coerce.Uint16(value)
// 		targetValue.SetUint(uint64(c))
// 		return err
// 	case reflect.Uint32:
// 		c, err := coerce.Uint32(value)
// 		targetValue.SetUint(uint64(c))
// 		return err
// 	case reflect.Uint64:
// 		c, err := coerce.Uint64(value)
// 		targetValue.SetUint(uint64(c))
// 		return err
// 	case reflect.String:
// 		c, err := coerce.String(value)
// 		targetValue.SetString(c)
// 		return err
// 	}
//
// 	if targetType.Kind() == reflect.Interface {
// 		// this is an interface
// 		targetValue.Set(sourceValue)
// 		return nil
// 	}
//
// 	switch sourceType.Kind() {
// 	case reflect.Slice:
// 		// this is a slice, so we expect the target to be a slice too
// 		if targetType.Kind() != reflect.Slice {
// 			return fmt.Errorf("expected an array to coerce an array into")
// 		}
// 		elemType := targetType.Elem()
// 		targetSliceValue := reflect.MakeSlice(reflect.SliceOf(elemType), 0, 0)
// 		for i := 0; i < sourceValue.Len(); i++ {
// 			// slicePath := append(path, i)
// 			sourceElemValue := sourceValue.Index(i)
// 			var targetValue reflect.Value
// 			if elemType.Kind() == reflect.Ptr {
// 				// the slice expects a pointer type
// 				targetValue = reflect.New(unpointType(elemType))
// 			} else {
// 				// the slice expects a literal type
// 				targetValue = reflect.New(elemType)
// 			}
// 			if err := DefaultResolve(targetValue, sourceElemValue.Interface()); err != nil {
// 				return err
// 			}
// 			if elemType.Kind() == reflect.Ptr {
// 				// the slice expects a pointer type
// 				targetSliceValue = reflect.Append(targetSliceValue, targetValue)
// 			} else {
// 				// the slice expects a literal type
// 				targetSliceValue = reflect.Append(targetSliceValue, unpointValue(targetValue))
// 			}
// 		}
// 		if !targetValue.CanSet() {
// 			return fmt.Errorf("cannot set '%s' ", targetType.Kind())
// 		}
//
// 		targetValue.Set(targetSliceValue)
// 	default:
// 		return fmt.Errorf("cannot coerce source of type '%s' into target of type '%s'",
// 			sourceType.Kind(), targetType.Kind())
// 	}
//
// 	return nil
// }
//
// func tryAssign(st, tt reflect.Type, sv, tv reflect.Value) bool {
// 	st = unpointType(st)
// 	tt = unpointType(tt)
// 	sv = unpointValue(sv)
// 	tv = unpointValue(tv)
//
// 	if !tv.CanSet() {
// 		return false
// 	}
//
// 	if tt.AssignableTo(st) {
// 		// the source can be directly assigned to the target
// 		tv.Set(sv)
// 		return true
// 	}
//
// 	if tt.Kind() != reflect.Slice && st.ConvertibleTo(tt) {
// 		// slices are handled manually
// 		tv.Set(sv.Convert(tt))
// 		return true
// 	}
//
// 	return false
// }
