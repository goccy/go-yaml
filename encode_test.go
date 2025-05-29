package yaml_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"net/netip"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

var zero = 0
var emptyStr = ""

type TestTextMarshaler string

func (t TestTextMarshaler) MarshalText() ([]byte, error) {
	return []byte(t), nil
}

type TestTextUnmarshalerContainer struct {
	V TestTextMarshaler
}

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
			map[string]int64{"v": int64(4294967296)},
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
			"v: 1e-06\n",
			map[string]float32{"v": 1e-06},
			nil,
		},
		{
			"v: 1e-06\n",
			map[string]float64{"v": 0.000001},
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
			"v: 1e-06\n",
			map[string]float64{"v": 0.000001},
			nil,
		},
		{
			"v: 1e-06\n",
			map[string]float64{"v": 1e-06},
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
			"a: \"-\"\n",
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
			"t2: \"2018-01-09T10:40:47Z\"\nt4: \"2098-01-09T10:40:47Z\"\n",
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
			"a: \" b\"\n",
			map[string]string{"a": " b"},
			nil,
		},
		{
			"a: \"b \"\n",
			map[string]string{"a": "b "},
			nil,
		},
		{
			"a: \" b \"\n",
			map[string]string{"a": " b "},
			nil,
		},
		{
			"a: \"`b` c\"\n",
			map[string]string{"a": "`b` c"},
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

		// Omitempty flag.
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
			"a:\n  \"y\": \"\"\n",
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
		// Highlighting differences of go-yaml omitempty vs std encoding/json
		// omitempty. Encoding/json will emit the following fields: https://go.dev/play/p/VvNpdM0GD4d
		{
			"{}\n",
			struct {
				// This type has a custom IsZero method.
				A netip.Addr         `yaml:"a,omitempty"`
				B struct{ X, y int } `yaml:"b,omitempty"`
			}{},
			nil,
		},

		// omitzero flag.
		{
			"a: 1\n",
			struct {
				A int `yaml:"a,omitzero"`
				B int `yaml:"b,omitzero"`
			}{1, 0},
			nil,
		},
		{
			"{}\n",
			struct {
				A int `yaml:"a,omitzero"`
				B int `yaml:"b,omitzero"`
			}{0, 0},
			nil,
		},
		{
			"a:\n  \"y\": \"\"\n",
			struct {
				A *struct {
					X string `yaml:"x,omitzero"`
					Y string
				}
			}{&struct {
				X string `yaml:"x,omitzero"`
				Y string
			}{}},
			nil,
		},
		{
			"a: {}\n",
			struct {
				A *struct {
					X string `yaml:"x,omitzero"`
					Y string `yaml:"y,omitzero"`
				}
			}{&struct {
				X string `yaml:"x,omitzero"`
				Y string `yaml:"y,omitzero"`
			}{}},
			nil,
		},
		{
			"a: {x: 1}\n",
			struct {
				A *struct{ X, y int } `yaml:"a,omitzero,flow"`
			}{&struct{ X, y int }{1, 2}},
			nil,
		},
		{
			"{}\n",
			struct {
				A *struct{ X, y int } `yaml:"a,omitzero,flow"`
			}{nil},
			nil,
		},
		{
			"a: {x: 0}\n",
			struct {
				A *struct{ X, y int } `yaml:"a,omitzero,flow"`
			}{&struct{ X, y int }{}},
			nil,
		},
		{
			"a: {x: 1}\n",
			struct {
				A struct{ X, y int } `yaml:"a,omitzero,flow"`
			}{struct{ X, y int }{1, 2}},
			nil,
		},
		{
			"{}\n",
			struct {
				A struct{ X, y int } `yaml:"a,omitzero,flow"`
			}{struct{ X, y int }{0, 1}},
			nil,
		},
		{
			"a: 1.0\n",
			struct {
				A float64 `yaml:"a,omitzero"`
				B float64 `yaml:"b,omitzero"`
			}{1, 0},
			nil,
		},
		{
			"a: 1\nb: []\n",
			struct {
				A int
				B []string `yaml:"b,omitzero"`
			}{
				1, []string{},
			},
			nil,
		},
		{
			"{}\n",
			struct {
				A netip.Addr         `yaml:"a,omitzero"`
				B struct{ X, y int } `yaml:"b,omitzero"`
			}{},
			nil,
		},

		// OmitEmpty global option.
		{
			"a: 1\n",
			struct {
				A int
				B int `yaml:"b,omitempty"`
			}{1, 0},
			[]yaml.EncodeOption{
				yaml.OmitEmpty(),
			},
		},
		{
			"{}\n",
			struct {
				A int
				B int `yaml:"b,omitempty"`
			}{0, 0},
			[]yaml.EncodeOption{
				yaml.OmitEmpty(),
			},
		},
		{
			"a: \"\"\nb: {}\n",
			struct {
				A netip.Addr         `yaml:"a"`
				B struct{ X, y int } `yaml:"b"`
			}{},
			[]yaml.EncodeOption{
				yaml.OmitEmpty(),
			},
		},

		// OmitZero global option.
		{
			"a: 1\n",
			struct {
				A int
				B int
			}{1, 0},
			[]yaml.EncodeOption{
				yaml.OmitZero(),
			},
		},
		{
			"{}\n",
			struct {
				A int
				B int
			}{0, 0},
			[]yaml.EncodeOption{
				yaml.OmitZero(),
			},
		},
		{
			"{}\n",
			struct {
				A netip.Addr         `yaml:"a"`
				B struct{ X, y int } `yaml:"b"`
			}{},
			[]yaml.EncodeOption{
				yaml.OmitZero(),
			},
		},

		// Flow flag.
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
		// Quoting in flow mode
		{
			`a: [b, "c,d", e]` + "\n",
			struct {
				A []string `yaml:"a,flow"`
			}{[]string{"b", "c,d", "e"}},
			[]yaml.EncodeOption{
				yaml.UseSingleQuote(false),
			},
		},
		{
			`a: [b, "c]", d]` + "\n",
			struct {
				A []string `yaml:"a,flow"`
			}{[]string{"b", "c]", "d"}},
			[]yaml.EncodeOption{
				yaml.UseSingleQuote(false),
			},
		},
		{
			`a: [b, "c}", d]` + "\n",
			struct {
				A []string `yaml:"a,flow"`
			}{[]string{"b", "c}", "d"}},
			[]yaml.EncodeOption{
				yaml.UseSingleQuote(false),
			},
		},
		{
			`a: [b, "c\"", d]` + "\n",
			struct {
				A []string `yaml:"a,flow"`
			}{[]string{"b", `c"`, "d"}},
			[]yaml.EncodeOption{
				yaml.UseSingleQuote(false),
			},
		},
		{
			`a: [b, "c'", d]` + "\n",
			struct {
				A []string `yaml:"a,flow"`
			}{[]string{"b", "c'", "d"}},
			[]yaml.EncodeOption{
				yaml.UseSingleQuote(false),
			},
		},
		// No quoting in non-flow mode
		{
			"a:\n- b\n- c,d\n- e\n",
			struct {
				A []string `yaml:"a"`
			}{[]string{"b", "c,d", "e"}},
			nil,
		},
		{
			`a: [b, "c]", d]` + "\n",
			struct {
				A []string `yaml:"a,flow"`
			}{[]string{"b", "c]", "d"}},
			nil,
		},
		{
			`a: [b, "c}", d]` + "\n",
			struct {
				A []string `yaml:"a,flow"`
			}{[]string{"b", "c}", "d"}},
			nil,
		},
		{
			`a: [b, "c\"", d]` + "\n",
			struct {
				A []string `yaml:"a,flow"`
			}{[]string{"b", `c"`, "d"}},
			nil,
		},
		{
			`a: [b, "c'", d]` + "\n",
			struct {
				A []string `yaml:"a,flow"`
			}{[]string{"b", "c'", "d"}},
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
			map[string]time.Time{"v": {}},
			nil,
		},
		{
			"v: 0001-01-01T00:00:00Z\n",
			map[string]*time.Time{"v": {}},
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
		{
			"v: 30s\n",
			map[string]*time.Duration{"v": ptr(30 * time.Second)},
			nil,
		},
		{
			"v: null\n",
			map[string]*time.Duration{"v": nil},
			nil,
		},
		{
			"v: test\n",
			TestTextUnmarshalerContainer{V: "test"},
			nil,
		},
		{
			"v: \"1\"\n",
			TestTextUnmarshalerContainer{V: "1"},
			nil,
		},
		{
			"v: \"#\"\n",
			TestTextUnmarshalerContainer{V: "#"},
			nil,
		},
		// Quote style
		{
			`v: '''a''b'` + "\n",
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
		{
			`a: '\.yaml'` + "\n",
			map[string]string{"a": `\.yaml`},
			[]yaml.EncodeOption{
				yaml.UseSingleQuote(true),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			var buf bytes.Buffer
			enc := yaml.NewEncoder(&buf, test.options...)
			if err := enc.Encode(test.value); err != nil {
				t.Fatalf("%+v", err)
			}
			if test.source != buf.String() {
				t.Fatalf("expect = [%s], actual = [%s]", test.source, buf.String())
			}
		})
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
	expect := "a:\n  m:\n    x: \"y\"\n"
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
	expect := "m:\n  x: \"y\"\n"
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
		C *T `yaml:"c"`
		D *T `yaml:"d"`
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
		C *T `yaml:"c"`
		D *T `yaml:"d"`
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
		*Person `yaml:",omitempty,inline"`
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

func TestEncoder_InlineNil(t *testing.T) {
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
		C: true,
	}); err != nil {
		t.Fatalf("%+v", err)
	}
	expect := `
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
		"test": {1, 2, 3},
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
		Name string `yaml:","`
		*Host
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
			nameNode, _ := anchor.Name.(*ast.StringNode)
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

func TestEncoder_CustomMarshaler(t *testing.T) {
	t.Run("override struct type", func(t *testing.T) {
		type T struct {
			Foo string `yaml:"foo"`
		}
		b, err := yaml.MarshalWithOptions(&T{Foo: "bar"}, yaml.CustomMarshaler[T](func(v T) ([]byte, error) {
			return []byte(`"override"`), nil
		}))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(b, []byte("\"override\"\n")) {
			t.Fatalf("failed to switch to custom marshaler. got: %q", b)
		}
	})
	t.Run("override bytes type", func(t *testing.T) {
		type T struct {
			Foo []byte `yaml:"foo"`
		}
		b, err := yaml.MarshalWithOptions(&T{Foo: []byte("bar")}, yaml.CustomMarshaler[[]byte](func(v []byte) ([]byte, error) {
			if !bytes.Equal(v, []byte("bar")) {
				t.Fatalf("failed to get src buffer: %q", v)
			}
			return []byte(`override`), nil
		}))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(b, []byte("foo: override\n")) {
			t.Fatalf("failed to switch to custom marshaler. got: %q", b)
		}
	})
	t.Run("override bytes type with context", func(t *testing.T) {
		type T struct {
			Foo []byte `yaml:"foo"`
		}
		ctx := context.WithValue(context.Background(), "plop", uint(42))
		b, err := yaml.MarshalContext(ctx, &T{Foo: []byte("bar")}, yaml.CustomMarshalerContext[[]byte](func(ctx context.Context, v []byte) ([]byte, error) {
			if !bytes.Equal(v, []byte("bar")) {
				t.Fatalf("failed to get src buffer: %q", v)
			}
			if ctx.Value("plop") != uint(42) {
				t.Fatalf("context value is not correct")
			}
			return []byte(`override`), nil
		}))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(b, []byte("foo: override\n")) {
			t.Fatalf("failed to switch to custom marshaler. got: %q", b)
		}
	})
}

func TestEncoder_AutoInt(t *testing.T) {
	for _, test := range []struct {
		desc     string
		input    any
		expected string
	}{
		{
			desc: "int-convertible float64",
			input: map[string]float64{
				"key": 1.0,
			},
			expected: "key: 1\n",
		},
		{
			desc: "non int-convertible float64",
			input: map[string]float64{
				"key": 1.1,
			},
			expected: "key: 1.1\n",
		},
		{
			desc: "int-convertible float32",
			input: map[string]float32{
				"key": 1.0,
			},
			expected: "key: 1\n",
		},
		{
			desc: "non int-convertible float32",
			input: map[string]float32{
				"key": 1.1,
			},
			expected: "key: 1.1\n",
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			var buf bytes.Buffer
			enc := yaml.NewEncoder(&buf, yaml.AutoInt())
			if err := enc.Encode(test.input); err != nil {
				t.Fatalf("failed to encode: %s", err)
			}
			if actual := buf.String(); actual != test.expected {
				t.Errorf("expect:\n%s\nactual\n%s\n", test.expected, actual)
			}
		})
	}
}

func TestEncoder_MultipleDocuments(t *testing.T) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	if err := enc.Encode(1); err != nil {
		t.Fatalf("failed to encode: %s", err)
	}
	if err := enc.Encode(2); err != nil {
		t.Fatalf("failed to encode: %s", err)
	}
	if actual, expect := buf.String(), "1\n---\n2\n"; actual != expect {
		t.Errorf("expect:\n%s\nactual\n%s\n", expect, actual)
	}
}

func TestEncoder_UnmarshallableTypes(t *testing.T) {
	for _, test := range []struct {
		desc        string
		input       any
		expectedErr string
	}{
		{
			desc:        "channel",
			input:       make(chan int),
			expectedErr: "unknown value type chan int",
		},
		{
			desc:        "function",
			input:       func() {},
			expectedErr: "unknown value type func()",
		},
		{
			desc:        "complex number",
			input:       complex(10, 11),
			expectedErr: "unknown value type complex128",
		},
		{
			desc:        "unsafe pointer",
			input:       unsafe.Pointer(&struct{}{}),
			expectedErr: "unknown value type unsafe.Pointer",
		},
		{
			desc:        "uintptr",
			input:       uintptr(0x1234),
			expectedErr: "unknown value type uintptr",
		},
		{
			desc:        "map with channel",
			input:       map[string]any{"key": make(chan string)},
			expectedErr: "unknown value type chan string",
		},
		{
			desc: "nested map with func",
			input: map[string]any{
				"a": map[string]any{
					"b": func(_ string) {},
				},
			},
			expectedErr: "unknown value type func(string)",
		},
		{
			desc:        "slice with channel",
			input:       []any{make(chan bool)},
			expectedErr: "unknown value type chan bool",
		},
		{
			desc:        "nested slice with complex number",
			input:       []any{[]any{complex(10, 11)}},
			expectedErr: "unknown value type complex128",
		},
		{
			desc: "struct with unsafe pointer",
			input: struct {
				Field unsafe.Pointer `yaml:"field"`
			}{},
			expectedErr: "unknown value type unsafe.Pointer",
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			var buf bytes.Buffer
			err := yaml.NewEncoder(&buf).Encode(test.input)
			if err == nil {
				t.Errorf("expect error:\n%s\nbut got none\n", test.expectedErr)
			} else if err.Error() != test.expectedErr {
				t.Errorf("expect error:\n%s\nactual\n%s\n", test.expectedErr, err)
			}
		})
	}
}

func ExampleMarshal_node() {
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

func ExampleMarshal_explicitAnchorAlias() {
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

func ExampleMarshal_implicitAnchorAlias() {
	type T struct {
		I int
		S string
	}
	var v struct {
		A *T `yaml:"a,anchor"`
		B *T `yaml:"b,anchor"`
		C *T `yaml:"c"`
		D *T `yaml:"d"`
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
		return nil, errors.New("cannot get valid context")
	}
	if v != 1 {
		return nil, errors.New("cannot get valid context")
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

func ExampleMarshal() {
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
	// field: "13"
}

func TestIssue356(t *testing.T) {
	tests := map[string]struct {
		in string
	}{
		"content on first line": {
			in: `args:
  - |

    key:
      nest1: something
      nest2:
        nest2a: b
`,
		},
		"empty first line": {
			in: `args:
  - |

    key:
      nest1: something
      nest2:
        nest2a: b
`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f, err := parser.ParseBytes([]byte(test.in), 0)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			got := f.String()
			if test.in != got {
				t.Fatalf("failed to encode.\nexpected:\n%s\nbut got:\n%s\n", test.in, got)
			}
		})
	}
}

func TestMarshalIndentWithMultipleText(t *testing.T) {
	tests := []struct {
		name   string
		input  map[string]interface{}
		indent yaml.EncodeOption
		want   string
	}{
		{
			name: "depth1",
			input: map[string]interface{}{
				"key": []string{`line1
line2
line3`},
			},
			indent: yaml.Indent(2),
			want: `key:
- |-
  line1
  line2
  line3
`,
		},
		{
			name: "depth2",
			input: map[string]interface{}{
				"key": map[string]interface{}{
					"key2": []string{`line1
line2
line3`},
				},
			},
			indent: yaml.Indent(2),
			want: `key:
  key2:
  - |-
    line1
    line2
    line3
`,
		},
		{
			name: "raw string new lines",
			input: map[string]interface{}{
				"key": "line1\nline2\nline3",
			},
			indent: yaml.Indent(4),
			want: `key: |-
    line1
    line2
    line3
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := yaml.MarshalWithOptions(tt.input, tt.indent)
			if err != nil {
				t.Fatalf("failed to marshal yaml: %v", err)
			}
			got := string(b)
			if tt.want != got {
				t.Fatalf("failed to encode.\nexpected:\n%s\nbut got:\n%s\n", tt.want, got)
			}
		})
	}
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

