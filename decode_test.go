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
			"v: false\n",
			map[string]bool{"v": false},
		},
		{
			"v: 10\n",
			map[string]int{"v": 10},
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
			map[string]interface{}{},
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

func TestDecodeByAnchorOfOtherFile(t *testing.T) {
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

