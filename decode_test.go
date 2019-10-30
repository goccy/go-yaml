package yaml_test

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

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

		// Some cross type conversions
		{
			"v: 42",
			map[string]uint{"v": 42},
		}, {
			"v: -42",
			map[string]uint{},
		}, {
			"v: 4294967296",
			map[string]uint64{"v": 4294967296},
		}, {
			"v: -4294967296",
			map[string]uint64{},
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
		{
			"v: -1",
			map[string]uint{},
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
		{
			"v: -1",
			map[string]uint64{},
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

		// Overflow cases.
		{
			"v: 4294967297",
			map[string]int32{},
		}, {
			"v: 128",
			map[string]int8{},
		},

		// Quoted values.
		{
			"'1': '\"2\"'",
			map[interface{}]interface{}{"1": "\"2\""},
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
			"v: .inf\n",
			map[string]interface{}{"v": math.Inf(0)},
		},
		{
			"v: -.inf\n",
			map[string]interface{}{"v": math.Inf(-1)},
		},
		{
			"v: .nan\n",
			map[string]interface{}{"v": math.NaN()},
		},
		{
			"v: null\n",
			map[string]interface{}{"v": nil},
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
	}
	for _, test := range tests {
		buf := bytes.NewBufferString(test.source)
		dec := yaml.NewDecoder(buf)
		typ := reflect.ValueOf(test.value).Type()
		value := reflect.New(typ)
		if err := dec.Decode(value.Interface()); err != nil {
			t.Fatalf("%s: %+v", test.source, err)
		}
		actual := fmt.Sprintf("%v", value.Elem().Interface())
		expect := fmt.Sprintf("%v", test.value)
		if actual != expect {
			t.Fatal("failed to test. ", actual, expect)
		}
	}
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
