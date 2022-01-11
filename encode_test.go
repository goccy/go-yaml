package yaml_test

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

var zero = 0
var emptyStr = ""

func TestEncoder(t *testing.T) {
	tests := []struct {
		source  string
		value   interface{}
		options []yaml.EncodeOption
	}{
		{
			"null\n",
			(*struct{})(nil),
			nil,
		},
		{
			"v: hi\n",
			map[string]string{"v": "hi"},
			nil,
		},
		{
			"v: \"true\"\n",
			map[string]string{"v": "true"},
			nil,
		},
		{
			"v: \"false\"\n",
			map[string]string{"v": "false"},
			nil,
		},
		{
			"v: true\n",
			map[string]interface{}{"v": true},
			nil,
		},
		{
			"v: false\n",
			map[string]bool{"v": false},
			nil,
		},
		{
			"v: 10\n",
			map[string]int{"v": 10},
			nil,
		},
		{
			"v: -10\n",
			map[string]int{"v": -10},
			nil,
		},
		{
			"v: 4294967296\n",
			map[string]int{"v": 4294967296},
			nil,
		},
		{
			"v: 0.1\n",
			map[string]interface{}{"v": 0.1},
			nil,
		},
		{
			"v: 0.99\n",
			map[string]float32{"v": 0.99},
			nil,
		},
		{
			"v: 0.123456789\n",
			map[string]float64{"v": 0.123456789},
			nil,
		},
		{
			"v: -0.1\n",
			map[string]float64{"v": -0.1},
			nil,
		},
		{
			"v: 1.0\n",
			map[string]float64{"v": 1.0},
			nil,
		},
		{
			"v: 1e+06\n",
			map[string]float64{"v": 1000000},
			nil,
		},
		{
			"v: .inf\n",
			map[string]interface{}{"v": math.Inf(0)},
			nil,
		},
		{
			"v: -.inf\n",
			map[string]interface{}{"v": math.Inf(-1)},
			nil,
		},
		{
			"v: .nan\n",
			map[string]interface{}{"v": math.NaN()},
			nil,
		},
		{
			"v: null\n",
			map[string]interface{}{"v": nil},
			nil,
		},
		{
			"v: \"\"\n",
			map[string]string{"v": ""},
			nil,
		},
		{
			"v:\n- A\n- B\n",
			map[string][]string{"v": {"A", "B"}},
			nil,
		},
		{
			"v:\n  - A\n  - B\n",
			map[string][]string{"v": {"A", "B"}},
			[]yaml.EncodeOption{
				yaml.IndentSequence(true),
			},
		},
		{
			"v:\n- A\n- B\n",
			map[string][2]string{"v": {"A", "B"}},
			nil,
		},
		{
			"v:\n  - A\n  - B\n",
			map[string][2]string{"v": {"A", "B"}},
			[]yaml.EncodeOption{
				yaml.IndentSequence(true),
			},
		},
		{
			"a: -\n",
			map[string]string{"a": "-"},
			nil,
		},
		{
			"123\n",
			123,
			nil,
		},
		{
			"hello: world\n",
			map[string]string{"hello": "world"},
			nil,
		},
		{
			"hello: |\n  hello\n  world\n",
			map[string]string{"hello": "hello\nworld\n"},
			nil,
		},
		{
			"hello: |-\n  hello\n  world\n",
			map[string]string{"hello": "hello\nworld"},
			nil,
		},
		{
			"hello: |+\n  hello\n  world\n\n",
			map[string]string{"hello": "hello\nworld\n\n"},
			nil,
		},
		{
			"hello:\n  hello: |\n    hello\n    world\n",
			map[string]map[string]string{"hello": {"hello": "hello\nworld\n"}},
			nil,
		},
		{
			"hello: |\r  hello\r  world\n",
			map[string]string{"hello": "hello\rworld\r"},
			nil,
		},
		{
			"hello: |\r\n  hello\r\n  world\n",
			map[string]string{"hello": "hello\r\nworld\r\n"},
			nil,
		},
		{
			"v: |-\n  username: hello\n  password: hello123\n",
			map[string]interface{}{"v": "username: hello\npassword: hello123"},
			[]yaml.EncodeOption{
				yaml.UseLiteralStyleIfMultiline(true),
			},
		},
		{
			"v: |-\n  # comment\n  username: hello\n  password: hello123\n",
			map[string]interface{}{"v": "# comment\nusername: hello\npassword: hello123"},
			[]yaml.EncodeOption{
				yaml.UseLiteralStyleIfMultiline(true),
			},
		},
		{
			"v: \"# comment\\nusername: hello\\npassword: hello123\"\n",
			map[string]interface{}{"v": "# comment\nusername: hello\npassword: hello123"},
			[]yaml.EncodeOption{
				yaml.UseLiteralStyleIfMultiline(false),
			},
		},
		{
			"v:\n- A\n- 1\n- B:\n  - 2\n  - 3\n",
			map[string]interface{}{
				"v": []interface{}{
					"A",
					1,
					map[string][]int{
						"B": {2, 3},
					},
				},
			},
			nil,
		},
		{
			"v:\n  - A\n  - 1\n  - B:\n      - 2\n      - 3\n  - 2\n",
			map[string]interface{}{
				"v": []interface{}{
					"A",
					1,
					map[string][]int{
						"B": {2, 3},
					},
					2,
				},
			},
			[]yaml.EncodeOption{
				yaml.IndentSequence(true),
			},
		},
		{
			"a:\n  b: c\n",
			map[string]interface{}{
				"a": map[string]string{
					"b": "c",
				},
			},
			nil,
		},
		{
			"t2: 2018-01-09T10:40:47Z\nt4: 2098-01-09T10:40:47Z\n",
			map[string]string{
				"t2": "2018-01-09T10:40:47Z",
				"t4": "2098-01-09T10:40:47Z",
			},
			nil,
		},
		{
			"a:\n  b: c\n  d: e\n",
			map[string]interface{}{
				"a": map[string]string{
					"b": "c",
					"d": "e",
				},
			},
			nil,
		},
		{
			"a: 3s\n",
			map[string]string{
				"a": "3s",
			},
			nil,
		},
		{
			"a: <foo>\n",
			map[string]string{"a": "<foo>"},
			nil,
		},
		{
			"a: \"1:1\"\n",
			map[string]string{"a": "1:1"},
			nil,
		},
		{
			"a: 1.2.3.4\n",
			map[string]string{"a": "1.2.3.4"},
			nil,
		},
		{
			"a: \"b: c\"\n",
			map[string]string{"a": "b: c"},
			nil,
		},
		{
			"a: \"Hello #comment\"\n",
			map[string]string{"a": "Hello #comment"},
			nil,
		},
		{
			"a: 100.5\n",
			map[string]interface{}{
				"a": 100.5,
			},
			nil,
		},
		{
			"a: \"\\\\0\"\n",
			map[string]string{"a": "\\0"},
			nil,
		},
		{
			"a: 1\nb: 2\nc: 3\nd: 4\nsub:\n  e: 5\n",
			map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
				"sub": map[string]int{
					"e": 5,
				},
			},
			nil,
		},
		{
			"a: 1\nb: []\n",
			struct {
				A int
				B []string
			}{
				1, ([]string)(nil),
			},
			nil,
		},
		{
			"a: 1\nb: []\n",
			struct {
				A int
				B []string
			}{
				1, []string{},
			},
			nil,
		},
		{
			"a: {}\n",
			struct {
				A map[string]interface{}
			}{
				map[string]interface{}{},
			},
			nil,
		},
		{
			"a: b\nc: d\n",
			struct {
				A string
				C string `yaml:"c"`
			}{
				"b", "d",
			},
			nil,
		},
		{
			"a: 1\n",
			struct {
				A int
				B int `yaml:"-"`
			}{
				1, 0,
			},
			nil,
		},
		{
			"a: \"\"\n",
			struct {
				A string
			}{
				"",
			},
			nil,
		},
		{
			"a: null\n",
			struct {
				A *string
			}{
				nil,
			},
			nil,
		},
		{
			"a: \"\"\n",
			struct {
				A *string
			}{
				&emptyStr,
			},
			nil,
		},
		{
			"a: null\n",
			struct {
				A *int
			}{
				nil,
			},
			nil,
		},
		{
			"a: 0\n",
			struct {
				A *int
			}{
				&zero,
			},
			nil,
		},

		// Conditional flag
		{
			"a: 1\n",
			struct {
				A int `yaml:"a,omitempty"`
				B int `yaml:"b,omitempty"`
			}{1, 0},
			nil,
		},
		{
			"{}\n",
			struct {
				A int `yaml:"a,omitempty"`
				B int `yaml:"b,omitempty"`
			}{0, 0},
			nil,
		},

		{
			"a:\n  y: \"\"\n",
			struct {
				A *struct {
					X string `yaml:"x,omitempty"`
					Y string
				}
			}{&struct {
				X string `yaml:"x,omitempty"`
				Y string
			}{}},
			nil,
		},

		{
			"a: {}\n",
			struct {
				A *struct {
					X string `yaml:"x,omitempty"`
					Y string `yaml:"y,omitempty"`
				}
			}{&struct {
				X string `yaml:"x,omitempty"`
				Y string `yaml:"y,omitempty"`
			}{}},
			nil,
		},

		{
			"a: {x: 1}\n",
			struct {
				A *struct{ X, y int } `yaml:"a,omitempty,flow"`
			}{&struct{ X, y int }{1, 2}},
			nil,
		},

		{
			"{}\n",
			struct {
				A *struct{ X, y int } `yaml:"a,omitempty,flow"`
			}{nil},
			nil,
		},

		{
			"a: {x: 0}\n",
			struct {
				A *struct{ X, y int } `yaml:"a,omitempty,flow"`
			}{&struct{ X, y int }{}},
			nil,
		},

		{
			"a: {x: 1}\n",
			struct {
				A struct{ X, y int } `yaml:"a,omitempty,flow"`
			}{struct{ X, y int }{1, 2}},
			nil,
		},
		{
			"{}\n",
			struct {
				A struct{ X, y int } `yaml:"a,omitempty,flow"`
			}{struct{ X, y int }{0, 1}},
			nil,
		},
		{
			"a: 1.0\n",
			struct {
				A float64 `yaml:"a,omitempty"`
				B float64 `yaml:"b,omitempty"`
			}{1, 0},
			nil,
		},
		{
			"a: 1\n",
			struct {
				A int
				B []string `yaml:"b,omitempty"`
			}{
				1, []string{},
			},
			nil,
		},

		// Flow flag
		{
			"a: [1, 2]\n",
			struct {
				A []int `yaml:"a,flow"`
			}{[]int{1, 2}},
			nil,
		},
		{
			"a: {b: c, d: e}\n",
			&struct {
				A map[string]string `yaml:"a,flow"`
			}{map[string]string{"b": "c", "d": "e"}},
			nil,
		},
		{
			"a: {b: c, d: e}\n",
			struct {
				A struct {
					B, D string
				} `yaml:"a,flow"`
			}{struct{ B, D string }{"c", "e"}},
			nil,
		},

		// Multi bytes
		{
			"v: あいうえお\nv2: かきくけこ\n",
			map[string]string{"v": "あいうえお", "v2": "かきくけこ"},
			nil,
		},

		// time value
		{
			"v: 0001-01-01T00:00:00Z\n",
			map[string]time.Time{"v": time.Time{}},
			nil,
		},
		{
			"v: 0001-01-01T00:00:00Z\n",
			map[string]*time.Time{"v": &time.Time{}},
			nil,
		},
		{
			"v: null\n",
			map[string]*time.Time{"v": nil},
			nil,
		},
		{
			"v: 30s\n",
			map[string]time.Duration{"v": 30 * time.Second},
			nil,
		},
		// Quote style
		{
			`v: '\'a\'b'` + "\n",
			map[string]string{"v": `'a'b`},
			[]yaml.EncodeOption{
				yaml.UseSingleQuote(true),
			},
		},
		{
			`v: "'a'b"` + "\n",
			map[string]string{"v": `'a'b`},
			[]yaml.EncodeOption{
				yaml.UseSingleQuote(false),
			},
		},
	}
	for _, test := range tests {
		var buf bytes.Buffer
		enc := yaml.NewEncoder(&buf, test.options...)
		if err := enc.Encode(test.value); err != nil {
			t.Fatalf("%+v", err)
		}
		if test.source != buf.String() {
			t.Fatalf("expect = [%s], actual = [%s]", test.source, buf.String())
		}
	}
}