type customMapSliceOneItemMarshaler struct{}

func (m *customMapSliceOneItemMarshaler) MarshalYAML() ([]byte, error) {
	var v yaml.MapSlice
	v = append(v, yaml.MapItem{"a", "b"})
	return yaml.Marshal(v)
}

type customMapSliceTwoItemMarshaler struct{}

func (m *customMapSliceTwoItemMarshaler) MarshalYAML() ([]byte, error) {
	var v yaml.MapSlice
	v = append(v, yaml.MapItem{"a", "b"})
	v = append(v, yaml.MapItem{"b", "c"})
	return yaml.Marshal(v)
}

func TestCustomMapSliceMarshaler(t *testing.T) {
	type T struct {
		A *customMapSliceOneItemMarshaler `yaml:"a"`
		B *customMapSliceTwoItemMarshaler `yaml:"b"`
	}
	b, err := yaml.Marshal(&T{
		A: &customMapSliceOneItemMarshaler{},
		B: &customMapSliceTwoItemMarshaler{},
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := `
a:
  a: b
b:
  a: b
  b: c
`
	got := "\n" + string(b)
	if expected != got {
		t.Fatalf("failed to encode. expected %s but got %s", expected, got)
	}
}

type Issue174 struct {
	K string
	V []int
}

func (v Issue174) MarshalYAML() ([]byte, error) {
	return yaml.MarshalWithOptions(map[string][]int{v.K: v.V}, yaml.Flow(true))
}

func TestIssue174(t *testing.T) {
	b, err := yaml.Marshal(Issue174{"00:00:00-23:59:59", []int{1, 2, 3}})
	if err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSuffix(string(b), "\n")
	if got != `{"00:00:00-23:59:59": [1, 2, 3]}` {
		t.Fatalf("failed to encode: %q", got)
	}
}

func TestIssue259(t *testing.T) {
	type AnchorValue struct {
		Foo uint64
		Bar string
	}

	type Value struct {
		Baz   string       `yaml:"baz"`
		Value *AnchorValue `yaml:"value,anchor"`
	}

	type Schema struct {
		Values []*Value
	}

	schema := Schema{}
	anchorValue := AnchorValue{Foo: 3, Bar: "bar"}
	schema.Values = []*Value{
		{Baz: "xxx", Value: &anchorValue},
		{Baz: "yyy", Value: &anchorValue},
		{Baz: "zzz", Value: &anchorValue},
	}
	b, err := yaml.Marshal(schema)
	if err != nil {
		t.Fatal(err)
	}
	expected := `
values:
- baz: xxx
  value: &value
    foo: 3
    bar: bar
- baz: yyy
  value: *value
- baz: zzz
  value: *value
`
	if strings.TrimPrefix(expected, "\n") != string(b) {
		t.Fatalf("failed to encode: got = %s", string(b))
	}
}

func TestTagMarshalling(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "scalar", input: "a: !mytag 1"},
		{name: "mapping", input: `
a: !mytag
  b: 2`},
		{name: "sequence", input: `
a: !mytag
- 1
- 2
- 3`},
		{name: "anchor before tag", input: `
a: &anc !mytag
- 1
- 2
- 3`},
		{name: "flow mapping", input: "a: !mytag {b: 2}"},
		{name: "flow sequence", input: "a: !mytag [1, 2, 3]"},
		{name: "explicit type", input: "a: !!timestamp test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := parser.ParseBytes([]byte(tt.input), 0)
			result, err := yaml.Marshal(res.Docs[0])
			if err != nil {
				t.Fatal(err)
			}

			expected := strings.TrimSpace(tt.input)
			output := strings.TrimSpace(string(result))
			if expected != output {
				t.Fatalf("input is not equal to output.\n\nexpected:\n%v\n actual:\n%v", expected, output)
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
