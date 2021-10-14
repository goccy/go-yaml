package yaml_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/internal/errors"
	"github.com/goccy/go-yaml/parser"
	"golang.org/x/xerrors"
)

type Child struct {
	B int
	C int `yaml:"-"`
}

func TestDecoder(t *testing.T) {
	tests := []struct {
		source string
		value  interface{}
	}{
		{
			"null\n",
			(*struct{})(nil),
		},
		{
			"v: hi\n",
			map[string]string{"v": "hi"},
		},
		{
			"v: \"true\"\n",
			map[string]string{"v": "true"},
		},
		{
			"v: \"false\"\n",
			map[string]string{"v": "false"},
		},
		{
			"v: true\n",
			map[string]interface{}{"v": true},
		},
		{
			"v: true\n",
			map[string]string{"v": "true"},
		},
		{
			"v: 10\n",
			map[string]string{"v": "10"},
		},
		{
			"v: -10\n",
			map[string]string{"v": "-10"},
		},
		{
			"v: 1.234\n",
			map[string]string{"v": "1.234"},
		},
		{
			"v: false\n",
			map[string]bool{"v": false},
		},
		{
			"v: 10\n",
			map[string]int{"v": 10},
		},
		{
			"v: 10",
			map[string]interface{}{"v": 10},
		},
		{
			"v: 0b10",
			map[string]interface{}{"v": 2},
		},
		{
			"v: -0b101010",
			map[string]interface{}{"v": -42},
		},
		{
			"v: -0b1000000000000000000000000000000000000000000000000000000000000000",
			map[string]interface{}{"v": -9223372036854775808},
		},
		{
			"v: 0xA",
			map[string]interface{}{"v": 10},
		},
		{
			"v: .1",
			map[string]interface{}{"v": 0.1},
		},
		{
			"v: -.1",
			map[string]interface{}{"v": -0.1},
		},
		{
			"v: -10\n",
			map[string]int{"v": -10},
		},
		{
			"v: 4294967296\n",
			map[string]int{"v": 4294967296},
		},
		{
			"v: 0.1\n",
			map[string]interface{}{"v": 0.1},
		},
		{
			"v: 0.99\n",
			map[string]float32{"v": 0.99},
		},
		{
			"v: -0.1\n",
			map[string]float64{"v": -0.1},
		},
		{
			"v: 6.8523e+5",
			map[string]interface{}{"v": 6.8523e+5},
		},
		{
			"v: 685.230_15e+03",
			map[string]interface{}{"v": 685.23015e+03},
		},
		{
			"v: 685_230.15",
			map[string]interface{}{"v": 685230.15},
		},
		{
			"v: 685_230.15",
			map[string]float64{"v": 685230.15},
		},
		{
			"v: 685230",
			map[string]interface{}{"v": 685230},
		},
		{
			"v: +685_230",
			map[string]interface{}{"v": 685230},
		},
		{
			"v: 02472256",
			map[string]interface{}{"v": 685230},
		},
		{
			"v: 0x_0A_74_AE",
			map[string]interface{}{"v": 685230},
		},
		{
			"v: 0b1010_0111_0100_1010_1110",
			map[string]interface{}{"v": 685230},
		},
		{
			"v: +685_230",
			map[string]int{"v": 685230},
		},

		// Bools from spec
		{
			"v: True",
			map[string]interface{}{"v": true},
		},
		{
			"v: TRUE",
			map[string]interface{}{"v": true},
		},
		{
			"v: False",
			map[string]interface{}{"v": false},
		},
		{
			"v: FALSE",
			map[string]interface{}{"v": false},
		},
		{
			"v: y",
			map[string]interface{}{"v": "y"}, // y or yes or Yes is string
		},
		{
			"v: NO",
			map[string]interface{}{"v": "NO"}, // no or No or NO is string
		},
		{
			"v: on",
			map[string]interface{}{"v": "on"}, // on is string
		},

		// Some cross type conversions
		{
			"v: 42",
			map[string]uint{"v": 42},
		}, {
			"v: 4294967296",
			map[string]uint64{"v": 4294967296},
		},

		// int
		{
			"v: 2147483647",
			map[string]int{"v": math.MaxInt32},
		},
		{
			"v: -2147483648",
			map[string]int{"v": math.MinInt32},
		},

		// int64
		{
			"v: 9223372036854775807",
			map[string]int64{"v": math.MaxInt64},
		},
		{
			"v: 0b111111111111111111111111111111111111111111111111111111111111111",
			map[string]int64{"v": math.MaxInt64},
		},
		{
			"v: -9223372036854775808",
			map[string]int64{"v": math.MinInt64},
		},
		{
			"v: -0b111111111111111111111111111111111111111111111111111111111111111",
			map[string]int64{"v": -math.MaxInt64},
		},

		// uint
		{
			"v: 0",
			map[string]uint{"v": 0},
		},
		{
			"v: 4294967295",
			map[string]uint{"v": math.MaxUint32},
		},

		// uint64
		{
			"v: 0",
			map[string]uint{"v": 0},
		},
		{
			"v: 18446744073709551615",
			map[string]uint64{"v": math.MaxUint64},
		},
		{
			"v: 0b1111111111111111111111111111111111111111111111111111111111111111",
			map[string]uint64{"v": math.MaxUint64},
		},
		{
			"v: 9223372036854775807",
			map[string]uint64{"v": math.MaxInt64},
		},

		// float32
		{
			"v: 3.40282346638528859811704183484516925440e+38",
			map[string]float32{"v": math.MaxFloat32},
		},
		{
			"v: 1.401298464324817070923729583289916131280e-45",
			map[string]float32{"v": math.SmallestNonzeroFloat32},
		},
		{
			"v: 18446744073709551615",
			map[string]float32{"v": float32(math.MaxUint64)},
		},
		{
			"v: 18446744073709551616",
			map[string]float32{"v": float32(math.MaxUint64 + 1)},
		},

		// float64
		{
			"v: 1.797693134862315708145274237317043567981e+308",
			map[string]float64{"v": math.MaxFloat64},
		},
		{
			"v: 4.940656458412465441765687928682213723651e-324",
			map[string]float64{"v": math.SmallestNonzeroFloat64},
		},
		{
			"v: 18446744073709551615",
			map[string]float64{"v": float64(math.MaxUint64)},
		},
		{
			"v: 18446744073709551616",
			map[string]float64{"v": float64(math.MaxUint64 + 1)},
		},

		// Timestamps
		{
			// Date only.
			"v: 2015-01-01\n",
			map[string]time.Time{"v": time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			// RFC3339
			"v: 2015-02-24T18:19:39.12Z\n",
			map[string]time.Time{"v": time.Date(2015, 2, 24, 18, 19, 39, .12e9, time.UTC)},
		},
		{
			// RFC3339 with short dates.
			"v: 2015-2-3T3:4:5Z",
			map[string]time.Time{"v": time.Date(2015, 2, 3, 3, 4, 5, 0, time.UTC)},
		},
		{
			// ISO8601 lower case t
			"v: 2015-02-24t18:19:39Z\n",
			map[string]time.Time{"v": time.Date(2015, 2, 24, 18, 19, 39, 0, time.UTC)},
		},
		{
			// space separate, no time zone
			"v: 2015-02-24 18:19:39\n",
			map[string]time.Time{"v": time.Date(2015, 2, 24, 18, 19, 39, 0, time.UTC)},
		},
		{
			"v: 60s\n",
			map[string]time.Duration{"v": time.Minute},
		},
		{
			"v: -0.5h\n",
			map[string]time.Duration{"v": -30 * time.Minute},
		},

		// Single Quoted values.
		{
			`'1': '2'`,
			map[interface{}]interface{}{"1": `2`},
		},
		{
			`'1': '"2"'`,
			map[interface{}]interface{}{"1": `"2"`},
		},
		{
			`'1': ''''`,
			map[interface{}]interface{}{"1": `'`},
		},
		{
			`'1': '''2'''`,
			map[interface{}]interface{}{"1": `'2'`},
		},
		{
			`'1': 'B''z'`,
			map[interface{}]interface{}{"1": `B'z`},
		},
		{
			`'1': '\'`,
			map[interface{}]interface{}{"1": `\`},
		},
		{
			`'1': '\\'`,
			map[interface{}]interface{}{"1": `\\`},
		},
		{
			`'1': '\"2\"'`,
			map[interface{}]interface{}{"1": `\"2\"`},
		},
		{
			`'1': '\\"2\\"'`,
			map[interface{}]interface{}{"1": `\\"2\\"`},
		},
		{
			"'1': '   1\n    2\n    3'",
			map[interface{}]interface{}{"1": "   1 2 3"},
		},
		{
			"'1': '\n    2\n    3'",
			map[interface{}]interface{}{"1": " 2 3"},
		},

		// Double Quoted values.
		{
			`"1": "2"`,
			map[interface{}]interface{}{"1": `2`},
		},
		{
			`"1": "\"2\""`,
			map[interface{}]interface{}{"1": `"2"`},
		},
		{
			`"1": "\""`,
			map[interface{}]interface{}{"1": `"`},
		},
		{
			`"1": "X\"z"`,
			map[interface{}]interface{}{"1": `X"z`},
		},
		{
			`"1": "\\"`,
			map[interface{}]interface{}{"1": `\`},
		},
		{
			`"1": "\\\\"`,
			map[interface{}]interface{}{"1": `\\`},
		},
		{
			`"1": "\\\"2\\\""`,
			map[interface{}]interface{}{"1": `\"2\"`},
		},
		{
			"'1': \"   1\n    2\n    3\"",
			map[interface{}]interface{}{"1": "   1 2 3"},
		},
		{
			"'1': \"\n    2\n    3\"",
			map[interface{}]interface{}{"1": " 2 3"},
		},
		{
			`"1": "a\x2Fb"`,
			map[interface{}]interface{}{"1": `a/b`},
		},
		{
			`"1": "a\u002Fb"`,
			map[interface{}]interface{}{"1": `a/b`},
		},
		{
			`"1": "a\x2Fb\u002Fc\U0000002Fd"`,
			map[interface{}]interface{}{"1": `a/b/c/d`},
		},

		{
			"a: -b_c",
			map[string]interface{}{"a": "-b_c"},
		},
		{
			"a: +b_c",
			map[string]interface{}{"a": "+b_c"},
		},
		{
			"a: 50cent_of_dollar",
			map[string]interface{}{"a": "50cent_of_dollar"},
		},

		// Nulls
		{
			"v:",
			map[string]interface{}{"v": nil},
		},
		{
			"v: ~",
			map[string]interface{}{"v": nil},
		},
		{
			"~: null key",
			map[interface{}]string{nil: "null key"},
		},
		{
			"v:",
			map[string]*bool{"v": nil},
		},
		{
			"v: null",
			map[string]*string{"v": nil},
		},
		{
			"v: null",
			map[string]string{"v": ""},
		},
		{
			"v: null",
			map[string]interface{}{"v": nil},
		},
		{
			"v: Null",
			map[string]interface{}{"v": nil},
		},
		{
			"v: NULL",
			map[string]interface{}{"v": nil},
		},
		{
			"v: ~",
			map[string]*string{"v": nil},
		},
		{
			"v: ~",
			map[string]string{"v": ""},
		},

		{
			"v: .inf\n",
			map[string]interface{}{"v": math.Inf(0)},
		},
		{
			"v: .Inf\n",
			map[string]interface{}{"v": math.Inf(0)},
		},
		{
			"v: .INF\n",
			map[string]interface{}{"v": math.Inf(0)},
		},
		{
			"v: -.inf\n",
			map[string]interface{}{"v": math.Inf(-1)},
		},
		{
			"v: -.Inf\n",
			map[string]interface{}{"v": math.Inf(-1)},
		},
		{
			"v: -.INF\n",
			map[string]interface{}{"v": math.Inf(-1)},
		},
		{
			"v: .nan\n",
			map[string]interface{}{"v": math.NaN()},
		},
		{
			"v: .NaN\n",
			map[string]interface{}{"v": math.NaN()},
		},
		{
			"v: .NAN\n",
			map[string]interface{}{"v": math.NaN()},
		},

		// Explicit tags.
		{
			"v: !!float '1.1'",
			map[string]interface{}{"v": 1.1},
		},
		{
			"v: !!float 0",
			map[string]interface{}{"v": float64(0)},
		},
		{
			"v: !!float -1",
			map[string]interface{}{"v": float64(-1)},
		},
		{
			"v: !!null ''",
			map[string]interface{}{"v": nil},
		},
		{
			"v: !!timestamp \"2015-01-01\"",
			map[string]time.Time{"v": time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			"v: !!timestamp 2015-01-01",
			map[string]time.Time{"v": time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
		},

		// Flow sequence
		{
			"v: [A,B]",
			map[string]interface{}{"v": []interface{}{"A", "B"}},
		},
		{
			"v: [A,B,C,]",
			map[string][]string{"v": []string{"A", "B", "C"}},
		},
		{
			"v: [A,1,C]",
			map[string][]string{"v": []string{"A", "1", "C"}},
		},
		{
			"v: [A,1,C]",
			map[string]interface{}{"v": []interface{}{"A", 1, "C"}},
		},

		// Block sequence
		{
			"v:\n - A\n - B",
			map[string]interface{}{"v": []interface{}{"A", "B"}},
		},
		{
			"v:\n - A\n - B\n - C",
			map[string][]string{"v": []string{"A", "B", "C"}},
		},
		{
			"v:\n - A\n - 1\n - C",
			map[string][]string{"v": []string{"A", "1", "C"}},
		},
		{
			"v:\n - A\n - 1\n - C",
			map[string]interface{}{"v": []interface{}{"A", 1, "C"}},
		},

		// Map inside interface with no type hints.
		{
			"a: {b: c}",
			map[interface{}]interface{}{"a": map[interface{}]interface{}{"b": "c"}},
		},

		{
			"v: \"\"\n",
			map[string]string{"v": ""},
		},
		{
			"v:\n- A\n- B\n",
			map[string][]string{"v": {"A", "B"}},
		},
		{
			"a: '-'\n",
			map[string]string{"a": "-"},
		},
		{
			"123\n",
			123,
		},
		{
			"hello: world\n",
			map[string]string{"hello": "world"},
		},
		{
			"hello: world\r\n",
			map[string]string{"hello": "world"},
		},
		{
			"hello: world\rGo: Gopher",
			map[string]string{"hello": "world", "Go": "Gopher"},
		},

		// Structs and type conversions.
		{
			"hello: world",
			struct{ Hello string }{"world"},
		},
		{
			"a: {b: c}",
			struct{ A struct{ B string } }{struct{ B string }{"c"}},
		},
		{
			"a: {b: c}",
			struct{ A map[string]string }{map[string]string{"b": "c"}},
		},
		{
			"a:",
			struct{ A map[string]string }{},
		},
		{
			"a: 1",
			struct{ A int }{1},
		},
		{
			"a: 1",
			struct{ A float64 }{1},
		},
		{
			"a: 1.0",
			struct{ A int }{1},
		},
		{
			"a: 1.0",
			struct{ A uint }{1},
		},
		{
			"a: [1, 2]",
			struct{ A []int }{[]int{1, 2}},
		},
		{
			"a: [1, 2]",
			struct{ A [2]int }{[2]int{1, 2}},
		},
		{
			"a: 1",
			struct{ B int }{0},
		},
		{
			"a: 1",
			struct {
				B int `yaml:"a"`
			}{1},
		},

		{
			"a: 1\n",
			yaml.MapItem{Key: "a", Value: 1},
		},
		{
			"a: 1\nb: 2\nc: 3\n",
			yaml.MapSlice{
				{Key: "a", Value: 1},
				{Key: "b", Value: 2},
				{Key: "c", Value: 3},
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
		},
		{
			"a:\n  b: c\n",
			map[string]interface{}{
				"a": map[string]string{
					"b": "c",
				},
			},
		},
		{
			"a: {x: 1}\n",
			map[string]map[string]int{
				"a": {
					"x": 1,
				},
			},
		},
		{
			"t2: 2018-01-09T10:40:47Z\nt4: 2098-01-09T10:40:47Z\n",
			map[string]string{
				"t2": "2018-01-09T10:40:47Z",
				"t4": "2098-01-09T10:40:47Z",
			},
		},
		{
			"a: [1, 2]\n",
			map[string][]int{
				"a": {1, 2},
			},
		},
		{
			"a: {b: c, d: e}\n",
			map[string]interface{}{
				"a": map[string]string{
					"b": "c",
					"d": "e",
				},
			},
		},
		{
			"a: 3s\n",
			map[string]string{
				"a": "3s",
			},
		},
		{
			"a: <foo>\n",
			map[string]string{"a": "<foo>"},
		},
		{
			"a: \"1:1\"\n",
			map[string]string{"a": "1:1"},
		},
		{
			"a: 1.2.3.4\n",
			map[string]string{"a": "1.2.3.4"},
		},
		{
			"a: 'b: c'\n",
			map[string]string{"a": "b: c"},
		},
		{
			"a: 'Hello #comment'\n",
			map[string]string{"a": "Hello #comment"},
		},
		{
			"a: 100.5\n",
			map[string]interface{}{
				"a": 100.5,
			},
		},
		{
			"a: \"\\0\"\n",
			map[string]string{"a": "\\0"},
		},
		{
			"b: 2\na: 1\nd: 4\nc: 3\nsub:\n  e: 5\n",
			map[string]interface{}{
				"b": 2,
				"a": 1,
				"d": 4,
				"c": 3,
				"sub": map[string]int{
					"e": 5,
				},
			},
		},
		{
			"       a       :          b        \n",
			map[string]string{"a": "b"},
		},
		{
			"a: b # comment\nb: c\n",
			map[string]string{
				"a": "b",
				"b": "c",
			},
		},
		{
			"---\na: b\n",
			map[string]string{"a": "b"},
		},
		{
			"a: b\n...\n",
			map[string]string{"a": "b"},
		},
		{
			"%YAML 1.2\n---\n",
			(*struct{})(nil),
		},
		{
			"---\n",
			(*struct{})(nil),
		},
		{
			"...",
			(*struct{})(nil),
		},
		{
			"v: go test ./...",
			map[string]string{"v": "go test ./..."},
		},
		{
			"v: echo ---",
			map[string]string{"v": "echo ---"},
		},
		{
			"v: |\n  hello\n  ...\n  world\n",
			map[string]string{"v": "hello\n...\nworld\n"},
		},
		{
			"a: !!binary gIGC\n",
			map[string]string{"a": "\x80\x81\x82"},
		},
		{
			"a: !!binary |\n  " + strings.Repeat("kJCQ", 17) + "kJ\n  CQ\n",
			map[string]string{"a": strings.Repeat("\x90", 54)},
		},
		{
			"v:\n- A\n- |-\n  B\n  C\n",
			map[string][]string{
				"v": {
					"A", "B\nC",
				},
			},
		},
		{
			"v:\n- A\n- >-\n  B\n  C\n",
			map[string][]string{
				"v": {
					"A", "B C",
				},
			},
		},
		{
			"a: b\nc: d\n",
			struct {
				A string
				C string `yaml:"c"`
			}{
				"b", "d",
			},
		},
		{
			"a: 1\nb: 2\n",
			struct {
				A int
				B int `yaml:"-"`
			}{
				1, 0,
			},
		},
		{
			"a: 1\nb: 2\n",
			struct {
				A     int
				Child `yaml:",inline"`
			}{
				1,
				Child{
					B: 2,
					C: 0,
				},
			},
		},

		// Anchors and aliases.
		{
			"a: &x 1\nb: &y 2\nc: *x\nd: *y\n",
			struct{ A, B, C, D int }{1, 2, 1, 2},
		},
		{
			"a: &a {c: 1}\nb: *a\n",
			struct {
				A, B struct {
					C int
				}
			}{struct{ C int }{1}, struct{ C int }{1}},
		},
		{
			"a: &a [1, 2]\nb: *a\n",
			struct{ B []int }{[]int{1, 2}},
		},

		{
			"tags:\n- hello-world\na: foo",
			struct {
				Tags []string
				A    string
			}{Tags: []string{"hello-world"}, A: "foo"},
		},
		{
			"",
			(*struct{})(nil),
		},
		{
			"{}", struct{}{},
		},
		{
			"v: /a/{b}",
			map[string]string{"v": "/a/{b}"},
		},
		{
			"v: 1[]{},!%?&*",
			map[string]string{"v": "1[]{},!%?&*"},
		},
		{
			"v: user's item",
			map[string]string{"v": "user's item"},
		},
		{
			"v: [1,[2,[3,[4,5],6],7],8]",
			map[string]interface{}{
				"v": []interface{}{
					1,
					[]interface{}{
						2,
						[]interface{}{
							3,
							[]int{4, 5},
							6,
						},
						7,
					},
					8,
				},
			},
		},
		{
			"v: {a: {b: {c: {d: e},f: g},h: i},j: k}",
			map[string]interface{}{
				"v": map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": map[string]string{
								"d": "e",
							},
							"f": "g",
						},
						"h": "i",
					},
					"j": "k",
				},
			},
		},
		{
			`---
- a:
    b:
- c: d
`,
			[]map[string]interface{}{
				{
					"a": map[string]interface{}{
						"b": nil,
					},
				},
				{
					"c": "d",
				},
			},
		},
		{
			`---
a:
  b:
c: d
`,
			map[string]interface{}{
				"a": map[string]interface{}{
					"b": nil,
				},
				"c": "d",
			},
		},
		{
			`---
a:
b:
c:
`,
			map[string]interface{}{
				"a": nil,
				"b": nil,
				"c": nil,
			},
		},
		{
			`---
a: go test ./...
b:
c:
`,
			map[string]interface{}{
				"a": "go test ./...",
				"b": nil,
				"c": nil,
			},
		},
		{
			`---
a: |
  hello
  ...
  world
b:
c:
`,
			map[string]interface{}{
				"a": "hello\n...\nworld\n",
				"b": nil,
				"c": nil,
			},
		},

		// Multi bytes
		{
			"v: あいうえお\nv2: かきくけこ",
			map[string]string{"v": "あいうえお", "v2": "かきくけこ"},
		},
	}
	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			buf := bytes.NewBufferString(test.source)
			dec := yaml.NewDecoder(buf)
			typ := reflect.ValueOf(test.value).Type()
			value := reflect.New(typ)
			if err := dec.Decode(value.Interface()); err != nil {
				if err == io.EOF {
					return
				}
				t.Fatalf("%s: %+v", test.source, err)
			}
			actual := fmt.Sprintf("%+v", value.Elem().Interface())
			expect := fmt.Sprintf("%+v", test.value)
			if actual != expect {
				t.Fatalf("failed to test [%s], actual=[%s], expect=[%s]", test.source, actual, expect)
			}
		})
	}
}