func TestEncodeStructIncludeMap(t *testing.T) {
	type U struct {
		M map[string]string
	}
	type T struct {
		A U
	}
	bytes, err := yaml.Marshal(T{
		A: U{
			M: map[string]string{"x": "y"},
		},
	})
	if err != nil {
		t.Fatalf("%+v", err)
	}
	expect := "a:\n  m:\n    x: y\n"
	actual := string(bytes)
	if actual != expect {
		t.Fatalf("unexpected output. expect:[%s] actual:[%s]", expect, actual)
	}
}

func TestEncodeDefinedTypeKeyMap(t *testing.T) {
	type K string
	type U struct {
		M map[K]string
	}
	bytes, err := yaml.Marshal(U{
		M: map[K]string{K("x"): "y"},
	})
	if err != nil {
		t.Fatalf("%+v", err)
	}
	expect := "m:\n  x: y\n"
	actual := string(bytes)
	if actual != expect {
		t.Fatalf("unexpected output. expect:[%s] actual:[%s]", expect, actual)
	}
}

func TestEncodeWithAnchorAndAlias(t *testing.T) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	type T struct {
		A int
		B string
	}
	var v struct {
		A *T `yaml:"a,anchor=c"`
		B *T `yaml:"b,alias=c"`
	}
	v.A = &T{A: 1, B: "hello"}
	v.B = v.A
	if err := enc.Encode(v); err != nil {
		t.Fatalf("%+v", err)
	}
	expect := "a: &c\n  a: 1\n  b: hello\nb: *c\n"
	if expect != buf.String() {
		t.Fatalf("expect = [%s], actual = [%s]", expect, buf.String())
	}
}

