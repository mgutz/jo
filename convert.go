package nestedjson

import (
	"fmt"

	"github.com/mgutz/to"
)

// Bool returns bool value from path.
func (n *Map) Bool(path string) (bool, error) {
	o, err := n.Get(path)
	if err != nil {
		return false, err
	}
	return to.Bool(o)
}

// MayBool should get value from path or return val.
func (n *Map) MayBool(path string, val bool) bool {
	v, err := n.Bool(path)
	if err != nil {
		return val
	}
	return v
}

// MustBool gets string value from path or panics.
func (n *Map) MustBool(path string) bool {
	v, err := n.Bool(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "bool"))
	}
	return v
}

// Float returns a float64 value from path.
func (n *Map) Float(path string) (float64, error) {
	o, err := n.Get(path)
	if err != nil {
		return 0, err
	}
	return to.Float64(o)
}

// MayFloat should get value from path or return val.
func (n *Map) MayFloat(path string, val float64) float64 {
	v, err := n.Float(path)
	if err != nil {
		return val
	}
	return v
}

// MustFloat gets string value from path or panics.
func (n *Map) MustFloat(path string) float64 {
	v, err := n.Float(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "float64"))
	}
	return v
}

// Int returns an integer value from path.
func (n *Map) Int(path string) (int, error) {
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

// MayInt should get value from path or return val.
func (n *Map) MayInt(path string, val int) int {
	v, err := n.Int(path)
	if err != nil {
		return val
	}
	return v
}

// MustInt gets string value from path or panics.
func (n *Map) MustInt(path string) int {
	v, err := n.Int(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "int"))
	}
	return v
}

// Int64 returns an integer value from path.
func (n *Map) Int64(path string) (int64, error) {
	o, err := n.Get(path)
	if err != nil {
		return 0, err
	}
	return to.Int64(o)
}

// MayInt64 should get value from path or return val.
func (n *Map) MayInt64(path string, val int64) int64 {
	v, err := n.Int64(path)
	if err != nil {
		return val
	}
	return v
}

// MustInt64 gets string value from path or panics.
func (n *Map) MustInt64(path string) int64 {
	v, err := n.Int64(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "int"))
	}
	return v
}

// String returns a string value from path.
func (n *Map) String(path string) (string, error) {
	o, err := n.Get(path)
	if err != nil {
		return "", err
	}
	return to.String(o), nil
}

// MayString should get value from path or return s.
func (n *Map) MayString(path string, s string) string {
	v, err := n.String(path)
	if err != nil {
		return s
	}
	return v
}

// MustString gets string value from path or panics.
func (n *Map) MustString(path string) string {
	v, err := n.String(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "string"))
	}
	return v
}

////////// Collections

// Map returns the map at path.
func (n *Map) Map(path string) (map[string]interface{}, error) {
	o, err := n.Get(path)
	if err != nil {
		return nil, err
	}
	switch rv := o.(type) {
	case map[string]interface{}:
		return rv, nil
	default:
		return nil, fmt.Errorf("%s is not a map: %q", path, o)
	}
}

// MayMap should get value from path or return val.
func (n *Map) MayMap(path string, val map[string]interface{}) map[string]interface{} {
	v, err := n.Map(path)
	if err != nil {
		return val
	}
	return v
}

// MustMap gets string value from path or panics.
func (n *Map) MustMap(path string) map[string]interface{} {
	v, err := n.Map(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "map[string]interface{}"))
	}
	return v
}

// Array returns slice of interface{} from path.
func (n *Map) Array(path string) ([]interface{}, error) {
	o, err := n.Get(path)
	if err != nil {
		return nil, err
	}
	switch rv := o.(type) {
	case []interface{}:
		return rv, nil
	default:
		return nil, fmt.Errorf("%s is not n Array: %q", path, o)
	}
}

// MayArray should get value from path or return val.
func (n *Map) MayArray(path string, val []interface{}) []interface{} {
	v, err := n.Array(path)
	if err != nil {
		return val
	}
	return v
}

// MustArray gets array value from path or panics.
func (n *Map) MustArray(path string) []interface{} {
	v, err := n.Array(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "[]interface{}"))
	}
	return v
}

// At returns a sub map at path.
func (n *Map) At(path string) (*Map, error) {
	o, err := n.Get(path)
	if err != nil {
		return nil, err
	}
	return &Map{o}, nil
}

// MustAt gets a sub map at path or panics.
func (n *Map) MustAt(path string) *Map {
	m, err := n.At(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "*Map"))
	}
	return m
}