func TestDecoder_TypeConversionError(t *testing.T) {
	t.Run("type conversion for struct", func(t *testing.T) {
		type T struct {
			A int
			B uint
			C float32
			D bool
		}
		type U struct {
			*T `yaml:",inline"`
		}
		t.Run("string to int", func(t *testing.T) {
			var v T
			err := yaml.Unmarshal([]byte(`a: str`), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := "cannot unmarshal string into Go struct field T.A of type int"
			if err.Error() != msg {
				t.Fatalf("unexpected error message: %s. expect: %s", err.Error(), msg)
			}
		})
		t.Run("string to bool", func(t *testing.T) {
			var v T
			err := yaml.Unmarshal([]byte(`d: str`), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := "cannot unmarshal string into Go struct field T.D of type bool"
			if err.Error() != msg {
				t.Fatalf("unexpected error message: %s. expect: %s", err.Error(), msg)
			}
		})
		t.Run("string to int at inline", func(t *testing.T) {
			var v U
			err := yaml.Unmarshal([]byte(`a: str`), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := "cannot unmarshal string into Go struct field U.T.A of type int"
			if err.Error() != msg {
				t.Fatalf("unexpected error message: %s. expect: %s", err.Error(), msg)
			}
		})
	})
	t.Run("type conversion for array", func(t *testing.T) {
		t.Run("string to int", func(t *testing.T) {
			var v map[string][]int
			err := yaml.Unmarshal([]byte(`v: [A,1,C]`), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := "cannot unmarshal string into Go value of type int"
			if err.Error() != msg {
				t.Fatalf("unexpected error message: %s. expect: %s", err.Error(), msg)
			}
			if len(v) == 0 || len(v["v"]) == 0 {
				t.Fatal("failed to decode value")
			}
			if v["v"][0] != 1 {
				t.Fatal("failed to decode value")
			}
		})
		t.Run("string to int", func(t *testing.T) {
			var v map[string][]int
			err := yaml.Unmarshal([]byte("v:\n - A\n - 1\n - C"), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := "cannot unmarshal string into Go value of type int"
			if err.Error() != msg {
				t.Fatalf("unexpected error message: %s. expect: %s", err.Error(), msg)
			}
			if len(v) == 0 || len(v["v"]) == 0 {
				t.Fatal("failed to decode value")
			}
			if v["v"][0] != 1 {
				t.Fatal("failed to decode value")
			}
		})
	})
	t.Run("overflow error", func(t *testing.T) {
		t.Run("negative number to uint", func(t *testing.T) {
			var v map[string]uint
			err := yaml.Unmarshal([]byte("v: -42"), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := "cannot unmarshal -42 into Go value of type uint ( overflow )"
			if err.Error() != msg {
				t.Fatalf("unexpected error message: %s. expect: %s", err.Error(), msg)
			}
			if v["v"] != 0 {
				t.Fatal("failed to decode value")
			}
		})
		t.Run("negative number to uint64", func(t *testing.T) {
			var v map[string]uint64
			err := yaml.Unmarshal([]byte("v: -4294967296"), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := "cannot unmarshal -4294967296 into Go value of type uint64 ( overflow )"
			if err.Error() != msg {
				t.Fatalf("unexpected error message: %s. expect: %s", err.Error(), msg)
			}
			if v["v"] != 0 {
				t.Fatal("failed to decode value")
			}
		})
		t.Run("larger number for int32", func(t *testing.T) {
			var v map[string]int32
			err := yaml.Unmarshal([]byte("v: 4294967297"), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := "cannot unmarshal 4294967297 into Go value of type int32 ( overflow )"
			if err.Error() != msg {
				t.Fatalf("unexpected error message: %s. expect: %s", err.Error(), msg)
			}
			if v["v"] != 0 {
				t.Fatal("failed to decode value")
			}
		})
		t.Run("larger number for int8", func(t *testing.T) {
			var v map[string]int8
			err := yaml.Unmarshal([]byte("v: 128"), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := "cannot unmarshal 128 into Go value of type int8 ( overflow )"
			if err.Error() != msg {
				t.Fatalf("unexpected error message: %s. expect: %s", err.Error(), msg)
			}
			if v["v"] != 0 {
				t.Fatal("failed to decode value")
			}
		})
	})
	t.Run("type conversion for time", func(t *testing.T) {
		type T struct {
			A time.Time
			B time.Duration
		}
		t.Run("int to time", func(t *testing.T) {
			var v T
			err := yaml.Unmarshal([]byte(`a: 123`), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := "cannot unmarshal uint64 into Go struct field T.A of type time.Time"
			if err.Error() != msg {
				t.Fatalf("unexpected error message: %s. expect: %s", err.Error(), msg)
			}
		})
		t.Run("string to duration", func(t *testing.T) {
			var v T
			err := yaml.Unmarshal([]byte(`b: str`), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := `time: invalid duration "str"`
			if err.Error() != msg {
				t.Fatalf("unexpected error message: %s. expect: %s", err.Error(), msg)
			}
		})
		t.Run("int to duration", func(t *testing.T) {
			var v T
			err := yaml.Unmarshal([]byte(`b: 10`), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := "cannot unmarshal uint64 into Go struct field T.B of type time.Duration"
			if err.Error() != msg {
				t.Fatalf("unexpected error message: %s. expect: %s", err.Error(), msg)
			}
		})
	})
}

func TestDecoder_AnchorReferenceDirs(t *testing.T) {
	buf := bytes.NewBufferString("a: *a\n")
	dec := yaml.NewDecoder(buf, yaml.ReferenceDirs("testdata"))
	var v struct {
		A struct {
			B int
			C string
		}
	}
	if err := dec.Decode(&v); err != nil {
		t.Fatalf("%+v", err)
	}
	if v.A.B != 1 {
		t.Fatal("failed to decode by reference dirs")
	}
	if v.A.C != "hello" {
		t.Fatal("failed to decode by reference dirs")
	}
}

func TestDecoder_AnchorReferenceDirsRecursive(t *testing.T) {
	buf := bytes.NewBufferString("a: *a\n")
	dec := yaml.NewDecoder(
		buf,
		yaml.RecursiveDir(true),
		yaml.ReferenceDirs("testdata"),
	)
	var v struct {
		A struct {
			B int
			C string
		}
	}
	if err := dec.Decode(&v); err != nil {
		t.Fatalf("%+v", err)
	}
	if v.A.B != 1 {
		t.Fatal("failed to decode by reference dirs")
	}
	if v.A.C != "hello" {
		t.Fatal("failed to decode by reference dirs")
	}
}

func TestDecoder_AnchorFiles(t *testing.T) {
	buf := bytes.NewBufferString("a: *a\n")
	dec := yaml.NewDecoder(buf, yaml.ReferenceFiles("testdata/anchor.yml"))
	var v struct {
		A struct {
			B int
			C string
		}
	}
	if err := dec.Decode(&v); err != nil {
		t.Fatalf("%+v", err)
	}
	if v.A.B != 1 {
		t.Fatal("failed to decode by reference dirs")
	}
	if v.A.C != "hello" {
		t.Fatal("failed to decode by reference dirs")
	}
}

func TestDecodeWithMergeKey(t *testing.T) {
	yml := `
a: &a
  b: 1
  c: hello
items:
- <<: *a
- <<: *a
  c: world
`
	type Item struct {
		B int
		C string
	}
	type T struct {
		Items []*Item
	}
	buf := bytes.NewBufferString(yml)
	dec := yaml.NewDecoder(buf)
	var v T
	if err := dec.Decode(&v); err != nil {
		t.Fatalf("%+v", err)
	}
	if len(v.Items) != 2 {
		t.Fatal("failed to decode with merge key")
	}
	if v.Items[0].B != 1 || v.Items[0].C != "hello" {
		t.Fatal("failed to decode with merge key")
	}
	if v.Items[1].B != 1 || v.Items[1].C != "world" {
		t.Fatal("failed to decode with merge key")
	}
	t.Run("decode with interface{}", func(t *testing.T) {
		buf := bytes.NewBufferString(yml)
		dec := yaml.NewDecoder(buf)
		var v interface{}
		if err := dec.Decode(&v); err != nil {
			t.Fatalf("%+v", err)
		}
		items := v.(map[string]interface{})["items"].([]interface{})
		if len(items) != 2 {
			t.Fatal("failed to decode with merge key")
		}
		b0 := items[0].(map[string]interface{})["b"]
		if _, ok := b0.(uint64); !ok {
			t.Fatal("failed to decode with merge key")
		}
		if b0.(uint64) != 1 {
			t.Fatal("failed to decode with merge key")
		}
		c0 := items[0].(map[string]interface{})["c"]
		if _, ok := c0.(string); !ok {
			t.Fatal("failed to decode with merge key")
		}
		if c0.(string) != "hello" {
			t.Fatal("failed to decode with merge key")
		}
		b1 := items[1].(map[string]interface{})["b"]
		if _, ok := b1.(uint64); !ok {
			t.Fatal("failed to decode with merge key")
		}
		if b1.(uint64) != 1 {
			t.Fatal("failed to decode with merge key")
		}
		c1 := items[1].(map[string]interface{})["c"]
		if _, ok := c1.(string); !ok {
			t.Fatal("failed to decode with merge key")
		}
		if c1.(string) != "world" {
			t.Fatal("failed to decode with merge key")
		}
	})
	t.Run("decode with map", func(t *testing.T) {
		var v struct {
			Items []map[string]interface{}
		}
		buf := bytes.NewBufferString(yml)
		dec := yaml.NewDecoder(buf)
		if err := dec.Decode(&v); err != nil {
			t.Fatalf("%+v", err)
		}
		if len(v.Items) != 2 {
			t.Fatal("failed to decode with merge key")
		}
		b0 := v.Items[0]["b"]
		if _, ok := b0.(uint64); !ok {
			t.Fatal("failed to decode with merge key")
		}
		if b0.(uint64) != 1 {
			t.Fatal("failed to decode with merge key")
		}
		c0 := v.Items[0]["c"]
		if _, ok := c0.(string); !ok {
			t.Fatal("failed to decode with merge key")
		}
		if c0.(string) != "hello" {
			t.Fatal("failed to decode with merge key")
		}
		b1 := v.Items[1]["b"]
		if _, ok := b1.(uint64); !ok {
			t.Fatal("failed to decode with merge key")
		}
		if b1.(uint64) != 1 {
			t.Fatal("failed to decode with merge key")
		}
		c1 := v.Items[1]["c"]
		if _, ok := c1.(string); !ok {
			t.Fatal("failed to decode with merge key")
		}
		if c1.(string) != "world" {
			t.Fatal("failed to decode with merge key")
		}
	})
}

func TestDecoder_Inline(t *testing.T) {
	type Base struct {
		A int
		B string
	}
	yml := `---
a: 1
b: hello
c: true
`
	var v struct {
		*Base `yaml:",inline"`
		C     bool
	}
	if err := yaml.NewDecoder(strings.NewReader(yml)).Decode(&v); err != nil {
		t.Fatalf("%+v", err)
	}
	if v.A != 1 {
		t.Fatal("failed to decode with inline key")
	}
	if v.B != "hello" {
		t.Fatal("failed to decode with inline key")
	}
	if !v.C {
		t.Fatal("failed to decode with inline key")
	}

	t.Run("multiple inline with strict", func(t *testing.T) {
		type Base struct {
			A int
			B string
		}
		type Base2 struct {
			Base *Base `yaml:",inline"`
		}
		yml := `---
a: 1
b: hello
`
		var v struct {
			Base2 *Base2 `yaml:",inline"`
		}
		if err := yaml.NewDecoder(strings.NewReader(yml), yaml.Strict()).Decode(&v); err != nil {
			t.Fatalf("%+v", err)
		}
		if v.Base2.Base.A != 1 {
			t.Fatal("failed to decode with inline key")
		}
		if v.Base2.Base.B != "hello" {
			t.Fatal("failed to decode with inline key")
		}
	})
}

func TestDecoder_InlineAndConflictKey(t *testing.T) {
	type Base struct {
		A int
		B string
	}
	yml := `---
a: 1
b: hello
c: true
`
	var v struct {
		*Base `yaml:",inline"`
		A     int
		C     bool
	}
	if err := yaml.NewDecoder(strings.NewReader(yml)).Decode(&v); err != nil {
		t.Fatalf("%+v", err)
	}
	if v.A != 1 {
		t.Fatal("failed to decode with inline key")
	}
	if v.B != "hello" {
		t.Fatal("failed to decode with inline key")
	}
	if !v.C {
		t.Fatal("failed to decode with inline key")
	}
	if v.Base.A != 0 {
		t.Fatal("failed to decode with inline key")
	}
}

func TestDecoder_InlineAndWrongTypeStrict(t *testing.T) {
	type Base struct {
		A int
		B string
	}
	yml := `---
a: notanint
b: hello
c: true
`
	var v struct {
		*Base `yaml:",inline"`
		C     bool
	}
	err := yaml.NewDecoder(strings.NewReader(yml), yaml.Strict()).Decode(&v)
	if err == nil {
		t.Fatalf("expected error")
	}

	//TODO: properly check if errors are colored/have source
	t.Logf("%s", err)
	t.Logf("%s", yaml.FormatError(err, true, false))
	t.Logf("%s", yaml.FormatError(err, false, true))
	t.Logf("%s", yaml.FormatError(err, true, true))
}

func TestDecoder_InvalidCases(t *testing.T) {
	const src = `---
a:
- b
  c: d
`
	var v struct {
		A []string
	}
	err := yaml.NewDecoder(strings.NewReader(src)).Decode(&v)
	if err == nil {
		t.Fatalf("expected error")
	}

	if err.Error() != yaml.FormatError(err, false, true) {
		t.Logf("err.Error() = %s", err.Error())
		t.Logf("yaml.FormatError(err, false, true) = %s", yaml.FormatError(err, false, true))
		t.Fatal(`err.Error() should match yaml.FormatError(err, false, true)`)
	}

	//TODO: properly check if errors are colored/have source
	t.Logf("%s", err)
	t.Logf("%s", yaml.FormatError(err, true, false))
	t.Logf("%s", yaml.FormatError(err, false, true))
	t.Logf("%s", yaml.FormatError(err, true, true))
}

func TestDecoder_JSONTags(t *testing.T) {
	var v struct {
		A string `json:"a_json"`               // no YAML tag
		B string `json:"b_json" yaml:"b_yaml"` // both tags
	}

	const src = `---
a_json: a_json_value
b_json: b_json_value
b_yaml: b_yaml_value
`
	if err := yaml.NewDecoder(strings.NewReader(src)).Decode(&v); err != nil {
		t.Fatalf(`parsing should succeed: %s`, err)
	}

	if v.A != "a_json_value" {
		t.Fatalf("v.A should be `a_json_value`, got `%s`", v.A)
	}

	if v.B != "b_yaml_value" {
		t.Fatalf("v.B should be `b_yaml_value`, got `%s`", v.B)
	}
}

func TestDecoder_DisallowUnknownField(t *testing.T) {
	t.Run("different level keys with same name", func(t *testing.T) {
		var v struct {
			C Child `yaml:"c"`
		}
		yml := `---
b: 1
c:
  b: 1
`

		err := yaml.NewDecoder(strings.NewReader(yml), yaml.DisallowUnknownField()).Decode(&v)
		if err == nil {
			t.Fatalf("error expected")
		}
	})
	t.Run("inline", func(t *testing.T) {
		var v struct {
			*Child `yaml:",inline"`
			A      string `yaml:"a"`
		}
		yml := `---
a: a
b: 1
`

		if err := yaml.NewDecoder(strings.NewReader(yml), yaml.DisallowUnknownField()).Decode(&v); err != nil {
			t.Fatalf(`parsing should succeed: %s`, err)
		}
		if v.A != "a" {
			t.Fatalf("v.A should be `a`, got `%s`", v.A)
		}
		if v.B != 1 {
			t.Fatalf("v.B should be 1, got %d", v.B)
		}
		if v.C != 0 {
			t.Fatalf("v.C should be 0, got %d", v.C)
		}
	})
	t.Run("list", func(t *testing.T) {
		type C struct {
			Child `yaml:",inline"`
		}

		var v struct {
			Children []C `yaml:"children"`
		}

		yml := `---
children:
- b: 1
- b: 2
`

		if err := yaml.NewDecoder(strings.NewReader(yml), yaml.DisallowUnknownField()).Decode(&v); err != nil {
			t.Fatalf(`parsing should succeed: %s`, err)
		}

		if len(v.Children) != 2 {
			t.Fatalf(`len(v.Children) should be 2, got %d`, len(v.Children))
		}

		if v.Children[0].B != 1 {
			t.Fatalf(`v.Children[0].B should be 1, got %d`, v.Children[0].B)
		}

		if v.Children[1].B != 2 {
			t.Fatalf(`v.Children[1].B should be 2, got %d`, v.Children[1].B)
		}
	})
}

func TestDecoder_DisallowDuplicateKey(t *testing.T) {
	yml := `
a: b
a: c
`
	expected := `
[3:1] duplicate key "a"
   2 | a: b
>  3 | a: c
       ^
`
	t.Run("map", func(t *testing.T) {
		var v map[string]string
		err := yaml.NewDecoder(strings.NewReader(yml), yaml.DisallowDuplicateKey()).Decode(&v)
		if err == nil {
			t.Fatal("decoding should fail")
		}
		actual := "\n" + err.Error()
		if expected != actual {
			t.Fatalf("expected:[%s] actual:[%s]", expected, actual)
		}
	})
	t.Run("struct", func(t *testing.T) {
		var v struct {
			A string
		}
		err := yaml.NewDecoder(strings.NewReader(yml), yaml.DisallowDuplicateKey()).Decode(&v)
		if err == nil {
			t.Fatal("decoding should fail")
		}
		actual := "\n" + err.Error()
		if expected != actual {
			t.Fatalf("expected:[%s] actual:[%s]", expected, actual)
		}
	})
}

func TestDecoder_DefaultValues(t *testing.T) {
	v := struct {
		A string `yaml:"a"`
		B string `yaml:"b"`
		c string // private
		D struct {
			E string `yaml:"e"`
			F struct {
				G string `yaml:"g"`
			} `yaml:"f"`
			H struct {
				I string `yaml:"i"`
			} `yaml:",inline"`
		} `yaml:"d"`
		J struct {
			K string `yaml:"k"`
			L struct {
				M string `yaml:"m"`
			} `yaml:"l"`
			N struct {
				O string `yaml:"o"`
			} `yaml:",inline"`
		} `yaml:",inline"`
		P struct {
			Q string `yaml:"q"`
			R struct {
				S string `yaml:"s"`
			} `yaml:"r"`
			T struct {
				U string `yaml:"u"`
			} `yaml:",inline"`
		} `yaml:"p"`
		V struct {
			W string `yaml:"w"`
			X struct {
				Y string `yaml:"y"`
			} `yaml:"x"`
			Z struct {
				Ä string `yaml:"ä"`
			} `yaml:",inline"`
		} `yaml:",inline"`
	}{
		B: "defaultBValue",
		c: "defaultCValue",
	}

	v.D.E = "defaultEValue"
	v.D.F.G = "defaultGValue"
	v.D.H.I = "defaultIValue"
	v.J.K = "defaultKValue"
	v.J.L.M = "defaultMValue"
	v.J.N.O = "defaultOValue"
	v.P.R.S = "defaultSValue"
	v.P.T.U = "defaultUValue"
	v.V.X.Y = "defaultYValue"
	v.V.Z.Ä = "defaultÄValue"

	const src = `---
a: a_value
p:
   q: q_value
w: w_value
`
	if err := yaml.NewDecoder(strings.NewReader(src)).Decode(&v); err != nil {
		t.Fatalf(`parsing should succeed: %s`, err)
	}
	if v.A != "a_value" {
		t.Fatalf("v.A should be `a_value`, got `%s`", v.A)
	}

	if v.B != "defaultBValue" {
		t.Fatalf("v.B should be `defaultValue`, got `%s`", v.B)
	}

	if v.c != "defaultCValue" {
		t.Fatalf("v.c should be `defaultCValue`, got `%s`", v.c)
	}

	if v.D.E != "defaultEValue" {
		t.Fatalf("v.D.E should be `defaultEValue`, got `%s`", v.D.E)
	}

	if v.D.F.G != "defaultGValue" {
		t.Fatalf("v.D.F.G should be `defaultGValue`, got `%s`", v.D.F.G)
	}

	if v.D.H.I != "defaultIValue" {
		t.Fatalf("v.D.H.I should be `defaultIValue`, got `%s`", v.D.H.I)
	}

	if v.J.K != "defaultKValue" {
		t.Fatalf("v.J.K should be `defaultKValue`, got `%s`", v.J.K)
	}

	if v.J.L.M != "defaultMValue" {
		t.Fatalf("v.J.L.M should be `defaultMValue`, got `%s`", v.J.L.M)
	}

	if v.J.N.O != "defaultOValue" {
		t.Fatalf("v.J.N.O should be `defaultOValue`, got `%s`", v.J.N.O)
	}

	if v.P.Q != "q_value" {
		t.Fatalf("v.P.Q should be `q_value`, got `%s`", v.P.Q)
	}

	if v.P.R.S != "defaultSValue" {
		t.Fatalf("v.P.R.S should be `defaultSValue`, got `%s`", v.P.R.S)
	}

	if v.P.T.U != "defaultUValue" {
		t.Fatalf("v.P.T.U should be `defaultUValue`, got `%s`", v.P.T.U)
	}

	if v.V.W != "w_value" {
		t.Fatalf("v.V.W should be `w_value`, got `%s`", v.V.W)
	}

	if v.V.X.Y != "defaultYValue" {
		t.Fatalf("v.V.X.Y should be `defaultYValue`, got `%s`", v.V.X.Y)
	}

	if v.V.Z.Ä != "defaultÄValue" {
		t.Fatalf("v.V.Z.Ä should be `defaultÄValue`, got `%s`", v.V.Z.Ä)
	}
}

func Example_YAMLTags() {
	yml := `---
foo: 1
bar: c
A: 2
B: d
`
	var v struct {
		A int    `yaml:"foo" json:"A"`
		B string `yaml:"bar" json:"B"`
	}
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		log.Fatal(err)
	}
	fmt.Println(v.A)
	fmt.Println(v.B)
	// OUTPUT:
	// 1
	// c
}

type useJSONUnmarshalerTest struct {
	s string
}

func (t *useJSONUnmarshalerTest) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	t.s = s
	return nil
}

func TestDecoder_UseJSONUnmarshaler(t *testing.T) {
	var v useJSONUnmarshalerTest
	if err := yaml.UnmarshalWithOptions([]byte(`"a"`), &v, yaml.UseJSONUnmarshaler()); err != nil {
		t.Fatal(err)
	}
	if v.s != "a" {
		t.Fatalf("unexpected decoded value: %s", v.s)
	}
}

type unmarshalContext struct {
	v int
}

func (c *unmarshalContext) UnmarshalYAML(ctx context.Context, b []byte) error {
	v, ok := ctx.Value("k").(int)
	if !ok {
		return fmt.Errorf("cannot get valid context")
	}
	if v != 1 {
		return fmt.Errorf("cannot get valid context")
	}
	if string(b) != "1" {
		return fmt.Errorf("cannot get valid bytes")
	}
	c.v = v
	return nil
}

func Test_UnmarshalerContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "k", 1)
	var v unmarshalContext
	if err := yaml.UnmarshalContext(ctx, []byte(`1`), &v); err != nil {
		t.Fatalf("%+v", err)
	}
	if v.v != 1 {
		t.Fatal("cannot call UnmarshalYAML")
	}
}

func TestDecoder_DecodeFromNode(t *testing.T) {
	t.Run("has reference", func(t *testing.T) {
		str := `
anchor: &map
  text: hello
map: *map`
		var buf bytes.Buffer
		dec := yaml.NewDecoder(&buf)
		f, err := parser.ParseBytes([]byte(str), 0)
		if err != nil {
			t.Fatalf("failed to parse: %s", err)
		}
		type T struct {
			Map map[string]string
		}
		var v T
		if err := dec.DecodeFromNode(f.Docs[0].Body, &v); err != nil {
			t.Fatalf("failed to decode: %s", err)
		}
		actual := fmt.Sprintf("%+v", v)
		expect := fmt.Sprintf("%+v", T{map[string]string{"text": "hello"}})
		if actual != expect {
			t.Fatalf("actual=[%s], expect=[%s]", actual, expect)
		}
	})
	t.Run("with reference option", func(t *testing.T) {
		anchor := strings.NewReader(`
map: &map
  text: hello`)
		var buf bytes.Buffer
		dec := yaml.NewDecoder(&buf, yaml.ReferenceReaders(anchor))
		f, err := parser.ParseBytes([]byte("map: *map"), 0)
		if err != nil {
			t.Fatalf("failed to parse: %s", err)
		}
		type T struct {
			Map map[string]string
		}
		var v T
		if err := dec.DecodeFromNode(f.Docs[0].Body, &v); err != nil {
			t.Fatalf("failed to decode: %s", err)
		}
		actual := fmt.Sprintf("%+v", v)
		expect := fmt.Sprintf("%+v", T{map[string]string{"text": "hello"}})
		if actual != expect {
			t.Fatalf("actual=[%s], expect=[%s]", actual, expect)
		}
	})
	t.Run("value is not pointer", func(t *testing.T) {
		var buf bytes.Buffer
		var v bool
		err := yaml.NewDecoder(&buf).DecodeFromNode(nil, v)
		if !xerrors.Is(err, errors.ErrDecodeRequiredPointerType) {
			t.Fatalf("unexpected error: %s", err)
		}
	})
}

func Example_JSONTags() {
	yml := `---
foo: 1
bar: c
`
	var v struct {
		A int    `json:"foo"`
		B string `json:"bar"`
	}
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		log.Fatal(err)
	}
	fmt.Println(v.A)
	fmt.Println(v.B)
	// OUTPUT:
	// 1
	// c
}

func Example_DisallowUnknownField() {
	var v struct {
		A string `yaml:"simple"`
		C string `yaml:"complicated"`
	}

	const src = `---
simple: string
complecated: string
`
	err := yaml.NewDecoder(strings.NewReader(src), yaml.DisallowUnknownField()).Decode(&v)
	fmt.Printf("%v\n", err)

	// OUTPUT:
	// [3:1] unknown field "complecated"
	//        1 | ---
	//        2 | simple: string
	//     >  3 | complecated: string
	//            ^
}

func Example_Unmarshal_Node() {
	f, err := parser.ParseBytes([]byte("text: node example"), 0)
	if err != nil {
		panic(err)
	}
	var v struct {
		Text string `yaml:"text"`
	}
	if err := yaml.NodeToValue(f.Docs[0].Body, &v); err != nil {
		panic(err)
	}
	fmt.Println(v.Text)
	// OUTPUT:
	// node example
}

type unmarshalableYAMLStringValue string

func (v *unmarshalableYAMLStringValue) UnmarshalYAML(b []byte) error {
	var s string
	if err := yaml.Unmarshal(b, &s); err != nil {
		return err
	}
	*v = unmarshalableYAMLStringValue(s)
	return nil
}

type unmarshalableTextStringValue string

func (v *unmarshalableTextStringValue) UnmarshalText(b []byte) error {
	*v = unmarshalableTextStringValue(string(b))
	return nil
}

type unmarshalableStringContainer struct {
	A unmarshalableYAMLStringValue `yaml:"a"`
	B unmarshalableTextStringValue `yaml:"b"`
}

func TestUnmarshalableString(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		yml := `
a: ""
b: ""
`
		var container unmarshalableStringContainer
		if err := yaml.Unmarshal([]byte(yml), &container); err != nil {
			t.Fatalf("failed to unmarshal %v", err)
		}
		if container.A != "" {
			t.Fatalf("expected empty string, but %q is set", container.A)
		}
		if container.B != "" {
			t.Fatalf("expected empty string, but %q is set", container.B)
		}
	})
	t.Run("filled string", func(t *testing.T) {
		t.Parallel()
		yml := `
a: "aaa"
b: "bbb"
`
		var container unmarshalableStringContainer
		if err := yaml.Unmarshal([]byte(yml), &container); err != nil {
			t.Fatalf("failed to unmarshal %v", err)
		}
		if container.A != "aaa" {
			t.Fatalf("expected \"aaa\", but %q is set", container.A)
		}
		if container.B != "bbb" {
			t.Fatalf("expected \"bbb\", but %q is set", container.B)
		}
	})
	t.Run("single-quoted string", func(t *testing.T) {
		t.Parallel()
		yml := `
a: 'aaa'
b: 'bbb'
`
		var container unmarshalableStringContainer
		if err := yaml.Unmarshal([]byte(yml), &container); err != nil {
			t.Fatalf("failed to unmarshal %v", err)
		}
		if container.A != "aaa" {
			t.Fatalf("expected \"aaa\", but %q is set", container.A)
		}
		if container.B != "bbb" {
			t.Fatalf("expected \"aaa\", but %q is set", container.B)
		}
	})
	t.Run("literal", func(t *testing.T) {
		t.Parallel()
		yml := `
a: |
 a
 b
 c
b: |
 a
 b
 c
`
		var container unmarshalableStringContainer
		if err := yaml.Unmarshal([]byte(yml), &container); err != nil {
			t.Fatalf("failed to unmarshal %v", err)
		}
		if container.A != "a\nb\nc\n" {
			t.Fatalf("expected \"a\nb\nc\n\", but %q is set", container.A)
		}
		if container.B != "a\nb\nc\n" {
			t.Fatalf("expected \"a\nb\nc\n\", but %q is set", container.B)
		}
	})
	t.Run("anchor/alias", func(t *testing.T) {
		yml := `
a: &x 1
b: *x
c: &y hello
d: *y
`
		var v struct {
			A, B, C, D unmarshalableTextStringValue
		}
		if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
			t.Fatal(err)
		}
		if v.A != "1" {
			t.Fatal("failed to unmarshal")
		}
		if v.B != "1" {
			t.Fatal("failed to unmarshal")
		}
		if v.C != "hello" {
			t.Fatal("failed to unmarshal")
		}
		if v.D != "hello" {
			t.Fatal("failed to unmarshal")
		}
	})
	t.Run("net.IP", func(t *testing.T) {
		yml := `
a: &a 127.0.0.1
b: *a
`
		var v struct {
			A, B net.IP
		}
		if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
			t.Fatal(err)
		}
		if v.A.String() != net.IPv4(127, 0, 0, 1).String() {
			t.Fatal("failed to unmarshal")
		}
		if v.B.String() != net.IPv4(127, 0, 0, 1).String() {
			t.Fatal("failed to unmarshal")
		}
	})
}

type unmarshalablePtrStringContainer struct {
	V *string `yaml:"value"`
}

func TestUnmarshalablePtrString(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		var container unmarshalablePtrStringContainer
		if err := yaml.Unmarshal([]byte(`value: ""`), &container); err != nil {
			t.Fatalf("failed to unmarshal %v", err)
		}
		if container.V == nil || *container.V != "" {
			t.Fatalf("expected empty string, but %q is set", *container.V)
		}
	})

	t.Run("null", func(t *testing.T) {
		t.Parallel()
		var container unmarshalablePtrStringContainer
		if err := yaml.Unmarshal([]byte(`value: null`), &container); err != nil {
			t.Fatalf("failed to unmarshal %v", err)
		}
		if container.V != (*string)(nil) {
			t.Fatalf("expected nil, but %q is set", *container.V)
		}
	})
}

type unmarshalableIntValue int

func (v *unmarshalableIntValue) UnmarshalYAML(raw []byte) error {
	i, err := strconv.Atoi(string(raw))
	if err != nil {
		return err
	}
	*v = unmarshalableIntValue(i)
	return nil
}

type unmarshalableIntContainer struct {
	V unmarshalableIntValue `yaml:"value"`
}

func TestUnmarshalableInt(t *testing.T) {
	t.Run("empty int", func(t *testing.T) {
		t.Parallel()
		var container unmarshalableIntContainer
		if err := yaml.Unmarshal([]byte(``), &container); err != nil {
			t.Fatalf("failed to unmarshal %v", err)
		}
		if container.V != 0 {
			t.Fatalf("expected empty int, but %d is set", container.V)
		}
	})
	t.Run("filled int", func(t *testing.T) {
		t.Parallel()
		var container unmarshalableIntContainer
		if err := yaml.Unmarshal([]byte(`value: 9`), &container); err != nil {
			t.Fatalf("failed to unmarshal %v", err)
		}
		if container.V != 9 {
			t.Fatalf("expected 9, but %d is set", container.V)
		}
	})
	t.Run("filled number", func(t *testing.T) {
		t.Parallel()
		var container unmarshalableIntContainer
		if err := yaml.Unmarshal([]byte(`value: 9`), &container); err != nil {
			t.Fatalf("failed to unmarshal %v", err)
		}
		if container.V != 9 {
			t.Fatalf("expected 9, but %d is set", container.V)
		}
	})
}

type unmarshalablePtrIntContainer struct {
	V *int `yaml:"value"`
}

func TestUnmarshalablePtrInt(t *testing.T) {
	t.Run("empty int", func(t *testing.T) {
		t.Parallel()
		var container unmarshalablePtrIntContainer
		if err := yaml.Unmarshal([]byte(`value: 0`), &container); err != nil {
			t.Fatalf("failed to unmarshal %v", err)
		}
		if container.V == nil || *container.V != 0 {
			t.Fatalf("expected 0, but %q is set", *container.V)
		}
	})

	t.Run("null", func(t *testing.T) {
		t.Parallel()
		var container unmarshalablePtrIntContainer
		if err := yaml.Unmarshal([]byte(`value: null`), &container); err != nil {
			t.Fatalf("failed to unmarshal %v", err)
		}
		if container.V != (*int)(nil) {
			t.Fatalf("expected nil, but %q is set", *container.V)
		}
	})
}

type literalContainer struct {
	v string
}

func (c *literalContainer) UnmarshalYAML(v []byte) error {
	var lit string
	if err := yaml.Unmarshal(v, &lit); err != nil {
		return err
	}
	c.v = lit
	return nil
}

func TestDecode_Literal(t *testing.T) {
	yml := `---
value: |
  {
     "key": "value"
  }
`
	var v map[string]*literalContainer
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		t.Fatalf("failed to unmarshal %+v", err)
	}
	if v["value"] == nil {
		t.Fatal("failed to unmarshal literal with bytes unmarshaler")
	}
	if v["value"].v == "" {
		t.Fatal("failed to unmarshal literal with bytes unmarshaler")
	}
}

