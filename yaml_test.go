package yaml_test

import (
	"encoding/json"
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
	data := `
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
	data := `
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
	data := []byte(`
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

type rawYAMLWrapper struct {
	StaticField  string          `json:"staticField" yaml:"staticField"`
	DynamicField yaml.RawMessage `json:"dynamicField" yaml:"dynamicField"`
}

type rawJSONWrapper struct {
	StaticField  string          `json:"staticField" yaml:"staticField"`
	DynamicField json.RawMessage `json:"dynamicField" yaml:"dynamicField"`
}

func (w rawJSONWrapper) Equals(o *rawJSONWrapper) bool {
	if w.StaticField != o.StaticField {
		return false
	}
	return reflect.DeepEqual(w.DynamicField, o.DynamicField)
}

type dynamicField struct {
	A int               `json:"a" yaml:"a"`
	B string            `json:"b" yaml:"b"`
	C map[string]string `json:"c" yaml:"c"`
}

func (t dynamicField) Equals(o *dynamicField) bool {
	if t.A != o.A {
		return false
	}
	if t.B != o.B {
		return false
	}
	if len(t.C) != len(o.C) {
		return false
	}
	for k, v := range t.C {
		ov, exists := o.C[k]
		if !exists {
			return false
		}
		if v != ov {
			return false
		}
	}
	return true
}

func TestRawMessageJSONCompatibility(t *testing.T) {
	rawData := []byte(`staticField: value
dynamicField:
  a: 1
  b: abcd
  c:
    foo: bar
    something: else
`)

	expectedDynamicFieldValue := &dynamicField{
		A: 1,
		B: "abcd",
		C: map[string]string{
			"foo":       "bar",
			"something": "else",
		},
	}

	t.Run("UseJSONUnmarshaler and json.RawMessage", func(t *testing.T) {
		var wrapper rawJSONWrapper
		if err := yaml.UnmarshalWithOptions(rawData, &wrapper, yaml.UseJSONUnmarshaler()); err != nil {
			t.Fatal(err)
		}
		if wrapper.StaticField != "value" {
			t.Fatalf("unexpected wrapper static field value: %s", wrapper.StaticField)
		}
		var dynamicFieldValue dynamicField
		if err := yaml.Unmarshal(wrapper.DynamicField, &dynamicFieldValue); err != nil {
			t.Fatal(err)
		}
		if !dynamicFieldValue.Equals(expectedDynamicFieldValue) {
			t.Fatalf("unexpected dynamic field value: %v", dynamicFieldValue)
		}
	})

	t.Run("UseJSONUnmarshaler and yaml.RawMessage", func(t *testing.T) {
		var wrapper rawYAMLWrapper
		if err := yaml.UnmarshalWithOptions(rawData, &wrapper, yaml.UseJSONUnmarshaler()); err != nil {
			t.Fatal(err)
		}
		if wrapper.StaticField != "value" {
			t.Fatalf("unexpected wrapper static field value: %s", wrapper.StaticField)
		}
		var dynamicFieldValue dynamicField
		if err := yaml.Unmarshal(wrapper.DynamicField, &dynamicFieldValue); err != nil {
			t.Fatal(err)
		}
		if !dynamicFieldValue.Equals(expectedDynamicFieldValue) {
			t.Fatalf("unexpected dynamic field value: %v", dynamicFieldValue)
		}
	})

	t.Run("UseJSONMarshaler and json.RawMessage", func(t *testing.T) {
		dynamicFieldBytes, err := yaml.Marshal(expectedDynamicFieldValue)
		if err != nil {
			t.Fatal(err)
		}
		wrapper := rawJSONWrapper{
			StaticField:  "value",
			DynamicField: json.RawMessage(dynamicFieldBytes),
		}
		wrapperBytes, err := yaml.MarshalWithOptions(&wrapper, yaml.UseJSONMarshaler())
		if err != nil {
			t.Fatal(err)
		}
		var unmarshaledWrapper rawJSONWrapper
		if err := yaml.UnmarshalWithOptions(wrapperBytes, &unmarshaledWrapper, yaml.UseJSONUnmarshaler()); err != nil {
			t.Fatal(err)
		}
		if unmarshaledWrapper.StaticField != wrapper.StaticField {
			t.Fatalf("unexpected unmarshaled static field value: %s", unmarshaledWrapper.StaticField)
		}
		var unmarshaledDynamicFieldValue dynamicField
		if err := yaml.UnmarshalWithOptions(unmarshaledWrapper.DynamicField, &unmarshaledDynamicFieldValue, yaml.UseJSONUnmarshaler()); err != nil {
			t.Fatal(err)
		}
		if !unmarshaledDynamicFieldValue.Equals(expectedDynamicFieldValue) {
			t.Fatalf("unexpected unmarshaled dynamic field value: %v", unmarshaledDynamicFieldValue)
		}
	})

	t.Run("UseJSONMarshaler and yaml.RawMessage", func(t *testing.T) {
		dynamicFieldBytes, err := yaml.Marshal(expectedDynamicFieldValue)
		if err != nil {
			t.Fatal(err)
		}
		wrapper := rawYAMLWrapper{
			StaticField:  "value",
			DynamicField: yaml.RawMessage(dynamicFieldBytes),
		}
		wrapperBytes, err := yaml.MarshalWithOptions(&wrapper, yaml.UseJSONMarshaler())
		if err != nil {
			t.Fatal(err)
		}
		var unmarshaledWrapper rawYAMLWrapper
		if err := yaml.UnmarshalWithOptions(wrapperBytes, &unmarshaledWrapper, yaml.UseJSONUnmarshaler()); err != nil {
			t.Fatal(err)
		}
		if unmarshaledWrapper.StaticField != wrapper.StaticField {
			t.Fatalf("unexpected unmarshaled static field value: %s", unmarshaledWrapper.StaticField)
		}
		var unmarshaledDynamicFieldValue dynamicField
		if err := yaml.UnmarshalWithOptions(unmarshaledWrapper.DynamicField, &unmarshaledDynamicFieldValue, yaml.UseJSONUnmarshaler()); err != nil {
			t.Fatal(err)
		}
		if !unmarshaledDynamicFieldValue.Equals(expectedDynamicFieldValue) {
			t.Fatalf("unexpected unmarshaled dynamic field value: %v", unmarshaledDynamicFieldValue)
		}
	})
}
