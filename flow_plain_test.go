package yaml_test

import (
	"testing"

	"github.com/goccy/go-yaml"
)

func TestFlowPlainScalarContinuation(t *testing.T) {
	t.Run("folded when indented less than the key", func(t *testing.T) {
		// A continuation only has to clear the enclosing block indent.
		// Comparing it to the key column splits "x y" into two values.
		var v struct {
			K struct {
				A string `yaml:"a"`
			} `yaml:"k"`
		}
		if err := yaml.Unmarshal([]byte("k: {a: x\n  y}\n"), &v); err != nil {
			t.Fatal(err)
		}
		if v.K.A != "x y" {
			t.Fatalf("got %q", v.K.A)
		}
	})

	t.Run("still folded when indented past the key", func(t *testing.T) {
		var v struct {
			K struct {
				A string `yaml:"a"`
			} `yaml:"k"`
		}
		if err := yaml.Unmarshal([]byte("k: {a: x\n     y}\n"), &v); err != nil {
			t.Fatal(err)
		}
		if v.K.A != "x y" {
			t.Fatalf("got %q", v.K.A)
		}
	})

	t.Run("folded inside a nested block", func(t *testing.T) {
		var v struct {
			Parent struct {
				K struct {
					A string `yaml:"a"`
				} `yaml:"k"`
			} `yaml:"parent"`
		}
		src := "parent:\n  k: {a: x\n    y}\n"
		if err := yaml.Unmarshal([]byte(src), &v); err != nil {
			t.Fatal(err)
		}
		if v.Parent.K.A != "x y" {
			t.Fatalf("got %q", v.Parent.K.A)
		}
	})

	t.Run("flow sequence still folds", func(t *testing.T) {
		var v struct {
			K []string `yaml:"k"`
		}
		if err := yaml.Unmarshal([]byte("k: [x\n  y]\n"), &v); err != nil {
			t.Fatal(err)
		}
		if len(v.K) != 1 || v.K[0] != "x y" {
			t.Fatalf("got %#v", v.K)
		}
	})

	t.Run("comma still separates a later key", func(t *testing.T) {
		var v struct {
			K struct {
				A string `yaml:"a"`
				B string `yaml:"b"`
			} `yaml:"k"`
		}
		if err := yaml.Unmarshal([]byte("k: {a: x,\n  b: y}\n"), &v); err != nil {
			t.Fatal(err)
		}
		if v.K.A != "x" || v.K.B != "y" {
			t.Fatalf("got %#v", v.K)
		}
	})

	t.Run("multiline flow key still folds", func(t *testing.T) {
		var v []map[string]string
		if err := yaml.Unmarshal([]byte("- { multi\n  line: value}\n"), &v); err != nil {
			t.Fatal(err)
		}
		if len(v) != 1 || v[0]["multi line"] != "value" {
			t.Fatalf("got %#v", v)
		}
	})

	t.Run("next key without a comma is still an error", func(t *testing.T) {
		var v any
		if err := yaml.Unmarshal([]byte("k: {a: x\n  b: y}\n"), &v); err == nil {
			t.Fatalf("expected error, got %#v", v)
		}
	})
}