func TestDecoder_UseOrderedMap(t *testing.T) {
	yml := `
a: b
c: d
e:
  f: g
  h: i
j: k
`
	var v interface{}
	if err := yaml.NewDecoder(strings.NewReader(yml), yaml.UseOrderedMap()).Decode(&v); err != nil {
		t.Fatalf("%+v", err)
	}
	if _, ok := v.(yaml.MapSlice); !ok {
		t.Fatalf("failed to convert to ordered map: %T", v)
	}
	bytes, err := yaml.Marshal(v)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if string(yml) != "\n"+string(bytes) {
		t.Fatalf("expected:[%s] actual:[%s]", string(yml), "\n"+string(bytes))
	}
}

func TestDecoder_Stream(t *testing.T) {
	yml := `
---
a: b
c: d
---
e: f
g: h
---
i: j
k: l
`
	dec := yaml.NewDecoder(strings.NewReader(yml))
	values := []map[string]string{}
	for {
		var v map[string]string
		if err := dec.Decode(&v); err != nil {
			if err == io.EOF {
				break
			}
			t.Fatalf("%+v", err)
		}
		values = append(values, v)
	}
	if len(values) != 3 {
		t.Fatal("failed to stream decoding")
	}
	if values[0]["a"] != "b" {
		t.Fatal("failed to stream decoding")
	}
	if values[1]["e"] != "f" {
		t.Fatal("failed to stream decoding")
	}
	if values[2]["i"] != "j" {
		t.Fatal("failed to stream decoding")
	}
}

