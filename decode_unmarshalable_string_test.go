package yaml

import (
	"encoding/json"
	"testing"
)

type unmarshalableStringValue string

func (v *unmarshalableStringValue) UnmarshalYAML(raw []byte) error {
	*v = unmarshalableStringValue(string(raw))
	return nil
}

type unmarshalableStringContainer struct {
	V unmarshalableStringValue `yaml:"value" json:"value"`
}

func TestUnmarshalableString(t *testing.T) {
	var container unmarshalableStringContainer
	t.Run("empty string", func(t *testing.T) {
		if err := Unmarshal([]byte(`value: ""`), &container); err != nil {
			t.Fatalf("failed to unmarshal %v", err)
		}
		if container.V != "" {
			t.Fatalf("expected empty string, but %q is set", container.V)
		}
	})
	t.Run("filled string", func(t *testing.T) {
		if err := Unmarshal([]byte(`value: "aaa"`), &container); err != nil {
			t.Fatalf("failed to unmarshal %v", err)
		}
		if container.V != "aaa" {
			t.Fatalf("expected \"aaa\", but %q is set", container.V)
		}
	})
	t.Run("single-quoted string", func(t *testing.T) {
		if err := Unmarshal([]byte(`value: 'aaa'`), &container); err != nil {
			t.Fatalf("failed to unmarshal %v", err)
		}
		if container.V != "aaa" {
			t.Fatalf("expected \"aaa\", but %q is set", container.V)
		}
	})
	t.Run("(json) empty string", func(t *testing.T) {
		if err := json.Unmarshal([]byte(`{"value": ""}`), &container); err != nil {
			t.Fatalf("failed to unmarshal %v", err)
		}
		if container.V != "" {
			t.Fatalf("expected empty string, but %q is set", container.V)
		}
	})
	t.Run("(json) filled string", func(t *testing.T) {
		if err := json.Unmarshal([]byte(`{"value": "aaa"}`), &container); err != nil {
			t.Fatalf("failed to unmarshal %v", err)
		}
		if container.V != "aaa" {
			t.Fatalf("expected \"aaa\", but %q is set", container.V)
		}
	})
}
