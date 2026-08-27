package yaml_test

import (
	"testing"

	"github.com/goccy/go-yaml"
)

// Regression test for #916: a float64 in scientific notation must survive a
// dynamic (map[string]any) Unmarshal -> Marshal -> Unmarshal roundtrip as a
// float64, not decode back to a string.
func TestRoundtripScientificNotationIssue916(t *testing.T) {
	var before map[string]any
	if err := yaml.Unmarshal([]byte("value: 100000000000000000000.\n"), &before); err != nil {
		t.Fatal(err)
	}
	if _, ok := before["value"].(float64); !ok {
		t.Fatalf("before: value = %T, want float64", before["value"])
	}

	encoded, err := yaml.Marshal(before)
	if err != nil {
		t.Fatal(err)
	}

	var after map[string]any
	if err := yaml.Unmarshal(encoded, &after); err != nil {
		t.Fatal(err)
	}
	if _, ok := after["value"].(float64); !ok {
		t.Fatalf("after: value = %T (%v), want float64; encoded = %q",
			after["value"], after["value"], string(encoded))
	}
}