type unmarshalYAMLWithAliasString string

func (v *unmarshalYAMLWithAliasString) UnmarshalYAML(b []byte) error {
	var s string
	if err := yaml.Unmarshal(b, &s); err != nil {
		return err
	}
	*v = unmarshalYAMLWithAliasString(s)
	return nil
}

type unmarshalYAMLWithAliasMap map[string]interface{}

func (v *unmarshalYAMLWithAliasMap) UnmarshalYAML(b []byte) error {
	var m map[string]interface{}
	if err := yaml.Unmarshal(b, &m); err != nil {
		return err
	}
	*v = unmarshalYAMLWithAliasMap(m)
	return nil
}

func TestDecoder_UnmarshalYAMLWithAlias(t *testing.T) {
	yml := `
anchors:
 x: &x "\"hello\" \"world\""
 map: &y
   a: b
   c: d
   d: *x
a: *x
b:
 <<: *y
 e: f
`
	var v struct {
		A unmarshalYAMLWithAliasString
		B unmarshalYAMLWithAliasMap
	}
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		t.Fatalf("%+v", err)
	}
	if v.A != `"hello" "world"` {
		t.Fatal("failed to unmarshal with alias")
	}
	if len(v.B) != 4 {
		t.Fatal("failed to unmarshal with alias")
	}
	if v.B["a"] != "b" {
		t.Fatal("failed to unmarshal with alias")
	}
	if v.B["c"] != "d" {
		t.Fatal("failed to unmarshal with alias")
	}
	if v.B["d"] != `"hello" "world"` {
		t.Fatal("failed to unmarshal with alias")
	}
}