func TestEncodeWithAutoAlias(t *testing.T) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	type T struct {
		I int
		S string
	}
	var v struct {
		A *T `yaml:"a,anchor=a"`
		B *T `yaml:"b,anchor=b"`
		C *T `yaml:"c,alias"`
		D *T `yaml:"d,alias"`
	}
	v.A = &T{I: 1, S: "hello"}
	v.B = &T{I: 2, S: "world"}
	v.C = v.A
	v.D = v.B
	if err := enc.Encode(v); err != nil {
		t.Fatalf("%+v", err)
	}
	expect := `a: &a
  i: 1
  s: hello
b: &b
  i: 2
  s: world
c: *a
d: *b
`
	if expect != buf.String() {
		t.Fatalf("expect = [%s], actual = [%s]", expect, buf.String())
	}
}

func TestEncodeWithImplicitAnchorAndAlias(t *testing.T) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	type T struct {
		I int
		S string
	}
	var v struct {
		A *T `yaml:"a,anchor"`
		B *T `yaml:"b,anchor"`
		C *T `yaml:"c,alias"`
		D *T `yaml:"d,alias"`
	}
	v.A = &T{I: 1, S: "hello"}
	v.B = &T{I: 2, S: "world"}
	v.C = v.A
	v.D = v.B
	if err := enc.Encode(v); err != nil {
		t.Fatalf("%+v", err)
	}
	expect := `a: &a
  i: 1
  s: hello
b: &b
  i: 2
  s: world
c: *a
d: *b
`
	if expect != buf.String() {
		t.Fatalf("expect = [%s], actual = [%s]", expect, buf.String())
	}
}

