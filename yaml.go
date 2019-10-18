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

func Unmarshal(data []byte, v interface{}) error {
	dec := NewDecoder(bytes.NewBuffer(data))
	if err := dec.Decode(v); err != nil {
		return xerrors.Errorf("failed to unmarshal: %w", err)
	}
	return nil
}