type unmarshalString string

func (u *unmarshalString) UnmarshalYAML(b []byte) error {
	*u = unmarshalString(string(b))
	return nil
}

type unmarshalList struct {
	v []map[string]unmarshalString
}

func (u *unmarshalList) UnmarshalYAML(b []byte) error {
	expected := `
 - b: c
   d: |
     hello

     hello
   f: g
 - h: i`
	actual := "\n" + string(b)
	if expected != actual {
		return xerrors.Errorf("unexpected bytes: expected [%q] but got [%q]", expected, actual)
	}
	var v []map[string]unmarshalString
	if err := yaml.Unmarshal(b, &v); err != nil {
		return err
	}
	u.v = v
	return nil
}

func TestDecoder_UnmarshalBytesWithSeparatedList(t *testing.T) {
	yml := `
a:
 - b: c
   d: |
     hello

     hello
   f: g
 - h: i
`
	var v struct {
		A unmarshalList
	}
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		t.Fatal(err)
	}
	if len(v.A.v) != 2 {
		t.Fatalf("failed to unmarshal %+v", v)
	}
	if len(v.A.v[0]) != 3 {
		t.Fatalf("failed to unmarshal %+v", v.A.v[0])
	}
	if len(v.A.v[1]) != 1 {
		t.Fatalf("failed to unmarshal %+v", v.A.v[1])
	}
}

