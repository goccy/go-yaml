package yaml_test

import (
	"testing"

	"github.com/goccy/go-yaml"
)

func TestMarshal(t *testing.T) {
	var v struct {
		A int
		B string
	}
	v.A = 1
	v.B = "hello"
	bytes, err := yaml.Marshal(v)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if string(bytes) != "a: 1\nb: hello\n" {
		t.Fatal("failed to marshal")
	}
}

func TestUnmarshal(t *testing.T) {
	yml := `
%YAML 1.2
---
a: 1
b: c
`
	var v struct {
		A int
		B string
	}
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		t.Fatalf("%+v", err)
	}
}
