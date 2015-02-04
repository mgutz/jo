package nestedjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//var partRe = regexp.MustCompile(`^([A-Za-z0-9_]*)((\[[0-9]+\])+)*$`)
var partRe = regexp.MustCompile(`^([\w-]*)((\[[0-9]+\])+)*$`)
var arrayIndexRe = regexp.MustCompile(`\[([0-9]+)\]`)

const mustFormat = "Path not found or could not convert to %s"

// Map is a generic map.
type Map struct {
	data map[string]interface{}
}

func splitPath(path string) ([]interface{}, error) {
	var rv []interface{}
	parts := strings.Split(path, ".")
	for i, part := range parts {
		if part == "" {
			return nil, errors.New("Invalid path: " + path)
		}
		partMatches := partRe.FindStringSubmatch(part)
		if len(partMatches) == 0 {
			return nil, errors.New("Invalid part: " + part)
		}
		// abc[0][1][2]
		objKey := partMatches[1]       //abc
		arrayIndexes := partMatches[2] // [0][1][2]

		if objKey == "" {
			if i > 0 {
				return nil, errors.New("Invalid path: " + path)
			}
		} else {
			rv = append(rv, objKey)
		}

		if arrayIndexes != "" {
			arrayIndexMatches := arrayIndexRe.FindAllStringSubmatch(arrayIndexes, -1)
			for _, indexMatch := range arrayIndexMatches {
				intIndex, _ := strconv.Atoi(indexMatch[1])
				rv = append(rv, intIndex)
			}
		}
	}
	return rv, nil
}

func getPart(obj interface{}, part interface{},
	createMissingObject bool) (interface{}, error) {

	switch p := part.(type) {
	case int:
		if arr, ok := obj.([]interface{}); ok {
			if p < len(arr) {
				return arr[p], nil
			}
			return nil, fmt.Errorf("Array index out of bounds: %d", p)
		}
		return nil, fmt.Errorf("%s is not an array: %T", obj, obj)

	case string:
		if m, ok := obj.(map[string]interface{}); ok {
			if rv, ok := m[p]; ok {
				return rv, nil
			}
			if createMissingObject {
				rv := make(map[string]interface{})
				m[p] = rv
				return rv, nil
			}
			return nil, fmt.Errorf("Key does not exist: %s", p)
		}
		return nil, fmt.Errorf("%s is not an object: %T", obj, obj)
	}
	return nil, fmt.Errorf("Invalid Part: %T", part)
}

// New creates a new JSON struct.
func New() *Map {
	return &Map{make(map[string]interface{})}
}

// NewFromMap creates a new JSON struct from an existing map
func NewFromMap(m map[string]interface{}) *Map {
	return &Map{m}
}

// Unmarshal decodes bytes to JSON.
func Unmarshal(b []byte) (*Map, error) {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return &Map{m}, nil
}

// Marshal encodes JSON to []byte.
func (n *Map) Marshal() ([]byte, error) {
	return json.Marshal(n.data)
}

// MarshalIndent pretty encodes JSON to indented []byte.
func (n *Map) MarshalIndent(prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(n.data, prefix, indent)
}

// Get gets value at path which may contain "." for path traversal.
func (n *Map) Get(path string) (interface{}, error) {
	parts, err := splitPath(path)
	if err != nil {
		return nil, err
	}

	var curr interface{} = n.data
	for _, part := range parts {
		curr, err = getPart(curr, part, false)
		if err != nil {
			return nil, err
		}
	}
	return curr, nil
}

// Set sets the value at path.
func (n *Map) Set(path string, val interface{}) error {
	parts, err := splitPath(path)
	if err != nil {
		return err
	}

	var curr interface{} = n.data
	for _, part := range parts[:len(parts)-1] {
		curr, err = getPart(curr, part, true)
		if err != nil {
			return err
		}
	}

	switch k := parts[len(parts)-1].(type) {
	case int:
		if arr, ok := curr.([]interface{}); ok {
			arr[k] = val
		} else {
			return fmt.Errorf("Not an array: %s", curr)
		}

	case string:
		if m, ok := curr.(map[string]interface{}); ok {
			m[k] = val
		} else {
			return fmt.Errorf("Not an object: %s", curr)
		}
	}

	return nil

}

// Data returns the entire data map.
func (n *Map) Data() map[string]interface{} {
	return n.data
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
		return "", fmt.Errorf("%s is not a string: %v", path, o)
	}
}

// SafeString should get value from path or return s.
func (n *Map) SafeString(path string, s string) string {
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
		return 0, fmt.Errorf("%s is not an integer: %v", path, o)
	}
}

// SafeInt should get value from path or return val.
func (n *Map) SafeInt(path string, val int) int {
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
		return 0, fmt.Errorf("%s is not a float: %v", path, o)
	}
}

// SafeFloat should get value from path or return val.
func (n *Map) SafeFloat(path string, val float64) float64 {
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
		return false, fmt.Errorf("%s is not a bool: %v", path, o)
	}
}

// SafeBool should get value from path or return val.
func (n *Map) SafeBool(path string, val bool) bool {
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
		return nil, fmt.Errorf("%s is not an Array: %v", path, o)
	}
}

// SafeArray should get value from path or return val.
func (n *Map) SafeArray(path string, val []interface{}) []interface{} {
	v, err := n.Array(path)
	if err != nil {
		return val
	}
	return v
}

// MustArray gets string value from path or panics.
func (n *Map) MustArray(path string) []interface{} {
	v, err := n.Array(path)
	if err != nil {
		panic(fmt.Sprintf(mustFormat, "[]interface{}"))
	}
	return v
}

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
		return nil, fmt.Errorf("%s is not a map: %v", path, o)
	}
}

// SafeMap should get value from path or return val.
func (n *Map) SafeMap(path string, val map[string]interface{}) map[string]interface{} {
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