func TestDecoder_LiteralWithNewLine(t *testing.T) {
	type A struct {
		Node     string `yaml:"b"`
		LastNode string `yaml:"last"`
	}
	tests := []A{
		A{
			Node: "hello\nworld",
		},
		A{
			Node: "hello\nworld\n",
		},
		A{
			Node: "hello\nworld\n\n",
		},
		A{
			LastNode: "hello\nworld",
		},
		A{
			LastNode: "hello\nworld\n",
		},
		A{
			LastNode: "hello\nworld\n\n",
		},
	}
	// struct(want) -> Marshal -> Unmarchal -> struct(got)
	for _, want := range tests {
		bytes, _ := yaml.Marshal(want)
		got := A{}
		if err := yaml.Unmarshal(bytes, &got); err != nil {
			t.Fatal(err)
		}
		if want.Node != got.Node {
			t.Fatalf("expected:%q but got %q", want.Node, got.Node)
		}
		if want.LastNode != got.LastNode {
			t.Fatalf("expected:%q but got %q", want.LastNode, got.LastNode)
		}
	}
}

func TestDecoder_TabCharacterAtRight(t *testing.T) {
	yml := `
- a: [2 , 2] 			
  b: [2 , 2] 			
  c: [2 , 2]`
	var v []map[string][]int
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		t.Fatal(err)
	}
	if len(v) != 1 {
		t.Fatalf("failed to unmarshal %+v", v)
	}
	if len(v[0]) != 3 {
		t.Fatalf("failed to unmarshal %+v", v)
	}
}