func TestEncodeWithMerge(t *testing.T) {
	type Person struct {
		*Person `yaml:",omitempty,inline,alias"`
		Name    string `yaml:",omitempty"`
		Age     int    `yaml:",omitempty"`
	}
	defaultPerson := &Person{
		Name: "John Smith",
		Age:  20,
	}
	people := []*Person{
		{
			Person: defaultPerson,
			Name:   "Ken",
			Age:    10,
		},
		{
			Person: defaultPerson,
		},
	}
	var doc struct {
		Default *Person   `yaml:"default,anchor"`
		People  []*Person `yaml:"people"`
	}
	doc.Default = defaultPerson
	doc.People = people
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	if err := enc.Encode(doc); err != nil {
		t.Fatalf("%+v", err)
	}
	expect := `default: &default
  name: John Smith
  age: 20
people:
- <<: *default
  name: Ken
  age: 10
- <<: *default
`
	if expect != buf.String() {
		t.Fatalf("expect = [%s], actual = [%s]", expect, buf.String())
	}
}

func TestEncodeWithNestedYAML(t *testing.T) {
	// Represents objects containing stringified YAML, and special chars
	tests := []struct {
		value interface{}
		// If true, expects a different result between when using forced literal style or not
		expectDifferent bool
	}{
		{
			value:           map[string]interface{}{"v": "# comment\nname: hello\npassword: hello123\nspecial: \":ghost:\"\ntext: |\n  nested multiline!"},
			expectDifferent: true,
		},
		{
			value:           map[string]interface{}{"v": "# comment\nusername: hello\npassword: hello123"},
			expectDifferent: true,
		},
		{
			value:           map[string]interface{}{"v": "# comment\n"},
			expectDifferent: true,
		},
		{
			value: map[string]interface{}{"v": "\n"},
		},
	}

	for _, test := range tests {
		yamlBytesForced, err := yaml.MarshalWithOptions(test.value, yaml.UseLiteralStyleIfMultiline(true))
		if err != nil {
			t.Fatalf("%+v", err)
		}

		// Convert it back for proper equality testing
		var unmarshaled interface{}

		if err := yaml.Unmarshal(yamlBytesForced, &unmarshaled); err != nil {
			t.Fatalf("%+v", err)
		}

		if !reflect.DeepEqual(test.value, unmarshaled) {
			t.Fatalf("expected %v(%T). but actual %v(%T)", test.value, test.value, unmarshaled, unmarshaled)
		}

		if test.expectDifferent {
			yamlBytesNotForced, err := yaml.MarshalWithOptions(test.value)
			if err != nil {
				t.Fatalf("%+v", err)
			}

			if string(yamlBytesForced) == string(yamlBytesNotForced) {
				t.Fatalf("expected different strings when force literal style is not enabled. forced: %s, not forced: %s", string(yamlBytesForced), string(yamlBytesNotForced))
			}
		}
	}
}

