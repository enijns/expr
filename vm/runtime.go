package vm

//go:generate go run ./generate

import (
	"fmt"
	"math"
	"reflect"
)

type Call struct {
	Name string
	Size int
}

type Scope map[string]interface{}

func fetch(from interface{}, i interface{}) interface{} {
	v := reflect.ValueOf(from)
	switch v.Kind() {

	case reflect.Array, reflect.Slice, reflect.String:
		index := toInt(i)
		value := v.Index(int(index))
		if value.IsValid() && value.CanInterface() {
			return value.Interface()
		}

	case reflect.Map:
		value := v.MapIndex(reflect.ValueOf(i))
		if value.IsValid() && value.CanInterface() {
			return value.Interface()
		}

	case reflect.Struct:
		value := v.FieldByName(reflect.ValueOf(i).String())
		if value.IsValid() && value.CanInterface() {
			return value.Interface()
		}

	case reflect.Ptr:
		value := v.Elem()
		if value.IsValid() && value.CanInterface() {
			return fetch(value.Interface(), i)
		}

	}
	panic(fmt.Sprintf("%v doesn't contains %v", from, i))
}

func slice(array, from, to interface{}) interface{} {
	v := reflect.ValueOf(array)
	switch v.Kind() {

	case reflect.Array, reflect.Slice, reflect.String:
		value := v.Slice(toInt(from), toInt(to))
		if value.IsValid() && value.CanInterface() {
			return value.Interface()
		}

	case reflect.Ptr:
		value := v.Elem()
		if value.IsValid() && value.CanInterface() {
			return slice(value.Interface(), from, to)
		}

	}
	panic(fmt.Sprintf("cannot slice %v", from))
}

func fetchFn(from interface{}, name string) reflect.Value {
	v := reflect.ValueOf(from)

	// Methods can be defined on any type.
	if v.NumMethod() > 0 {
		method := v.MethodByName(name)
		if method.IsValid() {
			return method
		}
	}

	d := v
	if v.Kind() == reflect.Ptr {
		d = v.Elem()
	}

	switch d.Kind() {
	case reflect.Map:
		value := d.MapIndex(reflect.ValueOf(name))
		if value.IsValid() && value.CanInterface() {
			return value.Elem()
		}
	case reflect.Struct:
		// If struct has not method, maybe it has func field.
		// To access this field we need dereference value.
		value := d.FieldByName(name)
		if value.IsValid() {
			return value
		}
	}
	panic(fmt.Sprintf(`cannot get "%v" from %T`, name, from))
}

func in(needle interface{}, array interface{}) bool {
	if array == nil {
		return false
	}
	v := reflect.ValueOf(array)

	switch v.Kind() {

	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			value := v.Index(i)
			if value.IsValid() && value.CanInterface() {
				if equal(value.Interface(), needle).(bool) {
					return true
				}
			}
		}
		return false

	case reflect.Map:
		n := reflect.ValueOf(needle)
		if !n.IsValid() {
			panic(fmt.Sprintf("cannot use %T as index to %T", needle, array))
		}
		value := v.MapIndex(n)
		if value.IsValid() {
			return true
		}
		return false

	case reflect.Struct:
		n := reflect.ValueOf(needle)
		if !n.IsValid() || n.Kind() != reflect.String {
			panic(fmt.Sprintf("cannot use %T as field name of %T", needle, array))
		}
		value := v.FieldByName(n.String())
		if value.IsValid() {
			return true
		}
		return false

	case reflect.Ptr:
		value := v.Elem()
		if value.IsValid() && value.CanInterface() {
			return in(needle, value.Interface())
		}
		return false
	}

	panic(fmt.Sprintf(`operator "in"" not defined on %T`, array))
}

func length(a interface{}) int {
	v := reflect.ValueOf(a)
	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		return v.Len()
	default:
		panic(fmt.Sprintf("invalid argument for len (type %T)", a))
	}
}

func negate(i interface{}) interface{} {
	switch v := i.(type) {
	case float32:
		return -v
	case float64:
		return -v

	case int:
		return -v
	case int8:
		return -v
	case int16:
		return -v
	case int32:
		return -v
	case int64:
		return -v

	case uint:
		return -v
	case uint8:
		return -v
	case uint16:
		return -v
	case uint32:
		return -v
	case uint64:
		return -v

	default:
		panic(fmt.Sprintf("invalid operation: - %T", v))
	}
}

func exponent(a, b interface{}) float64 {
	return math.Pow(toFloat64(a), toFloat64(b))
}

func makeRange(a, b interface{}) []int {
	min := toInt(a)
	max := toInt(b)
	size := max - min + 1
	rng := make([]int, size)
	for i := range rng {
		rng[i] = min + i
	}
	return rng
}

func toInt(a interface{}) int {
	switch x := a.(type) {
	case float32:
		return int(x)
	case float64:
		return int(x)

	case int:
		return int(x)
	case int8:
		return int(x)
	case int16:
		return int(x)
	case int32:
		return int(x)
	case int64:
		return int(x)

	case uint:
		return int(x)
	case uint8:
		return int(x)
	case uint16:
		return int(x)
	case uint32:
		return int(x)
	case uint64:
		return int(x)

	default:
		panic(fmt.Sprintf("invalid operation: int(%T)", x))
	}
}

func toInt64(a interface{}) int64 {
	switch x := a.(type) {
	case float32:
		return int64(x)
	case float64:
		return int64(x)

	case int:
		return int64(x)
	case int8:
		return int64(x)
	case int16:
		return int64(x)
	case int32:
		return int64(x)
	case int64:
		return int64(x)

	case uint:
		return int64(x)
	case uint8:
		return int64(x)
	case uint16:
		return int64(x)
	case uint32:
		return int64(x)
	case uint64:
		return int64(x)

	default:
		panic(fmt.Sprintf("invalid operation: int64(%T)", x))
	}
}

func toFloat64(a interface{}) float64 {
	switch x := a.(type) {
	case float32:
		return float64(x)
	case float64:
		return x

	case int:
		return float64(x)
	case int8:
		return float64(x)
	case int16:
		return float64(x)
	case int32:
		return float64(x)
	case int64:
		return float64(x)

	case uint:
		return float64(x)
	case uint8:
		return float64(x)
	case uint16:
		return float64(x)
	case uint32:
		return float64(x)
	case uint64:
		return float64(x)

	default:
		panic(fmt.Sprintf("invalid operation: float64(%T)", x))
	}
}
