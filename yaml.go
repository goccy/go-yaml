package yaml

import (
	"bytes"

	"golang.org/x/xerrors"
)

type Marshaler interface {
	MarshalYAML() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalYAML([]byte) error
}

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, xerrors.Errorf("failed to marshal: %w", err)
	}
	return buf.Bytes(), nil
}

func Unmarshal(data []byte, v interface{}) error {
	dec := NewDecoder(bytes.NewBuffer(data))
	if err := dec.Decode(v); err != nil {
		return xerrors.Errorf("failed to unmarshal: %w", err)
	}
	return nil
}