func TestEncoder_Inline(t *testing.T) {
	type base struct {
		A int
		B string
	}
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	if err := enc.Encode(struct {
		*base `yaml:",inline"`
		C     bool
	}{
		base: &base{
			A: 1,
			B: "hello",
		},
		C: true,
	}); err != nil {
		t.Fatalf("%+v", err)
	}
	expect := `
a: 1
b: hello
c: true
`
	actual := "\n" + buf.String()
	if expect != actual {
		t.Fatalf("inline marshal error: expect=[%s] actual=[%s]", expect, actual)
	}
}

func TestEncoder_InlineAndConflictKey(t *testing.T) {
	type base struct {
		A int
		B string
	}
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	if err := enc.Encode(struct {
		*base `yaml:",inline"`
		A     int // conflict
		C     bool
	}{
		base: &base{
			A: 1,
			B: "hello",
		},
		A: 0, // default value
		C: true,
	}); err != nil {
		t.Fatalf("%+v", err)
	}
	expect := `
b: hello
a: 0
c: true
`
	actual := "\n" + buf.String()
	if expect != actual {
		t.Fatalf("inline marshal error: expect=[%s] actual=[%s]", expect, actual)
	}
}

func TestEncoder_Flow(t *testing.T) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf, yaml.Flow(true))
	var v struct {
		A int
		B string
		C struct {
			D int
			E string
		}
		F []int
	}
	v.A = 1
	v.B = "hello"
	v.C.D = 3
	v.C.E = "world"
	v.F = []int{1, 2}
	if err := enc.Encode(v); err != nil {
		t.Fatalf("%+v", err)
	}
	expect := `
{a: 1, b: hello, c: {d: 3, e: world}, f: [1, 2]}
`
	actual := "\n" + buf.String()
	if expect != actual {
		t.Fatalf("flow style marshal error: expect=[%s] actual=[%s]", expect, actual)
	}
}

