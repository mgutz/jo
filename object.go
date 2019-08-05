package jo

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"
)

var partRe = regexp.MustCompile(`^([\w-]*)((\[[0-9]+\])+)*$`)
var arrayIndexRe = regexp.MustCompile(`\[([0-9]+)\]`)

// ErrKeyDoesNotExist is returned if any part of a path cannot be traversed
var ErrKeyDoesNotExist = errors.New("Key does not exist")

const mustFormat = `Path not found or could not convert to %s: "%s"`

// Object represents a JSON object. Remember a JSON object can be
// a literal value like a quoted string.
type Object struct {
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
		if arr, ok := obj.([]map[string]interface{}); ok {
			if p < len(arr) {
				return arr[p], nil
			}
			return nil, fmt.Errorf("Array index out of bounds: %d", p)
		} else if arr, ok := obj.([]interface{}); ok {
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
			return nil, ErrKeyDoesNotExist
		}
		return nil, fmt.Errorf("%s is not an object: %T", obj, obj)
	}
	return nil, fmt.Errorf("Invalid Part: %T", part)
}

// New creates a new JSON object struct.
func New() *Object {
	return &Object{make(map[string]interface{})}
}

// NewFromAny creates a new JSON struct from any marhsallable JSON
// object.
func NewFromAny(v interface{}) (*Object, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return NewFromBytes(b)
}

// NewFromJSONFile creates a new Object from a filename. filename can be prefixed
// with ~ for home directory.
func NewFromJSONFile(filename string) (*Object, error) {
	fname, err := Untildify(filename)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	return NewFromReadCloser(file)
}

// NewFromMap creates a new JSON struct from an existing map
func NewFromMap(m map[string]interface{}) *Object {
	return &Object{m}
}

// NewFromReadCloser decodes JSON from an io.ReadCloser
// suchs as Request.Body and Response.Body and closes it.
func NewFromReadCloser(body io.ReadCloser) (*Object, error) {
	decoder := json.NewDecoder(body)
	var result map[string]interface{}
	err := decoder.Decode(&result)
	if err != nil {
		body.Close()
		return nil, err
	}
	err = body.Close()
	return &Object{result}, err
}

// NewFromBytes creates an object directly from JSON encoded bytes.
func NewFromBytes(b []byte) (*Object, error) {
	obj := New()
	err := obj.UnmarshalJSON(b)
	return obj, err
}

// NewFromString creates an object directly from a JSON encoded string.
func NewFromString(json string) (*Object, error) {
	obj := New()
	err := obj.UnmarshalJSON([]byte(json))
	return obj, err
}

// MarshalIndent pretty encodes JSON to indented []byte.
func (n *Object) MarshalIndent(prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(n.data, prefix, indent)
}

// Prettify returns the JSON indented representation of this map.
func (n *Object) Prettify() string {
	b, err := n.MarshalIndent("", "  ")
	if err != nil {
		log.Println("Unable to prettify:", err)
		return "nil"
	}
	return string(b)
}

// Stringify returns the JSON string representation of this map.
func (n *Object) Stringify() string {
	b, err := n.MarshalJSON()
	if err != nil {
		log.Println("Unable to stringify:", err)
		return "nil"
	}
	return string(b)
}

// Delete deletes a path. If the path does not exist then no error is returned.
func (n *Object) Delete(path string) error {
	parts, err := splitPath(path)
	if err != nil {
		return err
	}

	L := len(parts)

	var parent interface{} = n.data
	for i, part := range parts {

		// last item
		if i == L-1 {
			m, ok := parent.(map[string]interface{})
			if ok {
				k, ok := part.(string)
				if !ok {
					return errors.New("Map key must be a string")
				}
				delete(m, k)
				continue
			}

			a, ok := parent.([]interface{})
			if ok {
				idx, ok := part.(int)
				if !ok {
					return errors.New("Index must be an int")
				}
				// deletes a slice entry safely
				// SEE https://github.com/golang/go/wiki/SliceTricks
				copy(a[idx:], a[idx+1:])
				a[len(a)-1] = nil // or the zero value of T
				// TODO this doesn't actually change teh slice, need to track
				// grandparents
				a = a[:len(a)-1]
				continue
			}

			return errors.New("Parent node is neither map[string]interface{} or []interface{}")
		}

		parent, err = getPart(parent, part, false)
		if err != nil {
			return err
		}
	}
	return nil
}

// Get gets value at path which may contain "." for path traversal.
func (n *Object) Get(path string) (interface{}, error) {
	if path == "." {
		return n.data, nil
	}

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
func (n *Object) Set(path string, val interface{}) error {
	parts, err := splitPath(path)
	if err != nil {
		return err
	}

	curr := n.data
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
			return fmt.Errorf("Not an object: %#v", parts)
		}
	}

	return nil
}

// UnmarshalJSON implements Unmarshaller interface.
func (n *Object) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &n.data)
}

// MarshalJSON implements Marshaller interface.
func (n *Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.data)
}

// Data returns the entire data map.
func (n *Object) Data() interface{} {
	return n.data
}

// Untildify replaces leading ~ with current user's home directory
func Untildify(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(path, "~") {
		return currentUser.HomeDir + path[1:], nil
	}
	return path, nil
}

// Scan implements datbase scanner
func (n *Object) Scan(value interface{}) error {
	if value == nil {
		n.data = "null"
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &n.data)
	case string:
		return json.Unmarshal([]byte(v), &n.data)
	default:
		b, err := json.Marshal(value)
		if err != nil {
			return err
		}

		return json.Unmarshal(b, &n.data)
	}
}

// Value implements database valuer
func (n Object) Value() (driver.Value, error) {
	return n.data, nil
}
