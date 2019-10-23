package yaml

import (
	"bytes"

	"github.com/goccy/go-yaml/errors"
)

// Marshaler interface may be implemented by types to customize their
// behavior when being marshaled into a YAML document. The returned value
// is marshaled in place of the original value implementing Marshaler.
//
// If an error is returned by MarshalYAML, the marshaling procedure stops
// and returns with the provided error.
type Marshaler interface {
	MarshalYAML() ([]byte, error)
}

// Unmarshaler interface may be implemented by types to customize their
// behavior when being unmarshaled from a YAML document.
type Unmarshaler interface {
	UnmarshalYAML([]byte) error
}

// Marshal serializes the value provided into a YAML document. The structure
// of the generated document will reflect the structure of the value itself.
// Maps and pointers (to struct, string, int, etc) are accepted as the in value.
//
// Struct fields are only marshalled if they are exported (have an upper case
// first letter), and are marshalled using the field name lowercased as the
// default key. Custom keys may be defined via the "yaml" name in the field
// tag: the content preceding the first comma is used as the key, and the
// following comma-separated options are used to tweak the marshalling process.
// Conflicting names result in a runtime error.
//
// The field tag format accepted is:
//
//     `(...) yaml:"[<key>][,<flag1>[,<flag2>]]" (...)`
//
// The following flags are currently supported:
//
//     omitempty    Only include the field if it's not set to the zero
//                  value for the type or to empty slices or maps.
//                  Zero valued structs will be omitted if all their public
//                  fields are zero, unless they implement an IsZero
//                  method (see the IsZeroer interface type), in which
//                  case the field will be included if that method returns true.
//
//     flow         Marshal using a flow style (useful for structs,
//                  sequences and maps).
//
//     inline       Inline the field, which must be a struct or a map,
//                  causing all of its fields or keys to be processed as if
//                  they were part of the outer struct. For maps, keys must
//                  not conflict with the yaml keys of other struct fields.
//
//     anchor       Marshal with anchor. If want to define anchor name explicitly, use anchor=name style.
//                  Otherwise, if used 'anchor' name only, used the field name lowercased as the anchor name
//
//     alias        Marshal with alias. If want to define alias name explicitly, use alias=name style.
//                  Otherwise, If omitted alias name and the field type is pointer type,
//                  assigned anchor name automatically from same pointer address.
//
// In addition, if the key is "-", the field is ignored.
//
// For example:
//
//     type T struct {
//         F int `yaml:"a,omitempty"`
//         B int
//     }
//     yaml.Marshal(&T{B: 2}) // Returns "b: 2\n"
//     yaml.Marshal(&T{F: 1}} // Returns "a: 1\nb: 0\n"
//
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, errors.Wrapf(err, "failed to marshal", err)
	}
	return buf.Bytes(), nil
}

// Unmarshal decodes the first document found within the in byte slice
// and assigns decoded values into the out value.
//
// Struct fields are only unmarshalled if they are exported (have an
// upper case first letter), and are unmarshalled using the field name
// lowercased as the default key. Custom keys may be defined via the
// "yaml" name in the field tag: the content preceding the first comma
// is used as the key, and the following comma-separated options are
// used to tweak the marshalling process (see Marshal).
// Conflicting names result in a runtime error.
//
// For example:
//
//     type T struct {
//         F int `yaml:"a,omitempty"`
//         B int
//     }
//     var t T
//     yaml.Unmarshal([]byte("a: 1\nb: 2"), &t)
//
// See the documentation of Marshal for the format of tags and a list of
// supported tag options.
//
func Unmarshal(data []byte, v interface{}) error {
	dec := NewDecoder(bytes.NewBuffer(data))
	if err := dec.Decode(v); err != nil {
		return errors.Wrapf(err, "failed to unmarshal")
	}
	return nil
}
