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
	data interface{}
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

func getPart(obj interface{}, part interface{}, createMissingObject bool) (interface{}, error) {

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
	var m interface{}
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
func (n *Map) Data() interface{} {
	return n.data
}
