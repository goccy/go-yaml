package yaml

import "testing"

type emptyStringValue string

func (v *emptyStringValue) UnmarshalYAML(raw []byte) error {
	*v = emptyStringValue(string(raw))
	return nil
}

type emptyStringContainer struct {
	V emptyStringValue `yaml:"value"`
}

func TestEmptyString(t *testing.T) {
	var container emptyStringContainer
	if err := Unmarshal([]byte(`value: ""`), &container); err != nil {
		t.Fatalf("failed to unmarshal %v", err)
	}
	if container.V != "" {
		t.Fatalf("expected empty string, but %q is set", container.V)
	}
}