func TestEncoder_FlowRecursive(t *testing.T) {
	var v struct {
		M map[string][]int `yaml:",flow"`
	}
	v.M = map[string][]int{
		"test": []int{1, 2, 3},
	}
	var buf bytes.Buffer
	if err := yaml.NewEncoder(&buf).Encode(v); err != nil {
		t.Fatalf("%+v", err)
	}
	expect := `
m: {test: [1, 2, 3]}
`
	actual := "\n" + buf.String()
	if expect != actual {
		t.Fatalf("flow style marshal error: expect=[%s] actual=[%s]", expect, actual)
	}
}

func TestEncoder_JSON(t *testing.T) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf, yaml.JSON())
	type st struct {
		I int8
		S string
		F float32
	}
	if err := enc.Encode(struct {
		I        int
		U        uint
		S        string
		F        float64
		Struct   *st
		Slice    []int
		Map      map[string]interface{}
		Time     time.Time
		Duration time.Duration
	}{
		I: -10,
		U: 10,
		S: "hello",
		F: 3.14,
		Struct: &st{
			I: 2,
			S: "world",
			F: 1.23,
		},
		Slice: []int{1, 2, 3, 4, 5},
		Map: map[string]interface{}{
			"a": 1,
			"b": 1.23,
			"c": "json",
		},
		Time:     time.Time{},
		Duration: 5 * time.Minute,
	}); err != nil {
		t.Fatalf("%+v", err)
	}
	expect := `
{"i": -10, "u": 10, "s": "hello", "f": 3.14, "struct": {"i": 2, "s": "world", "f": 1.23}, "slice": [1, 2, 3, 4, 5], "map": {"a": 1, "b": 1.23, "c": "json"}, "time": "0001-01-01T00:00:00Z", "duration": "5m0s"}
`
	actual := "\n" + buf.String()
	if expect != actual {
		t.Fatalf("JSON style marshal error: expect=[%s] actual=[%s]", expect, actual)
	}
}

