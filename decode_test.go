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
)

type Child struct {
	B int
	C int `yaml:"-"`
}

type TestString string

func TestDecoder(t *testing.T) {
	tests := []struct {
		source string
		value  interface{}
		eof    bool
	}{
		{
			source: "v: hi\n",
			value:  map[string]string{"v": "hi"},
		},
		{
			source: "v: hi\n",
			value:  map[string]TestString{"v": "hi"},
		},
		{
			source: "v: \"true\"\n",
			value:  map[string]string{"v": "true"},
		},
		{
			source: "v: \"false\"\n",
			value:  map[string]string{"v": "false"},
		},
		{
			source: "v: true\n",
			value:  map[string]interface{}{"v": true},
		},
		{
			source: "v: true\n",
			value:  map[string]string{"v": "true"},
		},
		{
			source: "v: 10\n",
			value:  map[string]string{"v": "10"},
		},
		{
			source: "v: 10\n",
			value:  map[string]TestString{"v": "10"},
		},
		{
			source: "v: -10\n",
			value:  map[string]string{"v": "-10"},
		},
		{
			source: "v: 1.234\n",
			value:  map[string]string{"v": "1.234"},
		},
		{
			source: "v: \" foo\"\n",
			value:  map[string]string{"v": " foo"},
		},
		{
			source: "v: \"foo \"\n",
			value:  map[string]string{"v": "foo "},
		},
		{
			source: "v: \" foo \"\n",
			value:  map[string]string{"v": " foo "},
		},
		{
			source: "v: false\n",
			value:  map[string]bool{"v": false},
		},
		{
			source: "v: 10\n",
			value:  map[string]int{"v": 10},
		},
		{
			source: "v: 10",
			value:  map[string]interface{}{"v": 10},
		},
		{
			source: "v: 0b10",
			value:  map[string]interface{}{"v": 2},
		},
		{
			source: "v: -0b101010",
			value:  map[string]interface{}{"v": -42},
		},
		{
			source: "v: -0b1000000000000000000000000000000000000000000000000000000000000000",
			value:  map[string]interface{}{"v": int64(-9223372036854775808)},
		},
		{
			source: "v: 0xA",
			value:  map[string]interface{}{"v": 10},
		},
		{
			source: "v: .1",
			value:  map[string]interface{}{"v": 0.1},
		},
		{
			source: "v: -.1",
			value:  map[string]interface{}{"v": -0.1},
		},
		{
			source: "v: -10\n",
			value:  map[string]int{"v": -10},
		},
		{
			source: "v: 4294967296\n",
			value:  map[string]int64{"v": int64(4294967296)},
		},
		{
			source: "v: 0.1\n",
			value:  map[string]interface{}{"v": 0.1},
		},
		{
			source: "v: 0.99\n",
			value:  map[string]float32{"v": 0.99},
		},
		{
			source: "v: -0.1\n",
			value:  map[string]float64{"v": -0.1},
		},
		{
			source: "v: 6.8523e+5",
			value:  map[string]interface{}{"v": 6.8523e+5},
		},
		{
			source: "v: 685.230_15e+03",
			value:  map[string]interface{}{"v": 685.23015e+03},
		},
		{
			source: "v: 685_230.15",
			value:  map[string]interface{}{"v": 685230.15},
		},
		{
			source: "v: 685_230.15",
			value:  map[string]float64{"v": 685230.15},
		},
		{
			source: "v: 685230",
			value:  map[string]interface{}{"v": 685230},
		},
		{
			source: "v: +685_230",
			value:  map[string]interface{}{"v": 685230},
		},
		{
			source: "v: 02472256",
			value:  map[string]interface{}{"v": 685230},
		},
		{
			source: "v: 0x_0A_74_AE",
			value:  map[string]interface{}{"v": 685230},
		},
		{
			source: "v: 0b1010_0111_0100_1010_1110",
			value:  map[string]interface{}{"v": 685230},
		},
		{
			source: "v: +685_230",
			value:  map[string]int{"v": 685230},
		},

		// Bools from spec
		{
			source: "v: True",
			value:  map[string]interface{}{"v": true},
		},
		{
			source: "v: TRUE",
			value:  map[string]interface{}{"v": true},
		},
		{
			source: "v: False",
			value:  map[string]interface{}{"v": false},
		},
		{
			source: "v: FALSE",
			value:  map[string]interface{}{"v": false},
		},
		{
			source: "v: y",
			value:  map[string]interface{}{"v": "y"}, // y or yes or Yes is string
		},
		{
			source: "v: NO",
			value:  map[string]interface{}{"v": "NO"}, // no or No or NO is string
		},
		{
			source: "v: on",
			value:  map[string]interface{}{"v": "on"}, // on is string
		},

		// Some cross type conversions
		{
			source: "v: 42",
			value:  map[string]uint{"v": 42},
		},
		{
			source: "v: 4294967296",
			value:  map[string]uint64{"v": uint64(4294967296)},
		},

		// int
		{
			source: "v: 2147483647",
			value:  map[string]int{"v": math.MaxInt32},
		},
		{
			source: "v: -2147483648",
			value:  map[string]int{"v": math.MinInt32},
		},

		// int64
		{
			source: "v: 9223372036854775807",
			value:  map[string]int64{"v": math.MaxInt64},
		},
		{
			source: "v: 0b111111111111111111111111111111111111111111111111111111111111111",
			value:  map[string]int64{"v": math.MaxInt64},
		},
		{
			source: "v: -9223372036854775808",
			value:  map[string]int64{"v": math.MinInt64},
		},
		{
			source: "v: -0b111111111111111111111111111111111111111111111111111111111111111",
			value:  map[string]int64{"v": -math.MaxInt64},
		},

		// uint
		{
			source: "v: 0",
			value:  map[string]uint{"v": 0},
		},
		{
			source: "v: 4294967295",
			value:  map[string]uint{"v": math.MaxUint32},
		},
		{
			source: "v: 1e3",
			value:  map[string]uint{"v": 1000},
		},

		// uint64
		{
			source: "v: 0",
			value:  map[string]uint{"v": 0},
		},
		{
			source: "v: 18446744073709551615",
			value:  map[string]uint64{"v": math.MaxUint64},
		},
		{
			source: "v: 0b1111111111111111111111111111111111111111111111111111111111111111",
			value:  map[string]uint64{"v": math.MaxUint64},
		},
		{
			source: "v: 9223372036854775807",
			value:  map[string]uint64{"v": math.MaxInt64},
		},
		{
			source: "v: 1e3",
			value:  map[string]uint64{"v": 1000},
		},

		// float32
		{
			source: "v: 3.40282346638528859811704183484516925440e+38",
			value:  map[string]float32{"v": math.MaxFloat32},
		},
		{
			source: "v: 1.401298464324817070923729583289916131280e-45",
			value:  map[string]float32{"v": math.SmallestNonzeroFloat32},
		},
		{
			source: "v: 18446744073709551615",
			value:  map[string]float32{"v": float32(math.MaxUint64)},
		},
		{
			source: "v: 18446744073709551616",
			value:  map[string]float32{"v": float32(math.MaxUint64 + 1)},
		},
		{
			source: "v: 1e-06",
			value:  map[string]float32{"v": 1e-6},
		},

		// float64
		{
			source: "v: 1.797693134862315708145274237317043567981e+308",
			value:  map[string]float64{"v": math.MaxFloat64},
		},
		{
			source: "v: 4.940656458412465441765687928682213723651e-324",
			value:  map[string]float64{"v": math.SmallestNonzeroFloat64},
		},
		{
			source: "v: 18446744073709551615",
			value:  map[string]float64{"v": float64(math.MaxUint64)},
		},
		{
			source: "v: 18446744073709551616",
			value:  map[string]float64{"v": float64(math.MaxUint64 + 1)},
		},
		{
			source: "v: 1e-06",
			value:  map[string]float64{"v": 1e-06},
		},

		// Timestamps
		{
			// Date only.
			source: "v: 2015-01-01\n",
			value:  map[string]time.Time{"v": time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			// RFC3339
			source: "v: 2015-02-24T18:19:39.12Z\n",
			value:  map[string]time.Time{"v": time.Date(2015, 2, 24, 18, 19, 39, .12e9, time.UTC)},
		},
		{
			// RFC3339 with short dates.
			source: "v: 2015-2-3T3:4:5Z",
			value:  map[string]time.Time{"v": time.Date(2015, 2, 3, 3, 4, 5, 0, time.UTC)},
		},
		{
			// ISO8601 lower case t
			source: "v: 2015-02-24t18:19:39Z\n",
			value:  map[string]time.Time{"v": time.Date(2015, 2, 24, 18, 19, 39, 0, time.UTC)},
		},
		{
			// space separate, no time zone
			source: "v: 2015-02-24 18:19:39\n",
			value:  map[string]time.Time{"v": time.Date(2015, 2, 24, 18, 19, 39, 0, time.UTC)},
		},
		{
			source: "v: 60s\n",
			value:  map[string]time.Duration{"v": time.Minute},
		},
		{
			source: "v: -0.5h\n",
			value:  map[string]time.Duration{"v": -30 * time.Minute},
		},

		// Single Quoted values.
		{
			source: `'1': '2'`,
			value:  map[interface{}]interface{}{"1": `2`},
		},
		{
			source: `'1': '"2"'`,
			value:  map[interface{}]interface{}{"1": `"2"`},
		},
		{
			source: `'1': ''''`,
			value:  map[interface{}]interface{}{"1": `'`},
		},
		{
			source: `'1': '''2'''`,
			value:  map[interface{}]interface{}{"1": `'2'`},
		},
		{
			source: `'1': 'B''z'`,
			value:  map[interface{}]interface{}{"1": `B'z`},
		},
		{
			source: `'1': '\'`,
			value:  map[interface{}]interface{}{"1": `\`},
		},
		{
			source: `'1': '\\'`,
			value:  map[interface{}]interface{}{"1": `\\`},
		},
		{
			source: `'1': '\"2\"'`,
			value:  map[interface{}]interface{}{"1": `\"2\"`},
		},
		{
			source: `'1': '\\"2\\"'`,
			value:  map[interface{}]interface{}{"1": `\\"2\\"`},
		},
		{
			source: "'1': '   1\n    2\n    3'",
			value:  map[interface{}]interface{}{"1": "   1 2 3"},
		},
		{
			source: "'1': '\n    2\n    3'",
			value:  map[interface{}]interface{}{"1": " 2 3"},
		},

		// Double Quoted values.
		{
			source: `"1": "2"`,
			value:  map[interface{}]interface{}{"1": `2`},
		},
		{
			source: `"1": "\"2\""`,
			value:  map[interface{}]interface{}{"1": `"2"`},
		},
		{
			source: `"1": "\""`,
			value:  map[interface{}]interface{}{"1": `"`},
		},
		{
			source: `"1": "X\"z"`,
			value:  map[interface{}]interface{}{"1": `X"z`},
		},
		{
			source: `"1": "\\"`,
			value:  map[interface{}]interface{}{"1": `\`},
		},
		{
			source: `"1": "\\\\"`,
			value:  map[interface{}]interface{}{"1": `\\`},
		},
		{
			source: `"1": "\\\"2\\\""`,
			value:  map[interface{}]interface{}{"1": `\"2\"`},
		},
		{
			source: "'1': \"   1\n    2\n    3\"",
			value:  map[interface{}]interface{}{"1": "   1 2 3"},
		},
		{
			source: "'1': \"\n    2\n    3\"",
			value:  map[interface{}]interface{}{"1": " 2 3"},
		},
		{
			source: `"1": "a\x2Fb"`,
			value:  map[interface{}]interface{}{"1": `a/b`},
		},
		{
			source: `"1": "a\u002Fb"`,
			value:  map[interface{}]interface{}{"1": `a/b`},
		},
		{
			source: `"1": "a\x2Fb\u002Fc\U0000002Fd"`,
			value:  map[interface{}]interface{}{"1": `a/b/c/d`},
		},
		{
			source: "'1': \"2\\n3\"",
			value:  map[interface{}]interface{}{"1": "2\n3"},
		},
		{
			source: "'1': \"2\\r\\n3\"",
			value:  map[interface{}]interface{}{"1": "2\r\n3"},
		},
		{
			source: "'1': \"a\\\nb\\\nc\"",
			value:  map[interface{}]interface{}{"1": "abc"},
		},
		{
			source: "'1': \"a\\\r\nb\\\r\nc\"",
			value:  map[interface{}]interface{}{"1": "abc"},
		},
		{
			source: "'1': \"a\\\rb\\\rc\"",
			value:  map[interface{}]interface{}{"1": "abc"},
		},

		{
			source: "a: -b_c",
			value:  map[string]interface{}{"a": "-b_c"},
		},
		{
			source: "a: +b_c",
			value:  map[string]interface{}{"a": "+b_c"},
		},
		{
			source: "a: 50cent_of_dollar",
			value:  map[string]interface{}{"a": "50cent_of_dollar"},
		},

		// Nulls
		{
			source: "null",
			value:  (*struct{})(nil),
		},
		{
			source: "~",
			value:  (*struct{})(nil),
		},
		{
			source: "v:",
			value:  map[string]interface{}{"v": nil},
		},
		{
			source: "v: ~",
			value:  map[string]interface{}{"v": nil},
		},
		{
			source: "~: null key",
			value:  map[interface{}]string{nil: "null key"},
		},
		{
			source: "v:",
			value:  map[string]*bool{"v": nil},
		},
		{
			source: "v: null",
			value:  map[string]*string{"v": nil},
		},
		{
			source: "v: null",
			value:  map[string]string{"v": ""},
		},
		{
			source: "v: null",
			value:  map[string]interface{}{"v": nil},
		},
		{
			source: "v: Null",
			value:  map[string]interface{}{"v": nil},
		},
		{
			source: "v: NULL",
			value:  map[string]interface{}{"v": nil},
		},
		{
			source: "v: ~",
			value:  map[string]*string{"v": nil},
		},
		{
			source: "v: ~",
			value:  map[string]string{"v": ""},
		},

		{
			source: "v: .inf\n",
			value:  map[string]interface{}{"v": math.Inf(0)},
		},
		{
			source: "v: .Inf\n",
			value:  map[string]interface{}{"v": math.Inf(0)},
		},
		{
			source: "v: .INF\n",
			value:  map[string]interface{}{"v": math.Inf(0)},
		},
		{
			source: "v: -.inf\n",
			value:  map[string]interface{}{"v": math.Inf(-1)},
		},
		{
			source: "v: -.Inf\n",
			value:  map[string]interface{}{"v": math.Inf(-1)},
		},
		{
			source: "v: -.INF\n",
			value:  map[string]interface{}{"v": math.Inf(-1)},
		},
		{
			source: "v: .nan\n",
			value:  map[string]interface{}{"v": math.NaN()},
		},
		{
			source: "v: .NaN\n",
			value:  map[string]interface{}{"v": math.NaN()},
		},
		{
			source: "v: .NAN\n",
			value:  map[string]interface{}{"v": math.NaN()},
		},

		// Explicit tags.
		{
			source: "v: !!float '1.1'",
			value:  map[string]interface{}{"v": 1.1},
		},
		{
			source: "v: !!float 0",
			value:  map[string]interface{}{"v": float64(0)},
		},
		{
			source: "v: !!float -1",
			value:  map[string]interface{}{"v": float64(-1)},
		},
		{
			source: "v: !!null ''",
			value:  map[string]interface{}{"v": nil},
		},
		{
			source: "v: !!timestamp \"2015-01-01\"",
			value:  map[string]time.Time{"v": time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			source: "v: !!timestamp 2015-01-01",
			value:  map[string]time.Time{"v": time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			source: "v: !!bool yes",
			value:  map[string]bool{"v": true},
		},
		{
			source: "v: !!bool False",
			value:  map[string]bool{"v": false},
		},
		{
			source: `
!!merge <<: { a: 1, b: 2 }
c: 3
`,
			value: map[string]any{"a": 1, "b": 2, "c": 3},
		},

		// merge
		{
			source: `
a: &a
 foo: 1
b: &b
 bar: 2
merge:
 <<: [*a, *b]
`,
			value: map[string]map[string]any{
				"a":     {"foo": 1},
				"b":     {"bar": 2},
				"merge": {"foo": 1, "bar": 2},
			},
		},
		{
			source: `
a: &a
 foo: 1
b: &b
 bar: 2
merge:
 <<: [*a, *b]
`,
			value: map[string]yaml.MapSlice{
				"a":     {{Key: "foo", Value: 1}},
				"b":     {{Key: "bar", Value: 2}},
				"merge": {{Key: "foo", Value: 1}, {Key: "bar", Value: 2}},
			},
		},

		// Flow sequence
		{
			source: "v: [A,B]",
			value:  map[string]interface{}{"v": []interface{}{"A", "B"}},
		},
		{
			source: "v: [A,B,C,]",
			value:  map[string][]string{"v": {"A", "B", "C"}},
		},
		{
			source: "v: [A,1,C]",
			value:  map[string][]string{"v": {"A", "1", "C"}},
		},
		{
			source: "v: [A,1,C]",
			value:  map[string]interface{}{"v": []interface{}{"A", 1, "C"}},
		},
		{
			source: "v: [a: b, c: d]",
			value: map[string]any{"v": []any{
				map[string]any{"a": "b"},
				map[string]any{"c": "d"},
			}},
		},
		{
			source: "v: [{a: b}, {c: d, e: f}]",
			value: map[string]any{"v": []any{
				map[string]any{"a": "b"},
				map[string]any{
					"c": "d",
					"e": "f",
				},
			}},
		},

		// Block sequence
		{
			source: "v:\n - A\n - B",
			value:  map[string]interface{}{"v": []interface{}{"A", "B"}},
		},
		{
			source: "v:\n - A\n - B\n - C",
			value:  map[string][]string{"v": {"A", "B", "C"}},
		},
		{
			source: "v:\n - A\n - 1\n - C",
			value:  map[string][]string{"v": {"A", "1", "C"}},
		},
		{
			source: "v:\n - A\n - 1\n - C",
			value:  map[string]interface{}{"v": []interface{}{"A", 1, "C"}},
		},

		// Map inside interface with no type hints.
		{
			source: "a: {b: c}",
			value:  map[interface{}]interface{}{"a": map[interface{}]interface{}{"b": "c"}},
		},

		{
			source: "v: \"\"\n",
			value:  map[string]string{"v": ""},
		},
		{
			source: "v:\n- A\n- B\n",
			value:  map[string][]string{"v": {"A", "B"}},
		},
		{
			source: "a: '-'\n",
			value:  map[string]string{"a": "-"},
		},
		{
			source: "123\n",
			value:  123,
		},
		{
			source: "hello: world\n",
			value:  map[string]string{"hello": "world"},
		},
		{
			source: "hello: world\r\n",
			value:  map[string]string{"hello": "world"},
		},
		{
			source: "hello: world\rGo: Gopher",
			value:  map[string]string{"hello": "world", "Go": "Gopher"},
		},

		// Structs and type conversions.
		{
			source: "hello: world",
			value:  struct{ Hello string }{"world"},
		},
		{
			source: "a: {b: c}",
			value:  struct{ A struct{ B string } }{struct{ B string }{"c"}},
		},
		{
			source: "a: {b: c}",
			value:  struct{ A map[string]string }{map[string]string{"b": "c"}},
		},
		{
			source: "a:",
			value:  struct{ A map[string]string }{},
		},
		{
			source: "a: 1",
			value:  struct{ A int }{1},
		},
		{
			source: "a: 1",
			value:  struct{ A float64 }{1},
		},
		{
			source: "a: 1.0",
			value:  struct{ A int }{1},
		},
		{
			source: "a: 1.0",
			value:  struct{ A uint }{1},
		},
		{
			source: "a: [1, 2]",
			value:  struct{ A []int }{[]int{1, 2}},
		},
		{
			source: "a: [1, 2]",
			value:  struct{ A [2]int }{[2]int{1, 2}},
		},
		{
			source: "a: 1",
			value:  struct{ B int }{0},
		},
		{
			source: "a: 1",
			value: struct {
				B int `yaml:"a"`
			}{1},
		},

		{
			source: "a: 1\n",
			value:  yaml.MapItem{Key: "a", Value: 1},
		},
		{
			source: "a: 1\nb: 2\nc: 3\n",
			value: yaml.MapSlice{
				{Key: "a", Value: 1},
				{Key: "b", Value: 2},
				{Key: "c", Value: 3},
			},
		},
		{
			source: "v:\n- A\n- 1\n- B:\n  - 2\n  - 3\n",
			value: map[string]interface{}{
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
			source: "a:\n  b: c\n",
			value: map[string]interface{}{
				"a": map[string]string{
					"b": "c",
				},
			},
		},
		{
			source: "a: {x: 1}\n",
			value: map[string]map[string]int{
				"a": {
					"x": 1,
				},
			},
		},
		{
			source: "t2: 2018-01-09T10:40:47Z\nt4: 2098-01-09T10:40:47Z\n",
			value: map[string]string{
				"t2": "2018-01-09T10:40:47Z",
				"t4": "2098-01-09T10:40:47Z",
			},
		},
		{
			source: "a: [1, 2]\n",
			value: map[string][]int{
				"a": {1, 2},
			},
		},
		{
			source: "a: {b: c, d: e}\n",
			value: map[string]interface{}{
				"a": map[string]string{
					"b": "c",
					"d": "e",
				},
			},
		},
		{
			source: "a: 3s\n",
			value: map[string]string{
				"a": "3s",
			},
		},
		{
			source: "a: <foo>\n",
			value:  map[string]string{"a": "<foo>"},
		},
		{
			source: "a: \"1:1\"\n",
			value:  map[string]string{"a": "1:1"},
		},
		{
			source: "a: 1.2.3.4\n",
			value:  map[string]string{"a": "1.2.3.4"},
		},
		{
			source: "a: 'b: c'\n",
			value:  map[string]string{"a": "b: c"},
		},
		{
			source: "a: 'Hello #comment'\n",
			value:  map[string]string{"a": "Hello #comment"},
		},
		{
			source: "a: 100.5\n",
			value: map[string]interface{}{
				"a": 100.5,
			},
		},
		{
			source: "a: \"\\0\"\n",
			value:  map[string]string{"a": "\x00"},
		},
		{
			source: "b: 2\na: 1\nd: 4\nc: 3\nsub:\n  e: 5\n",
			value: map[string]interface{}{
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
			source: "       a       :          b        \n",
			value:  map[string]string{"a": "b"},
		},
		{
			source: "a: b # comment\nb: c\n",
			value: map[string]string{
				"a": "b",
				"b": "c",
			},
		},
		{
			source: "---\na: b\n",
			value:  map[string]string{"a": "b"},
		},
		{
			source: "a: b\n...\n",
			value:  map[string]string{"a": "b"},
		},
		{
			source: "%YAML 1.2\n---\n",
			value:  (*struct{})(nil),
			eof:    true,
		},
		{
			source: "---\n",
			value:  (*struct{})(nil),
			eof:    true,
		},
		{
			source: "...",
			value:  (*struct{})(nil),
			eof:    true,
		},
		{
			source: "v: go test ./...",
			value:  map[string]string{"v": "go test ./..."},
		},
		{
			source: "v: echo ---",
			value:  map[string]string{"v": "echo ---"},
		},
		{
			source: "v: |\n  hello\n  ...\n  world\n",
			value:  map[string]string{"v": "hello\n...\nworld\n"},
		},
		{
			source: "a: !!binary gIGC\n",
			value:  map[string]string{"a": "\x80\x81\x82"},
		},
		{
			source: "a: !!binary |\n  " + strings.Repeat("kJCQ", 17) + "kJ\n  CQ\n",
			value:  map[string]string{"a": strings.Repeat("\x90", 54)},
		},
		{
			source: "v:\n- A\n- |-\n  B\n  C\n",
			value: map[string][]string{
				"v": {
					"A", "B\nC",
				},
			},
		},
		{
			source: "v:\n- A\n- |-\n  B\n  C\n\n\n",
			value: map[string][]string{
				"v": {
					"A", "B\nC",
				},
			},
		},
		{
			source: "v:\n- A\n- >-\n  B\n  C\n",
			value: map[string][]string{
				"v": {
					"A", "B C",
				},
			},
		},
		{
			source: "v:\n- A\n- >-\n  B\n  C\n\n\n",
			value: map[string][]string{
				"v": {
					"A", "B C",
				},
			},
		},
		{
			source: "a: b\nc: d\n",
			value: struct {
				A string
				C string `yaml:"c"`
			}{
				"b", "d",
			},
		},
		{
			source: "a: 1\nb: 2\n",
			value: struct {
				A int
				B int `yaml:"-"`
			}{
				1, 0,
			},
		},
		{
			source: "a: 1\nb: 2\n",
			value: struct {
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
			source: "a: &x 1\nb: &y 2\nc: *x\nd: *y\n",
			value:  struct{ A, B, C, D int }{1, 2, 1, 2},
		},
		{
			source: "a: &a {c: 1}\nb: *a\n",
			value: struct {
				A, B struct {
					C int
				}
			}{struct{ C int }{1}, struct{ C int }{1}},
		},
		{
			source: "a: &a [1, 2]\nb: *a\n",
			value:  struct{ B []int }{[]int{1, 2}},
		},
		{
			source: "key1: &anchor\n  subkey: *anchor\nkey2: *anchor\n",
			value: map[string]any{
				"key1": map[string]any{
					"subkey": nil,
				},
				"key2": map[string]any{
					"subkey": nil,
				},
			},
		},
		{
			source: `{a: &a c, *a : b}`,
			value:  map[string]string{"a": "c", "c": "b"},
		},
		{
			source: "tags:\n- hello-world\na: foo",
			value: struct {
				Tags []string
				A    string
			}{Tags: []string{"hello-world"}, A: "foo"},
		},
		{
			source: "",
			value:  (*struct{})(nil),
			eof:    true,
		},
		{
			source: "{}",
			value:  struct{}{},
		},
		{
			source: "{a: , b: c}",
			value:  map[string]any{"a": nil, "b": "c"},
		},
		{
			source: "v: /a/{b}",
			value:  map[string]string{"v": "/a/{b}"},
		},
		{
			source: "v: 1[]{},!%?&*",
			value:  map[string]string{"v": "1[]{},!%?&*"},
		},
		{
			source: "v: user's item",
			value:  map[string]string{"v": "user's item"},
		},
		{
			source: "v: [1,[2,[3,[4,5],6],7],8]",
			value: map[string]interface{}{
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
			source: "v: {a: {b: {c: {d: e},f: g},h: i},j: k}",
			value: map[string]interface{}{
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
			source: `---
- a:
    b:
- c: d
`,
			value: []map[string]interface{}{
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
			source: `---
a:
  b:
c: d
`,
			value: map[string]interface{}{
				"a": map[string]interface{}{
					"b": nil,
				},
				"c": "d",
			},
		},
		{
			source: `---
a:
b:
c:
`,
			value: map[string]interface{}{
				"a": nil,
				"b": nil,
				"c": nil,
			},
		},
		{
			source: `---
a: go test ./...
b:
c:
`,
			value: map[string]interface{}{
				"a": "go test ./...",
				"b": nil,
				"c": nil,
			},
		},
		{
			source: `---
a: |
  hello
  ...
  world
b:
c:
`,
			value: map[string]interface{}{
				"a": "hello\n...\nworld\n",
				"b": nil,
				"c": nil,
			},
		},

		// Multi bytes
		{
			source: "v: あいうえお\nv2: かきくけこ",
			value:  map[string]string{"v": "あいうえお", "v2": "かきくけこ"},
		},
		{
			source: `
- "Fun with \\"
- "\" \a \b \e \f"
- "\n \r \t \v \0"
- "\  \_ \N \L \P \
  \x41 \u0041 \U00000041"
`,
			value: []string{"Fun with \\", "\" \u0007 \b \u001b \f", "\n \r \t \u000b \u0000", "\u0020 \u00a0 \u0085 \u2028 \u2029 A A A"},
		},
		{
			source: `"\ud83e\udd23"`,
			value:  "🤣",
		},
		{
			source: `"\uD83D\uDE00\uD83D\uDE01"`,
			value:  "😀😁",
		},
		{
			source: `"\uD83D\uDE00a\uD83D\uDE01"`,
			value:  "😀a😁",
		},
	}
	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			buf := bytes.NewBufferString(test.source)
			dec := yaml.NewDecoder(buf)
			typ := reflect.ValueOf(test.value).Type()
			value := reflect.New(typ)
			if err := dec.Decode(value.Interface()); err != nil {
				if test.eof && err == io.EOF {
					return
				}
				t.Fatalf("%s: %+v", test.source, err)
			}
			if test.eof {
				t.Fatal("expected EOF but got no error")
			}
			actual := fmt.Sprintf("%+v", value.Elem().Interface())
			expect := fmt.Sprintf("%+v", test.value)
			if actual != expect {
				t.Fatalf("failed to test [%s], actual=[%s], expect=[%s]", test.source, actual, expect)
			}
		})
	}
}

func TestDecoder_Invalid(t *testing.T) {
	tests := []struct {
		src    string
		expect string
	}{
		{
			"*-0",
			`
[1:2] could not find alias "-0"
>  1 | *-0
        ^
`,
		},
	}
	for _, test := range tests {
		t.Run(test.src, func(t *testing.T) {
			var v any
			err := yaml.Unmarshal([]byte(test.src), &v)
			if err == nil {
				t.Fatal("cannot catch decode error")
			}
			actual := "\n" + err.Error()
			if test.expect != actual {
				t.Fatalf("expected: [%s] but got [%s]", test.expect, actual)
			}
		})
	}
}

func TestDecoder_ScientificNotation(t *testing.T) {
	tests := []struct {
		source string
		value  interface{}
	}{
		{
			"v: 1e3",
			map[string]uint{"v": 1000},
		},
		{
			"v: 1e-3",
			map[string]uint{"v": 0},
		},
		{
			"v: 1e3",
			map[string]int{"v": 1000},
		},
		{
			"v: 1e-3",
			map[string]int{"v": 0},
		},
		{
			"v: 1e3",
			map[string]float32{"v": 1000},
		},
		{
			"v: 1.0e3",
			map[string]float64{"v": 1000},
		},
		{
			"v: 1e-3",
			map[string]float64{"v": 0.001},
		},
		{
			"v: 1.0e-3",
			map[string]float64{"v": 0.001},
		},
		{
			"v: 1.0e+3",
			map[string]float64{"v": 1000},
		},
		{
			"v: 1.0e+3",
			map[string]float64{"v": 1000},
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
			if !strings.Contains(err.Error(), msg) {
				t.Fatalf("expected error message: %s to contain: %s", err.Error(), msg)
			}
		})
		t.Run("string to uint", func(t *testing.T) {
			var v T
			err := yaml.Unmarshal([]byte(`b: str`), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := "cannot unmarshal string into Go struct field T.B of type uint"
			if !strings.Contains(err.Error(), msg) {
				t.Fatalf("expected error message: %s to contain: %s", err.Error(), msg)
			}
		})
		t.Run("string to bool", func(t *testing.T) {
			var v T
			err := yaml.Unmarshal([]byte(`d: str`), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := "cannot unmarshal string into Go struct field T.D of type bool"
			if !strings.Contains(err.Error(), msg) {
				t.Fatalf("expected error message: %s to contain: %s", err.Error(), msg)
			}
		})
		t.Run("string to int at inline", func(t *testing.T) {
			var v U
			err := yaml.Unmarshal([]byte(`a: str`), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := "cannot unmarshal string into Go struct field U.T.A of type int"
			if !strings.Contains(err.Error(), msg) {
				t.Fatalf("expected error message: %s to contain: %s", err.Error(), msg)
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
			if !strings.Contains(err.Error(), msg) {
				t.Fatalf("expected error message: %s to contain: %s", err.Error(), msg)
			}
		})
		t.Run("string to int", func(t *testing.T) {
			var v map[string][]int
			err := yaml.Unmarshal([]byte("v:\n - A\n - 1\n - C"), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := "cannot unmarshal string into Go value of type int"
			if !strings.Contains(err.Error(), msg) {
				t.Fatalf("expected error message: %s to contain: %s", err.Error(), msg)
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
			if !strings.Contains(err.Error(), msg) {
				t.Fatalf("expected error message: %s to contain: %s", err.Error(), msg)
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
			if !strings.Contains(err.Error(), msg) {
				t.Fatalf("expected error message: %s to contain: %s", err.Error(), msg)
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
			if !strings.Contains(err.Error(), msg) {
				t.Fatalf("expected error message: %s to contain: %s", err.Error(), msg)
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
			if !strings.Contains(err.Error(), msg) {
				t.Fatalf("expected error message: %s to contain: %s", err.Error(), msg)
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
			if !strings.Contains(err.Error(), msg) {
				t.Fatalf("expected error message: %s to contain: %s", err.Error(), msg)
			}
		})
		t.Run("string to duration", func(t *testing.T) {
			var v T
			err := yaml.Unmarshal([]byte(`b: str`), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := `time: invalid duration "str"`
			if !strings.Contains(err.Error(), msg) {
				t.Fatalf("expected error message: %s to contain: %s", err.Error(), msg)
			}
		})
		t.Run("int to duration", func(t *testing.T) {
			var v T
			err := yaml.Unmarshal([]byte(`b: 10`), &v)
			if err == nil {
				t.Fatal("expected to error")
			}
			msg := "cannot unmarshal uint64 into Go struct field T.B of type time.Duration"
			if !strings.Contains(err.Error(), msg) {
				t.Fatalf("expected error message: %s to contain: %s", err.Error(), msg)
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
	dec := yaml.NewDecoder(buf, yaml.AllowDuplicateMapKey())
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
		items, _ := v.(map[string]interface{})["items"].([]interface{})
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

	// TODO: properly check if errors are colored/have source
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

	// TODO: properly check if errors are colored/have source
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

func TestDecoder_AllowDuplicateMapKey(t *testing.T) {
	yml := `
a: b
a: c
`
	t.Run("map", func(t *testing.T) {
		var v map[string]string
		if err := yaml.NewDecoder(strings.NewReader(yml), yaml.AllowDuplicateMapKey()).Decode(&v); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("struct", func(t *testing.T) {
		var v struct {
			A string
		}
		if err := yaml.NewDecoder(strings.NewReader(yml), yaml.AllowDuplicateMapKey()).Decode(&v); err != nil {
			t.Fatal(err)
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

func ExampleUnmarshal_yAMLTags() {
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

func TestDecoder_CustomUnmarshaler(t *testing.T) {
	t.Run("override struct type", func(t *testing.T) {
		type T struct {
			Foo string `yaml:"foo"`
		}
		src := []byte(`foo: "bar"`)
		var v T
		if err := yaml.UnmarshalWithOptions(src, &v, yaml.CustomUnmarshaler[T](func(dst *T, b []byte) error {
			if !bytes.Equal(src, b) {
				t.Fatalf("failed to get decode target buffer. expected %q but got %q", src, b)
			}
			var v T
			if err := yaml.Unmarshal(b, &v); err != nil {
				return err
			}
			if v.Foo != "bar" {
				t.Fatal("failed to decode")
			}
			dst.Foo = "bazbaz" // assign another value to target
			return nil
		})); err != nil {
			t.Fatal(err)
		}
		if v.Foo != "bazbaz" {
			t.Fatalf("failed to switch to custom unmarshaler. got: %v", v.Foo)
		}
	})
	t.Run("override bytes type", func(t *testing.T) {
		type T struct {
			Foo []byte `yaml:"foo"`
		}
		src := []byte(`foo: "bar"`)
		var v T
		if err := yaml.UnmarshalWithOptions(src, &v, yaml.CustomUnmarshaler[[]byte](func(dst *[]byte, b []byte) error {
			if !bytes.Equal(b, []byte(`"bar"`)) {
				t.Fatalf("failed to get target buffer: %q", b)
			}
			*dst = []byte("bazbaz")
			return nil
		})); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(v.Foo, []byte("bazbaz")) {
			t.Fatalf("failed to switch to custom unmarshaler. got: %q", v.Foo)
		}
	})
	t.Run("override bytes type with context", func(t *testing.T) {
		type T struct {
			Foo []byte `yaml:"foo"`
		}
		src := []byte(`foo: "bar"`)
		var v T
		ctx := context.WithValue(context.Background(), "plop", uint(42))
		if err := yaml.UnmarshalContext(ctx, src, &v, yaml.CustomUnmarshalerContext[[]byte](func(ctx context.Context, dst *[]byte, b []byte) error {
			if !bytes.Equal(b, []byte(`"bar"`)) {
				t.Fatalf("failed to get target buffer: %q", b)
			}
			if ctx.Value("plop") != uint(42) {
				t.Fatalf("context value is not correct")
			}
			*dst = []byte("bazbaz")
			return nil
		})); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(v.Foo, []byte("bazbaz")) {
			t.Fatalf("failed to switch to custom unmarshaler. got: %q", v.Foo)
		}
	})
}

type unmarshalContext struct {
	v int
}

func (c *unmarshalContext) UnmarshalYAML(ctx context.Context, b []byte) error {
	v, ok := ctx.Value("k").(int)
	if !ok {
		return errors.New("cannot get valid context")
	}
	if v != 1 {
		return errors.New("cannot get valid context")
	}
	if string(b) != "1" {
		return errors.New("cannot get valid bytes")
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
		if !errors.Is(err, yaml.ErrDecodeRequiredPointerType) {
			t.Fatalf("unexpected error: %s", err)
		}
	})
}

func TestCommentWithCustomUnmarshaler(t *testing.T) {
	type T struct{}

	for idx, test := range []string{
		`
foo:
  # comment
  - a: b
`,
		`
foo: # comment
  bar: 1
  baz: true
`,
	} {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			m := yaml.CommentMap{}
			var v T
			if err := yaml.UnmarshalWithOptions(
				[]byte(test),
				&v,
				yaml.CommentToMap(m),
				yaml.CustomUnmarshaler[T](func(dst *T, b []byte) error {
					expected := bytes.Trim([]byte(test), "\n")
					if !bytes.Equal(b, expected) {
						return fmt.Errorf("failed to decode: got\n%s", string(test))
					}
					return nil
				}),
			); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func ExampleUnmarshal_jSONTags() {
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

func ExampleDecoder_Decode_disallowUnknownField() {
	var v struct {
		A string `yaml:"simple"`
		C string `yaml:"complicated"`
	}

	const src = `---
simple: string
unknown: string
`
	err := yaml.NewDecoder(strings.NewReader(src), yaml.DisallowUnknownField()).Decode(&v)
	fmt.Printf("%v\n", err)

	// OUTPUT:
	// [3:1] unknown field "unknown"
	//    1 | ---
	//    2 | simple: string
	// >  3 | unknown: string
	//        ^
}

func ExampleNodeToValue() {
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
	t.Run("quoted map keys", func(t *testing.T) {
		t.Parallel()
		yml := `
a:
  "b"  : 2
  'c': true
`
		var v struct {
			A struct {
				B int
				C bool
			}
		}
		if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
			t.Fatalf("failed to unmarshal %v", err)
		}
		if v.A.B != 2 {
			t.Fatalf("expected a.b to equal 2 but was %d", v.A.B)
		}
		if !v.A.C {
			t.Fatal("expected a.c to be true but was false")
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
	type value struct {
		String unmarshalYAMLWithAliasString
		Map    unmarshalYAMLWithAliasMap
	}
	tests := []struct {
		name          string
		yaml          string
		expectedValue value
		err           string
	}{
		{
			name: "ok",
			yaml: `
anchors:
 w: &w "\"hello\" \"world\""
 map: &x
   a: b
   c: d
   d: *w
string: *w
map:
 <<: *x
 e: f
`,
			expectedValue: value{
				String: unmarshalYAMLWithAliasString(`"hello" "world"`),
				Map: unmarshalYAMLWithAliasMap(map[string]interface{}{
					"a": "b",
					"c": "d",
					"d": `"hello" "world"`,
					"e": "f",
				}),
			},
		},
		{
			name: "unknown alias",
			yaml: `
anchors:
 w: &w "\"hello\" \"world\""
 map: &x
   a: b
   c: d
   d: *w
string: *y
map:
 <<: *z
 e: f
`,
			err: `could not find alias "y"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var v value
			err := yaml.Unmarshal([]byte(test.yaml), &v)

			if test.err != "" {
				if err == nil {
					t.Fatal("expected to error")
				}
				if !strings.Contains(err.Error(), test.err) {
					t.Fatalf("expected error message: %s to contain: %s", err.Error(), test.err)
				}
			} else {
				if err != nil {
					t.Fatalf("%+v", err)
				}
				if !reflect.DeepEqual(test.expectedValue, v) {
					t.Fatalf("non matching values:\nexpected[%s]\ngot     [%s]", test.expectedValue, v)
				}
			}
		})
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
- b: c # comment
  # comment
  d: | # comment
    hello

    hello
  f: g
- h: i`
	actual := "\n" + string(b)
	if expected != actual {
		return fmt.Errorf("unexpected bytes: expected [%q] but got [%q]", expected, actual)
	}
	var v []map[string]unmarshalString
	if err := yaml.Unmarshal(b, &v); err != nil {
		return err
	}
	u.v = v
	return nil
}

func TestDecoder_DecodeWithAnchorAnyValue(t *testing.T) {
	type Config struct {
		Env []string `json:"env"`
	}

	type Schema struct {
		Def    map[string]any `json:"def"`
		Config Config         `json:"config"`
	}

	data := `
def:
  myenv: &my_env
    - VAR1=1
    - VAR2=2
config:
  env: *my_env
`

	var cfg Schema
	if err := yaml.NewDecoder(strings.NewReader(data)).Decode(&cfg); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg.Config.Env, []string{"VAR1=1", "VAR2=2"}) {
		t.Fatalf("failed to decode value. actual = %+v", cfg)
	}
}

func TestDecoder_UnmarshalBytesWithSeparatedList(t *testing.T) {
	yml := `
a:
 - b: c # comment
   # comment
   d: | # comment
     hello

     hello
   f: g
 - h: i
`
	var v struct {
		A unmarshalList
	}
	cm := yaml.CommentMap{}
	if err := yaml.UnmarshalWithOptions([]byte(yml), &v, yaml.CommentToMap(cm)); err != nil {
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
		{
			Node: "hello\nworld",
		},
		{
			Node: "hello\nworld\n",
		},
		{
			LastNode: "hello\nworld",
		},
		{
			LastNode: "hello\nworld\n",
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

func TestDecodeWithSameAnchor(t *testing.T) {
	yml := `
a: &a 1
b: &a 2
c: &a 3
d: *a
`
	type T struct {
		A int
		B int
		C int
		D int
	}
	var v T
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v, T{A: 1, B: 2, C: 3, D: 3}) {
		t.Fatalf("failed to decode same anchor: %+v", v)
	}
}

func TestUnmarshalMapSliceParallel(t *testing.T) {
	content := `
steps:
  req0:
    desc: Get /users/1
    req:
      /users/1:
        get: nil
    test: |
      current.res.status == 200
  req1:
    desc: Get /private
    req:
      /private:
        get: nil
    test: |
      current.res.status == 403
  req2:
    desc: Get /users
    req:
      /users:
        get: nil
    test: |
      current.res.status == 200
`
	type mappedSteps struct {
		Steps yaml.MapSlice `yaml:"steps,omitempty"`
	}
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("i=%d", i), func(t *testing.T) {
			t.Parallel()
			for i := 0; i < 10; i++ {
				m := mappedSteps{
					Steps: yaml.MapSlice{},
				}
				if err := yaml.Unmarshal([]byte(content), &m); err != nil {
					t.Fatal(err)
				}
				for _, s := range m.Steps {
					_, ok := s.Value.(map[string]interface{})
					if !ok {
						t.Fatal("unexpected error")
					}
				}
			}
		})
	}
}

func TestSameNameInineStruct(t *testing.T) {
	type X struct {
		X float64 `yaml:"x"`
	}

	type T struct {
		X `yaml:",inline"`
	}

	var v T
	if err := yaml.Unmarshal([]byte(`x: 0.7`), &v); err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(v.X.X) != "0.7" {
		t.Fatalf("failed to decode")
	}
}

type unmarshableMapKey struct {
	Key string
}

func (mk *unmarshableMapKey) UnmarshalYAML(b []byte) error {
	mk.Key = string(b)
	return nil
}

type testNodeUnmarshalerCtx struct {
	outErr   error
	received ast.Node
}

func (u *testNodeUnmarshalerCtx) UnmarshalYAML(ctx context.Context, node ast.Node) error {
	if u.outErr != nil {
		return u.outErr
	}

	if ctx == nil {
		return errors.New("nil context")
	}

	u.received = node
	return nil
}

func TestNodeUnmarshalerContext(t *testing.T) {
	type testNodeUnmarshalerBody struct {
		Root testNodeUnmarshalerCtx `yaml:"root"`
	}

	cases := []struct {
		name      string
		expectErr string
		src       []string
		body      testNodeUnmarshalerBody
	}{
		{
			name: "should pass node",
			src: []string{
				"root:",
				"  foo: bar",
				"  fizz: buzz",
			},
		},
		{
			name: "should pass returned error",
			body: testNodeUnmarshalerBody{
				Root: testNodeUnmarshalerCtx{
					outErr: errors.New("test error"),
				},
			},
			expectErr: "test error",
			src: []string{
				"root:",
				"  foo: bar",
				"  fizz: buzz",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			src := []byte(strings.Join(c.src, "\n"))
			out := c.body
			err := yaml.Unmarshal(src, &out)
			if c.expectErr != "" {
				if err == nil {
					t.Fatal("expected error but got nil")
					return
				}

				if !strings.Contains(err.Error(), c.expectErr) {
					t.Fatalf("error message %q should contain %q", err.Error(), c.expectErr)
				}
				return
			}

			expect := struct {
				Root ast.Node `yaml:"root"`
			}{}
			if err := yaml.UnmarshalContext(context.TODO(), src, &expect); err != nil {
				t.Fatal("invalid test yaml:", err)
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !reflect.DeepEqual(out.Root.received, expect.Root) {
				t.Fatalf("expected:\n%#v\n but got:\n%#v", expect.Root, out.Root.received)
			}
		})
	}

}

type testNodeUnmarshaler struct {
	outErr   error
	received ast.Node
}

func (u *testNodeUnmarshaler) UnmarshalYAML(node ast.Node) error {
	if u.outErr != nil {
		return u.outErr
	}

	u.received = node
	return nil
}

func TestNodeUnmarshaler(t *testing.T) {
	type testNodeUnmarshalerBody struct {
		Root testNodeUnmarshaler `yaml:"root"`
	}

	cases := []struct {
		name      string
		expectErr string
		src       []string
		body      testNodeUnmarshalerBody
	}{
		{
			name: "should pass node",
			src: []string{
				"root:",
				"  foo: bar",
				"  fizz: buzz",
			},
		},
		{
			name: "should pass returned error",
			body: testNodeUnmarshalerBody{
				Root: testNodeUnmarshaler{
					outErr: errors.New("test error"),
				},
			},
			expectErr: "test error",
			src: []string{
				"root:",
				"  foo: bar",
				"  fizz: buzz",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			src := []byte(strings.Join(c.src, "\n"))
			out := c.body
			err := yaml.Unmarshal(src, &out)
			if c.expectErr != "" {
				if err == nil {
					t.Fatal("expected error but got nil")
					return
				}

				if !strings.Contains(err.Error(), c.expectErr) {
					t.Fatalf("error message %q should contain %q", err.Error(), c.expectErr)
				}
				return
			}

			expect := struct {
				Root ast.Node `yaml:"root"`
			}{}
			if err := yaml.Unmarshal(src, &expect); err != nil {
				t.Fatal("invalid test yaml:", err)
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !reflect.DeepEqual(out.Root.received, expect.Root) {
				t.Fatalf("expected:\n%#v\n but got:\n%#v", expect.Root, out.Root.received)
			}
		})
	}

}

func TestMapKeyCustomUnmarshaler(t *testing.T) {
	var m map[unmarshableMapKey]string
	if err := yaml.Unmarshal([]byte(`key: value`), &m); err != nil {
		t.Fatalf("failed to unmarshal %v", err)
	}
	if len(m) != 1 {
		t.Fatalf("expected 1 element in map, but got %d", len(m))
	}
	val, ok := m[unmarshableMapKey{Key: "key"}]
	if !ok {
		t.Fatal("expected to have element 'key' in map")
	}
	if val != "value" {
		t.Fatalf("expected to have value \"value\", but got %q", val)
	}
}

type bytesUnmershalerWithMapAlias struct{}

func (*bytesUnmershalerWithMapAlias) UnmarshalYAML(b []byte) error {
	expected := strings.TrimPrefix(`
aaaaa:
  bbbbb:
    bar:
      - |
        foo
          bar
      - name: |
          foo
            bar

`, "\n")
	if string(b) != expected {
		return fmt.Errorf("failed to decode: expected:\n[%s]\nbut got:\n[%s]\n", expected, string(b))
	}
	return nil
}

func TestBytesUnmarshalerWithMapAlias(t *testing.T) {
	yml := `
x-foo: &data
  bar:
    - |
      foo
        bar
    - name: |
        foo
          bar

foo:
  aaaaa:
    bbbbb: *data
`
	type T struct {
		Foo bytesUnmershalerWithMapAlias `yaml:"foo"`
	}
	var v T
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		t.Fatal(err)
	}
}

func TestBytesUnmarshalerWithEmptyValue(t *testing.T) {
	type T struct{}

	unmarshaler := func(dst *T, b []byte) error {
		var v any
		return yaml.Unmarshal(b, &v)
	}

	yml := `
map: &m {}
seq: &seq []
foo: # comment
bar: *m
baz: *seq
`
	m := yaml.CommentMap{}
	var v T
	if err := yaml.UnmarshalWithOptions(
		[]byte(yml),
		&v,
		yaml.CommentToMap(m),
		yaml.CustomUnmarshaler[T](unmarshaler),
	); err != nil {
		t.Fatal(err)
	}
	if err := yaml.UnmarshalWithOptions(
		[]byte(yml),
		&v,
		yaml.CustomUnmarshaler[T](unmarshaler),
	); err != nil {
		t.Fatal(err)
	}

}

func TestIssue650(t *testing.T) {
	type Disk struct {
		Name   string `yaml:"name"`
		Format *bool  `yaml:"format"`
	}

	type Sample struct {
		Disks []Disk `yaml:"disks"`
	}

	unmarshalDisk := func(dst *Disk, b []byte) error {
		var s string
		if err := yaml.Unmarshal(b, &s); err == nil {
			*dst = Disk{Name: s}
			return nil
		}
		return yaml.Unmarshal(b, dst)
	}

	data := []byte(`
disks:
    -      name: foo
           format: true
`)

	var sample Sample
	if err := yaml.UnmarshalWithOptions(data, &sample, yaml.CustomUnmarshaler[Disk](unmarshalDisk)); err != nil {
		t.Fatal(err)
	}
}

func TestBytesUnmarshalerWithLiteral(t *testing.T) {
	t.Run("map value", func(t *testing.T) {
		type Literal string

		unmarshalLit := func(dst *Literal, b []byte) error {
			var s string
			if err := yaml.Unmarshal(b, &s); err != nil {
				return err
			}
			*dst = Literal(s)
			return nil
		}

		data := []byte(`
-         name:  |
           foo
             bar
-         name:
           |
           foo
           bar
`)

		var v []map[string]Literal
		if err := yaml.UnmarshalWithOptions(data, &v, yaml.CustomUnmarshaler[Literal](unmarshalLit)); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(v, []map[string]Literal{{"name": "foo\n  bar\n"}, {"name": "foo\nbar\n"}}) {
			t.Fatalf("failed to get decoded value. got: %q", v)
		}
	})
	t.Run("sequence value", func(t *testing.T) {
		type Literal string

		unmarshalLit := func(dst *Literal, b []byte) error {
			var s string
			if err := yaml.Unmarshal(b, &s); err != nil {
				return err
			}
			*dst = Literal(s)
			return nil
		}

		data := []byte(`
-            |
  foo
    bar
-
 |
 foo
 bar
`)

		var v []Literal
		if err := yaml.UnmarshalWithOptions(data, &v, yaml.CustomUnmarshaler[Literal](unmarshalLit)); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(v, []Literal{"foo\n  bar\n", "foo\nbar\n"}) {
			t.Fatalf("failed to get decoded value. got: %q", v)
		}
	})
}

func TestDecoderPreservesDefaultValues(t *testing.T) {
	type nested struct {
		Val string `yaml:"val"`
	}

	type test struct {
		First   string `yaml:"first"`
		Default nested `yaml:"nested"`
	}

	yml := `
first: "Test"
nested:
  # Just some comment here
#  val: "default"
`
	v := test{Default: nested{Val: "default"}}
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		t.Fatal(err)
	}
	if v.Default.Val != "default" {
		t.Fatal("decoder doesn't preserve struct defaults")
	}
}

func TestDecodeError(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "duplicated map key name with anchor-alias",
			source: "&0: *0\n*0:\n*0:",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var v any
			if err := yaml.Unmarshal([]byte(test.source), &v); err == nil {
				t.Fatal("cannot catch decode error")
			}
		})
	}
}

func TestIssue617(t *testing.T) {
	data := `
a: !Not [!Equals [!Ref foo, 'bar']]
`
	var v map[string][]any
	if err := yaml.Unmarshal([]byte(data), &v); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v, map[string][]any{
		"a": {[]any{"foo", "bar"}},
	}) {
		t.Fatalf("found unexpected value: %v", v)
	}
}

type issue337Template struct{}

func (i *issue337Template) UnmarshalYAML(b []byte) error {
	expected := strings.TrimPrefix(`
|
  apiVersion: v1
  kind: ConfigMap
  metadata:
    name: "abc"
    namespace: "abc"
  data:
    foo: FOO
`, "\n")
	if !bytes.Equal(b, []byte(expected)) {
		return fmt.Errorf("expected:\n%s\nbut got:\n%s\n", expected, string(b))
	}
	return nil
}

func TestIssue337(t *testing.T) {
	yml := `
releases:
- name: foo
  chart: ./raw
  values:
  - templates:
    - |
      apiVersion: v1
      kind: ConfigMap
      metadata:
        name: "abc"
        namespace: "abc"
      data:
        foo: FOO
`
	type Value struct {
		Templates []*issue337Template `yaml:"templates"`
	}
	type Release struct {
		Values []*Value `yaml:"values"`
	}
	var v struct {
		Releases []*Release `yaml:"releases"`
	}
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		t.Fatal(err)
	}
}

func TestSetNullValue(t *testing.T) {
	tests := []struct {
		name string
		src  string
	}{
		{
			name: "empty document",
			src:  "",
		},
		{
			name: "null value",
			src:  "null",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Run("set null", func(t *testing.T) {
				var v any
				v = 0x1
				if err := yaml.Unmarshal([]byte(test.src), &v); err != nil {
					t.Fatal(err)
				}
				if v != nil {
					t.Fatal("failed to set nil value")
				}
			})
			t.Run("invalid value", func(t *testing.T) {
				var v *struct{}
				if err := yaml.Unmarshal([]byte(test.src), v); err != nil {
					t.Fatal(err)
				}
				if v != nil {
					t.Fatal("failed to set nil value")
				}
			})
		})
	}
}
