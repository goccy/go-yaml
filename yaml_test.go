package yaml_test

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
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
# comment
---
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

func TestDecodeKeepAddress(t *testing.T) {
	var data = `
a: &a [_]
b: &b [*a,*a]
c: &c [*b,*b]
d: &d [*c,*c]
`
	var v map[string]any
	if err := yaml.Unmarshal([]byte(data), &v); err != nil {
		t.Fatal(err)
	}
	a := v["a"]
	b := v["b"]
	for _, vv := range v["b"].([]any) {
		if fmt.Sprintf("%p", a) != fmt.Sprintf("%p", vv) {
			t.Fatalf("failed to decode b element as keep address")
		}
	}
	for _, vv := range v["c"].([]any) {
		if fmt.Sprintf("%p", b) != fmt.Sprintf("%p", vv) {
			t.Fatalf("failed to decode c element as keep address")
		}
	}
}

func TestSmartAnchor(t *testing.T) {
	var data = `
a: &a [_,_,_,_,_,_,_,_,_,_,_,_,_,_,_]
b: &b [*a,*a,*a,*a,*a,*a,*a,*a,*a,*a]
c: &c [*b,*b,*b,*b,*b,*b,*b,*b,*b,*b]
d: &d [*c,*c,*c,*c,*c,*c,*c,*c,*c,*c]
e: &e [*d,*d,*d,*d,*d,*d,*d,*d,*d,*d]
f: &f [*e,*e,*e,*e,*e,*e,*e,*e,*e,*e]
g: &g [*f,*f,*f,*f,*f,*f,*f,*f,*f,*f]
h: &h [*g,*g,*g,*g,*g,*g,*g,*g,*g,*g]
i: &i [*h,*h,*h,*h,*h,*h,*h,*h,*h,*h]
`
	var v any
	if err := yaml.Unmarshal([]byte(data), &v); err != nil {
		t.Fatal(err)
	}
	got, err := yaml.MarshalWithOptions(v, yaml.WithSmartAnchor())
	if err != nil {
		t.Fatal(err)
	}
	expected := `
a: &a
- _
- _
- _
- _
- _
- _
- _
- _
- _
- _
- _
- _
- _
- _
- _
b: &b
- *a
- *a
- *a
- *a
- *a
- *a
- *a
- *a
- *a
- *a
c: &c
- *b
- *b
- *b
- *b
- *b
- *b
- *b
- *b
- *b
- *b
d: &d
- *c
- *c
- *c
- *c
- *c
- *c
- *c
- *c
- *c
- *c
e: &e
- *d
- *d
- *d
- *d
- *d
- *d
- *d
- *d
- *d
- *d
f: &f
- *e
- *e
- *e
- *e
- *e
- *e
- *e
- *e
- *e
- *e
g: &g
- *f
- *f
- *f
- *f
- *f
- *f
- *f
- *f
- *f
- *f
h: &h
- *g
- *g
- *g
- *g
- *g
- *g
- *g
- *g
- *g
- *g
i:
- *h
- *h
- *h
- *h
- *h
- *h
- *h
- *h
- *h
- *h
`
	if strings.TrimPrefix(expected, "\n") != string(got) {
		t.Fatalf("failed to encode: %s", string(got))
	}
}

func TestCustomErrorMessage(t *testing.T) {
	data := `
foo:
  bar:
    foo: 2
  baz:
    foo: 3
foo: 2
`
	if _, err := parser.ParseBytes([]byte(data), 0); err == nil {
		t.Fatalf("expected error")
	} else {
		yamlErr, ok := err.(yaml.Error)
		if !ok {
			t.Fatalf("failed to get yaml.Error from error: %T", err)
		}
		expected := `
[7:1] custom message
   4 |     foo: 2
   5 |   baz:
   6 |     foo: 3
>  7 | foo: 2
       ^
`
		got := "\n" + yaml.FormatErrorWithToken("custom message", yamlErr.GetToken(), false, true)
		if expected != got {
			t.Fatalf("unexpected error message:\nexpected:\n%s\nbut got:\n%s", expected, got)
		}
	}
}