func TestEncoder_MarshalAnchor(t *testing.T) {
	type Host struct {
		Hostname string
		Username string
		Password string
	}
	type HostDecl struct {
		Host *Host `yaml:",anchor"`
	}
	type Queue struct {
		Name  string `yaml:","`
		*Host `yaml:",alias"`
	}
	var doc struct {
		Hosts  []*HostDecl `yaml:"hosts"`
		Queues []*Queue    `yaml:"queues"`
	}
	host1 := &Host{
		Hostname: "host1.example.com",
		Username: "userA",
		Password: "pass1",
	}
	host2 := &Host{
		Hostname: "host2.example.com",
		Username: "userB",
		Password: "pass2",
	}
	doc.Hosts = []*HostDecl{
		{
			Host: host1,
		},
		{
			Host: host2,
		},
	}
	doc.Queues = []*Queue{
		{
			Name: "queue",
			Host: host1,
		}, {
			Name: "queue2",
			Host: host2,
		},
	}
	hostIdx := 1
	opt := yaml.MarshalAnchor(func(anchor *ast.AnchorNode, value interface{}) error {
		if _, ok := value.(*Host); ok {
			nameNode := anchor.Name.(*ast.StringNode)
			nameNode.Value = fmt.Sprintf("host%d", hostIdx)
			hostIdx++
		}
		return nil
	})

	var buf bytes.Buffer
	if err := yaml.NewEncoder(&buf, opt).Encode(doc); err != nil {
		t.Fatalf("%+v", err)
	}
	expect := `
hosts:
- host: &host1
    hostname: host1.example.com
    username: userA
    password: pass1
- host: &host2
    hostname: host2.example.com
    username: userB
    password: pass2
queues:
- name: queue
  host: *host1
- name: queue2
  host: *host2
`
	if "\n"+buf.String() != expect {
		t.Fatalf("unexpected output. %s", buf.String())
	}
}

type useJSONMarshalerTest struct{}

func (t useJSONMarshalerTest) MarshalJSON() ([]byte, error) {
	return []byte(`{"a":[1, 2, 3]}`), nil
}

func TestEncoder_UseJSONMarshaler(t *testing.T) {
	got, err := yaml.MarshalWithOptions(useJSONMarshalerTest{}, yaml.UseJSONMarshaler())
	if err != nil {
		t.Fatal(err)
	}
	expected := `
a:
- 1
- 2
- 3
`
	if expected != "\n"+string(got) {
		t.Fatalf("failed to use json marshaler. expected [%q] but got [%q]", expected, string(got))
	}
}

func Example_Marshal_Node() {
	type T struct {
		Text ast.Node `yaml:"text"`
	}
	stringNode, err := yaml.ValueToNode("node example")
	if err != nil {
		panic(err)
	}
	bytes, err := yaml.Marshal(T{Text: stringNode})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
	// OUTPUT:
	// text: node example
}

func Example_Marshal_ExplicitAnchorAlias() {
	type T struct {
		A int
		B string
	}
	var v struct {
		C *T `yaml:"c,anchor=x"`
		D *T `yaml:"d,alias=x"`
	}
	v.C = &T{A: 1, B: "hello"}
	v.D = v.C
	bytes, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
	// OUTPUT:
	// c: &x
	//   a: 1
	//   b: hello
	// d: *x
}

func Example_Marshal_ImplicitAnchorAlias() {
	type T struct {
		I int
		S string
	}
	var v struct {
		A *T `yaml:"a,anchor"`
		B *T `yaml:"b,anchor"`
		C *T `yaml:"c,alias"`
		D *T `yaml:"d,alias"`
	}
	v.A = &T{I: 1, S: "hello"}
	v.B = &T{I: 2, S: "world"}
	v.C = v.A // C has same pointer address to A
	v.D = v.B // D has same pointer address to B
	bytes, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
	// OUTPUT:
	// a: &a
	//   i: 1
	//   s: hello
	// b: &b
	//   i: 2
	//   s: world
	// c: *a
	// d: *b
}

type tMarshal []string

func (t *tMarshal) MarshalYAML() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("tags:\n")
	for i, v := range *t {
		if i == 0 {
			fmt.Fprintf(&buf, "- %s\n", v)
		} else {
			fmt.Fprintf(&buf, "  %s\n", v)
		}
	}
	return buf.Bytes(), nil
}
func Test_Marshaler(t *testing.T) {
	const expected = `- hello-world
`

	// sanity check
	var l []string
	if err := yaml.Unmarshal([]byte(expected), &l); err != nil {
		t.Fatalf("failed to parse string: %s", err)
	}

	buf, err := yaml.Marshal(tMarshal{"hello-world"})
	if err != nil {
		t.Fatalf("failed to marshal: %s", err)
	}

	if string(buf) != expected {
		t.Fatalf("expected '%s', got '%s'", expected, buf)
	}

	t.Logf("%s", buf)
}