func TestDecoder_Canonical(t *testing.T) {
	yml := `
!!map {
  ? !!str "explicit":!!str "entry",
  ? !!str "implicit" : !!str "entry",
  ? !!null "" : !!null "",
}
`
	var v interface{}
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		t.Fatalf("%+v", err)
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		t.Fatalf("failed to decode canonical yaml: %+v", v)
	}
	if m["explicit"] != "entry" {
		t.Fatalf("failed to decode canonical yaml: %+v", m)
	}
	if m["implicit"] != "entry" {
		t.Fatalf("failed to decode canonical yaml: %+v", m)
	}
	if m["null"] != nil {
		t.Fatalf("failed to decode canonical yaml: %+v", m)
	}
}

func TestDecoder_DecodeFromFile(t *testing.T) {
	yml := `
a: b
c: d
`
	file, err := parser.ParseBytes([]byte(yml), 0)
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]string
	if err := yaml.NewDecoder(file).Decode(&v); err != nil {
		t.Fatal(err)
	}
	if len(v) != 2 {
		t.Fatal("failed to decode from ast.File")
	}
	if v["a"] != "b" {
		t.Fatal("failed to decode from ast.File")
	}
	if v["c"] != "d" {
		t.Fatal("failed to decode from ast.File")
	}
}

