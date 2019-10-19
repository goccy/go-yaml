package yaml_test

import (
	"bytes"
	"math"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestEncoder(t *testing.T) {
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
			"v: 1.0\n",
			map[string]float64{"v": 1.0},
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
			map[string][]string{"v": []string{"A", "B"}},
		},
		{
			"a: -\n",
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
						"B": []int{2, 3},
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
			"t2: 2018-01-09T10:40:47Z\nt4: 2098-01-09T10:40:47Z\n",
			map[string]string{
				"t2": "2018-01-09T10:40:47Z",
				"t4": "2098-01-09T10:40:47Z",
			},
		},
		{
			"a:\n  b: c\n  d: e\n",
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
			"a: \"b: c\"\n",
			map[string]string{"a": "b: c"},
		},
		{
			"a: \"Hello #comment\"\n",
			map[string]string{"a": "Hello #comment"},
		},
		{
			"a: 100.5\n",
			map[string]interface{}{
				"a": 100.5,
			},
		},
		{
			"a: \"\\\\0\"\n",
			map[string]string{"a": "\\0"},
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
			"a: 1\n",
			struct {
				A int
				B int `yaml:"-"`
			}{
				1, 0,
			},
		},
	}
	for _, test := range tests {
		var buf bytes.Buffer
		enc := yaml.NewEncoder(&buf)
		if err := enc.Encode(test.value); err != nil {
			t.Fatalf("%+v", err)
		}
		if test.source != buf.String() {
			t.Fatalf("expect = [%s], actual = [%s]", test.source, buf.String())
		}
	}
}