type marshalContext struct{}

func (c *marshalContext) MarshalYAML(ctx context.Context) ([]byte, error) {
	v, ok := ctx.Value("k").(int)
	if !ok {
		return nil, fmt.Errorf("cannot get valid context")
	}
	if v != 1 {
		return nil, fmt.Errorf("cannot get valid context")
	}
	return []byte("1"), nil
}

func Test_MarshalerContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "k", 1)
	bytes, err := yaml.MarshalContext(ctx, &marshalContext{})
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if string(bytes) != "1\n" {
		t.Fatalf("failed marshal: %q", string(bytes))
	}
}

type SlowMarshaler struct {
	A string
	B int
}
type FastMarshaler struct {
	A string
	B int
}
type TextMarshaler int64
type TextMarshalerContainer struct {
	Field TextMarshaler `yaml:"field"`
}

func (v SlowMarshaler) MarshalYAML() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("tags:\n")
	buf.WriteString("- slow-marshaler\n")
	buf.WriteString("a: " + v.A + "\n")
	buf.WriteString("b: " + strconv.FormatInt(int64(v.B), 10) + "\n")
	return buf.Bytes(), nil
}

func (v FastMarshaler) MarshalYAML() (interface{}, error) {
	return yaml.MapSlice{
		{"tags", []string{"fast-marshaler"}},
		{"a", v.A},
		{"b", v.B},
	}, nil
}

func (t TextMarshaler) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(t), 8)), nil
}

func Example_MarshalYAML() {
	var slow SlowMarshaler
	slow.A = "Hello slow poke"
	slow.B = 100
	buf, err := yaml.Marshal(slow)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(string(buf))

	var fast FastMarshaler
	fast.A = "Hello speed demon"
	fast.B = 100
	buf, err = yaml.Marshal(fast)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(string(buf))

	text := TextMarshalerContainer{
		Field: 11,
	}
	buf, err = yaml.Marshal(text)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(string(buf))
	// OUTPUT:
	// tags:
	// - slow-marshaler
	// a: Hello slow poke
	// b: 100
	//
	// tags:
	// - fast-marshaler
	// a: Hello speed demon
	// b: 100
	//
	// field: 13
}

func TestMarshalIndentWithMultipleText(t *testing.T) {
	t.Run("depth1", func(t *testing.T) {
		b, err := yaml.MarshalWithOptions(map[string]interface{}{
			"key": []string{`line1
line2
line3`},
		}, yaml.Indent(2))
		if err != nil {
			t.Fatal(err)
		}
		got := string(b)
		expected := `key:
- |-
  line1
  line2
  line3
`
		if expected != got {
			t.Fatalf("failed to encode.\nexpected:\n%s\nbut got:\n%s\n", expected, got)
		}
	})
	t.Run("depth2", func(t *testing.T) {
		b, err := yaml.MarshalWithOptions(map[string]interface{}{
			"key": map[string]interface{}{
				"key2": []string{`line1
line2
line3`},
			},
		}, yaml.Indent(2))
		if err != nil {
			t.Fatal(err)
		}
		got := string(b)
		expected := `key:
  key2:
  - |-
    line1
    line2
    line3
`
		if expected != got {
			t.Fatalf("failed to encode.\nexpected:\n%s\nbut got:\n%s\n", expected, got)
		}
	})
}

type bytesMarshaler struct{}

func (b *bytesMarshaler) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(map[string]interface{}{"d": "foo"})
}

func TestBytesMarshaler(t *testing.T) {
	b, err := yaml.Marshal(map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": &bytesMarshaler{},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := `
a:
  b:
    c:
      d: foo
`
	got := "\n" + string(b)
	if expected != got {
		t.Fatalf("failed to encode. expected %s but got %s", expected, got)
	}
}
