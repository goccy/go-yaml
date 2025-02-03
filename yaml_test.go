package yaml_test

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestRoundTripWithComment(t *testing.T) {
	yml := `
# head comment
key: value # line comment
`
	var v struct {
		Key string
	}
	comments := yaml.CommentMap{}

	if err := yaml.UnmarshalWithOptions([]byte(yml), &v, yaml.Strict(), yaml.CommentToMap(comments)); err != nil {
		t.Fatal(err)
	}
	out, err := yaml.MarshalWithOptions(v, yaml.WithComment(comments))
	if err != nil {
		t.Fatal(err)
	}
	got := "\n" + string(out)
	if yml != got {
		t.Fatalf("failed to get round tripped yaml: %s", got)
	}
}

func TestStreamDecodingWithComment(t *testing.T) {
	yml := `
a:
  b:
    c: # comment
---
foo: bar # comment
---
- a
- b
- c # comment
`
	cm := yaml.CommentMap{}
	dec := yaml.NewDecoder(strings.NewReader(yml), yaml.CommentToMap(cm))
	var commentPathsWithDocIndex [][]string
	for {
		var v any
		if err := dec.Decode(&v); err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
		paths := make([]string, 0, len(cm))
		for k := range cm {
			paths = append(paths, k)
		}
		commentPathsWithDocIndex = append(commentPathsWithDocIndex, paths)
		for k := range cm {
			delete(cm, k)
		}
	}
	if !reflect.DeepEqual(commentPathsWithDocIndex, [][]string{
		{"$.a.b.c"},
		{"$.foo"},
		{"$[2]"},
	}) {
		t.Fatalf("failed to get comment: %v", commentPathsWithDocIndex)
	}
}
