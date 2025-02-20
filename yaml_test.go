package yaml_test

import (
	"encoding/json"
	"fmt"
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

func checkRawValue[T any](t *testing.T, v yaml.RawMessage, expected T) {
	t.Helper()

	var actual T

	if err := yaml.Unmarshal(v, &actual); err != nil {
		t.Errorf("failed to unmarshal: %v", err)
		return
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func checkJSONRawValue[T any](t *testing.T, v json.RawMessage, expected T) {
	t.Helper()

	var actual T

	if err := json.Unmarshal(v, &actual); err != nil {
		t.Errorf("failed to unmarshal: %v", err)
		return
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}

	checkRawValue(t, yaml.RawMessage(v), expected)
}

func TestRawMessage(t *testing.T) {
	var data = []byte(`
a: 1
b: "asdf"
c:
  foo: bar
`)

	var m map[string]yaml.RawMessage
	if err := yaml.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if len(m) != 3 {
		t.Fatalf("failed to decode: %d", len(m))
	}

	checkRawValue(t, m["a"], 1)
	checkRawValue(t, m["b"], "asdf")
	checkRawValue(t, m["c"], map[string]string{"foo": "bar"})

	dt, err := yaml.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	var m2 map[string]yaml.RawMessage
	if err := yaml.Unmarshal(dt, &m2); err != nil {
		t.Fatal(err)
	}

	checkRawValue(t, m2["a"], 1)
	checkRawValue(t, m2["b"], "asdf")
	checkRawValue(t, m2["c"], map[string]string{"foo": "bar"})

	dt, err = json.Marshal(m2)
	if err != nil {
		t.Fatal(err)
	}

	var m3 map[string]yaml.RawMessage
	if err := yaml.Unmarshal(dt, &m3); err != nil {
		t.Fatal(err)
	}
	checkRawValue(t, m3["a"], 1)
	checkRawValue(t, m3["b"], "asdf")
	checkRawValue(t, m3["c"], map[string]string{"foo": "bar"})

	var m4 map[string]json.RawMessage
	if err := json.Unmarshal(dt, &m4); err != nil {
		t.Fatal(err)
	}
	checkJSONRawValue(t, m4["a"], 1)
	checkJSONRawValue(t, m4["b"], "asdf")
	checkJSONRawValue(t, m4["c"], map[string]string{"foo": "bar"})
}