func TestDecoder_DecodeWithNode(t *testing.T) {
	t.Run("abstract node", func(t *testing.T) {
		type T struct {
			Text ast.Node `yaml:"text"`
		}
		var v T
		if err := yaml.Unmarshal([]byte(`text: hello`), &v); err != nil {
			t.Fatalf("%+v", err)
		}
		expected := "hello"
		got := v.Text.String()
		if expected != got {
			t.Fatalf("failed to decode to ast.Node: expected %s but got %s", expected, got)
		}
	})
	t.Run("concrete node", func(t *testing.T) {
		type T struct {
			Text *ast.StringNode `yaml:"text"`
		}
		var v T
		if err := yaml.Unmarshal([]byte(`text: hello`), &v); err != nil {
			t.Fatalf("%+v", err)
		}
		expected := "hello"
		got := v.Text.String()
		if expected != got {
			t.Fatalf("failed to decode to ast.Node: expected %s but got %s", expected, got)
		}
	})
}

func TestRoundtripAnchorAlias(t *testing.T) {
	t.Run("irreversible", func(t *testing.T) {
		type foo struct {
			K1 string
			K2 string
		}

		type bar struct {
			K1 string
			K3 string
		}

		type doc struct {
			Foo foo
			Bar bar
		}
		yml := `
foo:
 <<: &test-anchor
   k1: "One"
 k2: "Two"

bar:
 <<: *test-anchor
 k3: "Three"
`
		var v doc
		if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
			t.Fatalf("%+v", err)
		}
		bytes, err := yaml.Marshal(v)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		expected := `
foo:
  k1: One
  k2: Two
bar:
  k1: One
  k3: Three
`
		got := "\n" + string(bytes)
		if expected != got {
			t.Fatalf("expected:[%s] but got [%s]", expected, got)
		}
	})
	t.Run("reversible", func(t *testing.T) {
		type TestAnchor struct {
			K1 string
		}
		type foo struct {
			*TestAnchor `yaml:",inline,alias"`
			K2          string
		}
		type bar struct {
			*TestAnchor `yaml:",inline,alias"`
			K3          string
		}
		type doc struct {
			TestAnchor *TestAnchor `yaml:"test-anchor,anchor"`
			Foo        foo
			Bar        bar
		}
		yml := `
test-anchor: &test-anchor
  k1: One
foo:
  <<: *test-anchor
  k2: Two
bar:
  <<: *test-anchor
  k3: Three
`
		var v doc
		if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
			t.Fatalf("%+v", err)
		}
		bytes, err := yaml.Marshal(v)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		got := "\n" + string(bytes)
		if yml != got {
			t.Fatalf("expected:[%s] but got [%s]", yml, got)
		}
	})
}
