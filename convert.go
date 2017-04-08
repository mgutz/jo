package jo

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/mgutz/to"
)

var mapT = reflect.TypeOf(map[string]interface{}{})
var arrayT = reflect.TypeOf([]interface{}{})

func toMap(any interface{}) (map[string]interface{}, error) {
	if m, ok := any.(map[string]interface{}); ok {
		return m, nil
	}

	value := reflect.ValueOf(any)
	if !value.IsValid() {
		return nil, errors.New("Zero value is not convertible to map[string]interface{}")
	}

	if T := value.Type(); T.AssignableTo(mapT) {
		return value.Convert(mapT).Interface().(map[string]interface{}), nil
	}

	return nil, fmt.Errorf("Type %T is not convertible to map[string]interface{}", any)
}

func toArray(any interface{}) ([]interface{}, error) {
	if any == nil {
		return nil, nil
	}

	if arr, ok := any.([]interface{}); ok {
		return arr, nil
	}

	value := reflect.ValueOf(any)
	if !value.IsValid() {
		return nil, errors.New("Zero value is not convertible to []interface{}")
	}
	if T := value.Type(); T.Kind() == reflect.Slice {
		L := value.Len()
		arr := make([]interface{}, L)
		for i := 0; i < L; i++ {
			arr[i] = value.Index(i).Interface()
		}
		return arr, nil
	}

	return nil, fmt.Errorf("Type %T is not convertible to []interface{}", any)
}

// Bool returns bool value from path.
func (n *Object) Bool(path string) (bool, error) {
	o, err := n.Get(path)
	if err != nil {
		return false, err
	}
	return to.Bool(o)
}

// AsBool returns bool value from path else false.
func (n *Object) AsBool(path string) bool {
	v, err := n.Bool(path)
	if err != nil {
		return false
	}
	return v
}

// OrBool should get value from path or return val.
func (n *Object) OrBool(path string, val bool) bool {
	v, err := n.Bool(path)
	if err != nil {
		return val
	}
	return v
}

// MustBool gets string value from path or panics.
func (n *Object) MustBool(path string) bool {
	v, err := n.Bool(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "bool", path))
	}
	return v
}

// Float returns a float64 value from path.
func (n *Object) Float(path string) (float64, error) {
	o, err := n.Get(path)
	if err != nil {
		return 0, err
	}
	return to.Float64(o)
}

// AsFloat returns float64 value from path else 0.
func (n *Object) AsFloat(path string) float64 {
	v, err := n.Float(path)
	if err != nil {
		return 0
	}
	return v
}

// OrFloat should get value from path or return val.
func (n *Object) OrFloat(path string, val float64) float64 {
	v, err := n.Float(path)
	if err != nil {
		return val
	}
	return v
}

// MustFloat gets string value from path or panics.
func (n *Object) MustFloat(path string) float64 {
	v, err := n.Float(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "float64", path))
	}
	return v
}

// Int returns an integer value from path.
func (n *Object) Int(path string) (int, error) {
	o, err := n.Get(path)
	if err != nil {
		return 0, err
	}

	n64, err := to.Int64(o)
	if err != nil {
		return 0, err
	}
	return int(n64), nil
}

// AsInt returns int value from path else 0.
func (n *Object) AsInt(path string) int {
	v, err := n.Int(path)
	if err != nil {
		return 0
	}
	return v
}

// OrInt should get value from path or return val.
func (n *Object) OrInt(path string, val int) int {
	v, err := n.Int(path)
	if err != nil {
		return val
	}
	return v
}

// MustInt gets string value from path or panics.
func (n *Object) MustInt(path string) int {
	v, err := n.Int(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "int", path))
	}
	return v
}

// Int64 returns an integer value from path.
func (n *Object) Int64(path string) (int64, error) {
	o, err := n.Get(path)
	if err != nil {
		return 0, err
	}
	return to.Int64(o)
}

// AsInt64 returns int64 value from path else 0.
func (n *Object) AsInt64(path string) int64 {
	v, err := n.Int64(path)
	if err != nil {
		return 0
	}
	return v
}

// OrInt64 should get value from path or return val.
func (n *Object) OrInt64(path string, val int64) int64 {
	v, err := n.Int64(path)
	if err != nil {
		return val
	}
	return v
}

// MustInt64 gets string value from path or panics.
func (n *Object) MustInt64(path string) int64 {
	v, err := n.Int64(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "int", path))
	}
	return v
}

// String returns a string value from path.
func (n *Object) String(path string) (string, error) {
	o, err := n.Get(path)
	if err != nil {
		return "", err
	}
	return to.String(o), nil
}

// AsString returns string value from path else "".
func (n *Object) AsString(path string) string {
	v, err := n.String(path)
	if err != nil {
		return ""
	}
	return v
}

// OrString should get value from path or return s.
func (n *Object) OrString(path string, s string) string {
	v, err := n.String(path)
	if err != nil {
		return s
	}
	return v
}

// MustString gets string value from path or panics.
func (n *Object) MustString(path string) string {
	v, err := n.String(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "string", path))
	}
	return v
}

////////// Collections

// Map returns the map at path.
func (n *Object) Map(path string) (map[string]interface{}, error) {
	o, err := n.Get(path)
	if err != nil {
		return nil, err
	}
	return toMap(o)
}

// AsMap returns map  from path else nil.
func (n *Object) AsMap(path string) map[string]interface{} {
	v, err := n.Map(path)
	if err != nil {
		return nil
	}
	return v
}

// OrMap should get value from path or return val.
func (n *Object) OrMap(path string, val map[string]interface{}) map[string]interface{} {
	v, err := n.Map(path)
	if err != nil {
		return val
	}
	return v
}

// MustMap gets string value from path or panics.
func (n *Object) MustMap(path string) map[string]interface{} {
	v, err := n.Map(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "map[string]interface{}", path))
	}
	return v
}

// Array returns slice of interface{} from path.
func (n *Object) Array(path string) ([]interface{}, error) {
	o, err := n.Get(path)
	if err != nil {
		return nil, err
	}
	return toArray(o)
}

// StringArray returns a string slice from path. Null values are converted
// to "".
func (n *Object) StringArray(path string) ([]string, error) {
	arr, err := n.Array(path)
	if err != nil {
		return nil, err
	}
	result := make([]string, len(arr))
	for i, val := range arr {
		if val == nil {
			result[i] = ""
		} else if s, ok := val.(string); ok {
			result[i] = s
		} else {
			return nil, fmt.Errorf("Array has non-string values")
		}
	}
	return result, nil
}

// AsObjects returns slice of Objects from path.
func (n *Object) AsObjects(path string) []*Object {
	o, err := n.Array(path)
	if err != nil {
		return nil
	}
	objects := []*Object{}
	for _, iface := range o {
		objects = append(objects, &Object{iface})
	}
	return objects
}

// // OrArray should get value from path or return val.
// func (n *Object) OrArray(path string, val []interface{}) []Object {
// 	v, err := n.Array(path)
// 	if err != nil {
// 		return val
// 	}
// 	return v
// }

// MustArray gets array value from path or panics.
func (n *Object) MustArray(path string) []interface{} {
	v, err := n.Array(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "[]interface{}", path))
	}
	return v
}

// At returns a sub map at path.
func (n *Object) At(path string) (*Object, error) {
	o, err := n.Get(path)
	if err != nil {
		return nil, err
	}
	return &Object{o}, nil
}

// MustAt gets a sub map at path or panics.
func (n *Object) MustAt(path string) *Object {
	m, err := n.At(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "*Object", path))
	}
	return m
}
