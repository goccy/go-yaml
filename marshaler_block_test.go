package yaml_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
)

type blockDocument string

func (d blockDocument) MarshalYAML() ([]byte, error) {
	return []byte(d), nil
}

func TestEncoder_MarshalerBlockIndentation(t *testing.T) {
	for _, doc := range []string{
		"message: |-\n  first\n  second\n",
		"|-\n  first\n  second\n",
		"message: |-\n",
		"message: |\n  first\n\n    indented\n  last\n",
		"message: >-\n  first\n  second\n\n  last\n",
		"message: |2-\n    indented\n  last\n",
	} {
		for _, depth := range []int{0, 1, 3} {
			for _, sequence := range []bool{false, true} {
				t.Run(fmt.Sprintf("%q/depth=%d/sequence=%t", doc, depth, sequence), func(t *testing.T) {
					var expected any
					if err := yaml.Unmarshal([]byte(doc), &expected); err != nil {
						t.Fatal(err)
					}
					var value any = blockDocument(doc)
					for i := 0; i < depth; i++ {
						if sequence {
							value = []any{value}
							expected = []any{expected}
						} else {
							value = map[string]any{"container": value}
							expected = map[string]any{"container": expected}
						}
					}
					for _, custom := range []bool{false, true} {
						var opts []yaml.EncodeOption
						if custom {
							opts = append(opts, yaml.CustomMarshaler[blockDocument](func(d blockDocument) ([]byte, error) {
								return []byte(d), nil
							}))
						}
						encoded, err := yaml.MarshalWithOptions(value, opts...)
						if err != nil {
							t.Fatal(err)
						}
						var actual any
						if err := yaml.Unmarshal(encoded, &actual); err != nil {
							t.Fatalf("decode emitted YAML: %v\n%s", err, encoded)
						}
						if !reflect.DeepEqual(actual, expected) {
							t.Fatalf("want %#v, got %#v\n%s", expected, actual, encoded)
						}
					}
				})
			}
		}
	}
}

func TestBlockDocument_AddColumn(t *testing.T) {
	source := "message: |-\n  first\n\n    indented\n  last\n"
	file, err := parser.ParseBytes([]byte(source), 0)
	if err != nil {
		t.Fatal(err)
	}
	node := file.Docs[0].Body
	original := node.String()
	node.AddColumn(0)
	if node.String() != original {
		t.Fatal("zero column shift changed document")
	}
	node.AddColumn(2)
	node.AddColumn(4)
	var actual, expected any
	if err := yaml.Unmarshal([]byte(source), &expected); err != nil {
		t.Fatal(err)
	}
	if err := yaml.Unmarshal([]byte(node.String()), &actual); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("column shift changed value: %#v", actual)
	}
	node.AddColumn(-4)
	node.AddColumn(-2)
	if node.String() != original {
		t.Fatalf("shift did not restore original block:\n%s", node.String())
	}
}
