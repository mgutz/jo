package nestedjson

import "fmt"

// Bool returns bool value from path.
func (n *Map) Bool(path string) (bool, error) {
	o, err := n.Get(path)
	if err != nil {
		return false, err
	}
	switch rv := o.(type) {
	case bool:
		return rv, nil
	default:
		return false, fmt.Errorf("%s is not a bool: %q", path, o)
	}
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
	switch rv := o.(type) {
	case int:
		return float64(rv), nil
	case float64:
		return rv, nil
	default:
		return 0, fmt.Errorf("%s is not a float: %q", path, o)
	}
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
	switch rv := o.(type) {
	case int:
		return rv, nil
	case float64:
		return int(rv), nil
	default:
		return 0, fmt.Errorf("%s is not an integer: %q", path, o)
	}
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

// String returns a string value from path.
func (n *Map) String(path string) (string, error) {
	o, err := n.Get(path)
	if err != nil {
		return "", err
	}
	switch rv := o.(type) {
	case string:
		return rv, nil
	default:
		return "", fmt.Errorf("%s is not a string: %q", path, o)
	}
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

// Slice returns slice of interface{} from path.
func (n *Map) Slice(path string) ([]interface{}, error) {
	o, err := n.Get(path)
	if err != nil {
		return nil, err
	}
	switch rv := o.(type) {
	case []interface{}:
		return rv, nil
	default:
		return nil, fmt.Errorf("%s is not n Slice: %q", path, o)
	}
}

// MaySlice should get value from path or return val.
func (n *Map) MaySlice(path string, val []interface{}) []interface{} {
	v, err := n.Slice(path)
	if err != nil {
		return val
	}
	return v
}

// MustSlice gets string value from path or panics.
func (n *Map) MustSlice(path string) []interface{} {
	v, err := n.Slice(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "[]interface{}"))
	}
	return v
}
