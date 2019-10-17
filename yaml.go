package yaml

type Marshaler interface {
	MarshalYAML() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalYAML([]byte) error
}

func Unmarshal(data []byte, v interface{}) error {
	return nil
}
