package yaml_test

import (
	"bytes"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"golang.org/x/xerrors"
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

type marshalTest struct{}

func (t *marshalTest) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(yaml.MapSlice{
		{
			"a", 1,
		},
		{
			"b", "hello",
		},
		{
			"c", true,
		},
		{
			"d", map[string]string{"x": "y"},
		},
	})
}

type marshalTest2 struct{}

func (t *marshalTest2) MarshalYAML() (interface{}, error) {
	return yaml.MapSlice{
		{
			"a", 2,
		},
		{
			"b", "world",
		},
		{
			"c", true,
		},
	}, nil
}

func TestMarshalYAML(t *testing.T) {
	var v struct {
		A *marshalTest
		B *marshalTest2
	}
	v.A = &marshalTest{}
	v.B = &marshalTest2{}
	bytes, err := yaml.Marshal(v)
	if err != nil {
		t.Fatalf("failed to Marshal: %+v", err)
	}
	expect := `
a:
  a: 1
  b: hello
  c: true
  d:
    x: y
b:
  a: 2
  b: world
  c: true
`
	actual := "\n" + string(bytes)
	if expect != actual {
		t.Fatalf("failed to MarshalYAML expect:[%s], actual:[%s]", expect, actual)
	}
}

type unmarshalTest struct {
	a int
	b string
	c bool
}

func (t *unmarshalTest) UnmarshalYAML(b []byte) error {
	if t.a != 0 {
		return xerrors.New("unexpected field value to a")
	}
	if t.b != "" {
		return xerrors.New("unexpected field value to b")
	}
	if t.c {
		return xerrors.New("unexpected field value to c")
	}
	var v struct {
		A int
		B string
		C bool
	}
	if err := yaml.Unmarshal(b, &v); err != nil {
		return err
	}
	t.a = v.A
	t.b = v.B
	t.c = v.C
	return nil
}

type unmarshalTest2 struct {
	a int
	b string
	c bool
}

func (t *unmarshalTest2) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v struct {
		A int
		B string
		C bool
	}
	if t.a != 0 {
		return xerrors.New("unexpected field value to a")
	}
	if t.b != "" {
		return xerrors.New("unexpected field value to b")
	}
	if t.c {
		return xerrors.New("unexpected field value to c")
	}
	if err := unmarshal(&v); err != nil {
		return err
	}
	t.a = v.A
	t.b = v.B
	t.c = v.C
	return nil
}

func TestUnmarshalYAML(t *testing.T) {
	yml := `
a:
  a: 1
  b: hello
  c: true
b:
  a: 2
  b: world
  c: true
`
	var v struct {
		A *unmarshalTest
		B *unmarshalTest2
	}
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		t.Fatalf("failed to Unmarshal: %+v", err)
	}
	if v.A == nil {
		t.Fatal("failed to UnmarshalYAML")
	}
	if v.A.a != 1 {
		t.Fatal("failed to UnmarshalYAML")
	}
	if v.A.b != "hello" {
		t.Fatal("failed to UnmarshalYAML")
	}
	if !v.A.c {
		t.Fatal("failed to UnmarshalYAML")
	}
	if v.B == nil {
		t.Fatal("failed to UnmarshalYAML")
	}
	if v.B.a != 2 {
		t.Fatal("failed to UnmarshalYAML")
	}
	if v.B.b != "world" {
		t.Fatal("failed to UnmarshalYAML")
	}
	if !v.B.c {
		t.Fatal("failed to UnmarshalYAML")
	}
}

type ObjectMap map[string]*Object
type ObjectDecl struct {
	Name    string `yaml:"-"`
	*Object `yaml:",inline,anchor"`
}

func (m ObjectMap) MarshalYAML() (interface{}, error) {
	newMap := map[string]*ObjectDecl{}
	for k, v := range m {
		newMap[k] = &ObjectDecl{Name: k, Object: v}
	}
	return newMap, nil
}

type rootObject struct {
	Single     ObjectMap            `yaml:"single"`
	Collection map[string][]*Object `yaml:"collection"`
}

type Object struct {
	*Object  `yaml:",omitempty,inline,alias"`
	MapValue map[string]interface{} `yaml:",omitempty,inline"`
}

func TestInlineAnchorAndAlias(t *testing.T) {
	yml := `---
single:
  default: &default
    id: 1
    name: john
  user_1: &user_1
    id: 1
    name: ken
  user_2: &user_2
    <<: *default
    id: 2
collection:
  defaults:
  - *default
  - <<: *default
  - <<: *default
    id: 2
  users:
  - <<: *user_1
  - <<: *user_2
  - <<: *user_1
    id: 3
  - <<: *user_1
    id: 4
  - <<: *user_1
    id: 5
`
	var v rootObject
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		panic(err)
	}
	opt := yaml.MarshalAnchor(func(anchor *ast.AnchorNode, value interface{}) error {
		if o, ok := value.(*ObjectDecl); ok {
			anchor.Name.(*ast.StringNode).Value = o.Name
		}
		return nil
	})
	var buf bytes.Buffer
	if err := yaml.NewEncoder(&buf, opt).Encode(v); err != nil {
		t.Fatalf("%+v", err)
	}
	actual := "---\n" + buf.String()
	if yml != actual {
		t.Fatalf("failed to marshal: expected:[%s] actual:[%s]", yml, actual)
	}
}

func TestMapSlice_Map(t *testing.T) {
	yml := `
a: b
c: d
`
	var v yaml.MapSlice
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		t.Fatal(err)
	}
	m := v.ToMap()
	if len(m) != 2 {
		t.Fatal("failed to convert MapSlice to map")
	}
	if m["a"] != "b" {
		t.Fatal("failed to convert MapSlice to map")
	}
	if m["c"] != "d" {
		t.Fatal("failed to convert MapSlice to map")
	}
}
