package yaml_test

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
)

func builder() *yaml.PathBuilder { return &yaml.PathBuilder{} }

func TestPathBuilder(t *testing.T) {
	tests := []struct {
		expected string
		path     *yaml.Path
	}{
		{
			expected: `$.a.b[0]`,
			path:     builder().Root().Child("a").Child("b").Index(0).Build(),
		},
		{
			expected: `$.'a.b'.'c*d'`,
			path:     builder().Root().Child("a.b").Child("c*d").Build(),
		},
		{
			expected: `$.'a.b-*'.c`,
			path:     builder().Root().Child("a.b-*").Child("c").Build(),
		},
		{
			expected: `$.'a'.b`,
			path:     builder().Root().Child("'a'").Child("b").Build(),
		},
		{
			expected: `$.'a.b'.c`,
			path:     builder().Root().Child("'a.b'").Child("c").Build(),
		},
	}
	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			expected := test.expected
			got := test.path.String()
			if expected != got {
				t.Fatalf("failed to build path. expected:[%q] but got:[%q]", expected, got)
			}
		})
	}
}

func TestPath(t *testing.T) {
	yml := `
store:
  book:
    - author: john
      price: 10
    - author: ken
      price: 12
  bicycle:
    color: red
    price: 19.95
  bicycle*unicycle:
    price: 20.25
`
	tests := []struct {
		name     string
		path     *yaml.Path
		expected interface{}
	}{
		{
			name:     "$.store.book[0].author",
			path:     builder().Root().Child("store").Child("book").Index(0).Child("author").Build(),
			expected: "john",
		},
		{
			name:     "$.store.book[1].price",
			path:     builder().Root().Child("store").Child("book").Index(1).Child("price").Build(),
			expected: uint64(12),
		},
		{
			name:     "$.store.book[*].author",
			path:     builder().Root().Child("store").Child("book").IndexAll().Child("author").Build(),
			expected: []interface{}{"john", "ken"},
		},
		{
			name: "$.store.book[*]",
			path: builder().Root().Child("store").Child("book").IndexAll().Build(),
			expected: []interface{}{
				map[string]interface{}{
					"author": "john",
					"price":  uint64(10),
				},
				map[string]interface{}{
					"author": "ken",
					"price":  uint64(12),
				},
			},
		},
		{
			name: "$..book[*]",
			path: builder().Root().Recursive("book").IndexAll().Build(),
			expected: []interface{}{
				[]interface{}{
					map[string]interface{}{
						"author": "john",
						"price":  uint64(10),
					},
					map[string]interface{}{
						"author": "ken",
						"price":  uint64(12),
					},
				},
			},
		},
		{
			name:     "$.store.book[0]",
			path:     builder().Root().Child("store").Child("book").Index(0).Build(),
			expected: map[string]interface{}{"author": "john", "price": uint64(10)},
		},
		{
			name:     "$..author",
			path:     builder().Root().Recursive("author").Build(),
			expected: []interface{}{"john", "ken"},
		},
		{
			name:     "$.store.bicycle.price",
			path:     builder().Root().Child("store").Child("bicycle").Child("price").Build(),
			expected: float64(19.95),
		},
		{
			name:     `$.store.'bicycle*unicycle'.price`,
			path:     builder().Root().Child("store").Child(`bicycle*unicycle`).Child("price").Build(),
			expected: float64(20.25),
		},
		{
			name: "$",
			path: builder().Root().Build(),
			expected: map[string]interface{}{
				"store": map[string]interface{}{
					"book": []interface{}{
						map[string]interface{}{
							"author": "john",
							"price":  uint64(10),
						},
						map[string]interface{}{
							"author": "ken",
							"price":  uint64(12),
						},
					},
					"bicycle": map[string]interface{}{
						"color": "red",
						"price": 19.95,
					},
					"bicycle*unicycle": map[string]interface{}{
						"price": 20.25,
					},
				},
			},
		},
	}
	t.Run("PathString", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				path, err := yaml.PathString(test.name)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				if test.name != path.String() {
					t.Fatalf("expected %s but actual %s", test.name, path.String())
				}
			})
		}
	})
	t.Run("string", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				if test.name != test.path.String() {
					t.Fatalf("expected %s but actual %s", test.name, test.path.String())
				}
			})
		}
	})
	t.Run("read", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				var v interface{}
				if err := test.path.Read(strings.NewReader(yml), &v); err != nil {
					t.Fatalf("%+v", err)
				}
				if !reflect.DeepEqual(test.expected, v) {
					t.Fatalf("expected %v(%T). but actual %v(%T)", test.expected, test.expected, v, v)
				}
			})
		}
	})
	t.Run("filter", func(t *testing.T) {
		var target interface{}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatalf("failed to unmarshal: %+v", err)
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				var v interface{}
				if err := test.path.Filter(target, &v); err != nil {
					t.Fatalf("%+v", err)
				}
				if !reflect.DeepEqual(test.expected, v) {
					t.Fatalf("expected %v(%T). but actual %v(%T)", test.expected, test.expected, v, v)
				}
			})
		}
	})
}

func TestPath_ReservedKeyword(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		src      string
		expected interface{}
		failure  bool
	}{
		{
			name: "quoted path",
			path: `$.'a.b.c'.foo`,
			src: `
a.b.c:
  foo: bar
`,
			expected: "bar",
		},
		{
			name:     "contains quote key",
			path:     `$.a'b`,
			src:      `a'b: 10`,
			expected: uint64(10),
		},
		{
			name:     "escaped quote",
			path:     `$.'alice\'s age'`,
			src:      `alice's age: 10`,
			expected: uint64(10),
		},
		{
			name:     "directly use white space",
			path:     `$.a  b`,
			src:      `a  b: 10`,
			expected: uint64(10),
		},
		{
			name:    "empty quoted key",
			path:    `$.''`,
			src:     `a: 10`,
			failure: true,
		},
		{
			name:    "unterminated quote",
			path:    `$.'abcd`,
			src:     `abcd: 10`,
			failure: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path, err := yaml.PathString(test.path)
			if test.failure {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			} else {
				if err != nil {
					t.Fatalf("%+v", err)
				}
			}
			file, err := parser.ParseBytes([]byte(test.src), 0)
			if err != nil {
				t.Fatal(err)
			}
			var v interface{}
			if err := path.Read(file, &v); err != nil {
				t.Fatalf("%+v", err)
			}
			if v != test.expected {
				t.Fatalf("failed to get value. expected:[%v] but got:[%v]", test.expected, v)
			}
		})
	}
}

func TestPath_Invalid(t *testing.T) {
	tests := []struct {
		path string
		src  string
	}{
		{
			path: "$.wrong",
			src:  "foo: bar",
		},
	}
	for _, test := range tests {
		path, err := yaml.PathString(test.path)
		if err != nil {
			t.Fatal(err)
		}
		t.Run("path.Read", func(t *testing.T) {
			file, err := parser.ParseBytes([]byte(test.src), 0)
			if err != nil {
				t.Fatal(err)
			}
			var v interface{}
			err = path.Read(file, &v)
			if err == nil {
				t.Fatal("expected error")
			}
			if !yaml.IsNotFoundNodeError(err) {
				t.Fatalf("unexpected error %s", err)
			}
		})
		t.Run("path.ReadNode", func(t *testing.T) {
			file, err := parser.ParseBytes([]byte(test.src), 0)
			if err != nil {
				t.Fatal(err)
			}
			_, err = path.ReadNode(file)
			if err == nil {
				t.Fatal("expected error")
			}
			if !yaml.IsNotFoundNodeError(err) {
				t.Fatalf("unexpected error %s", err)
			}
		})
	}
}

func TestPath_ReadNode(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		src      string
		expected interface{}
	}{
		{
			name: "nested array sequence",
			path: `$.a.b[0].c`,
			src: `
a:
  b:
   - c: 123
  e: |
   Line1
   Line2
`,
			expected: uint64(123),
		},
		{
			name: "nested array sequence issue#281",
			path: `$..a.c`,
			src: `
s:
  - a:
      b: u1
      c: get1
      d: i1
  - w:
      c: bad
      e:
        - a:
           b: u2
           c: get2
           d: i2
`,
			// The expected values are
			// - get1
			// - get2
			expected: []interface{}{
				map[string]interface{}{
					"b": "u1",
					"c": "get1",
					"d": "i1",
				},
				map[string]interface{}{
					"b": "u2",
					"c": "get2",
					"d": "i2",
				},
			},
		},
		{
			name: "nested array sequence issue#281",
			path: `$..c`,
			src: `
s:
  - a:
      b: u1
      c: get1
      d: i1
  - w:
      c: bad
      e:
        - a:
            b: u2
            c: get2
            d: i2
`,
			expected: []interface{}{"get1", "bad", "get2"},
		},
		{
			name: "nested array sequence issue#281",
			path: `$.s[0].a.c`,
			src: `
s:
  - a:
      b: u1
      c: get1
      d: i1
  - w:
      c: bad
      e:
        - a:
            b: u2
            c: get2
            d: i2
`,
			expected: "get1",
		},
		{
			name: "nested array sequence issue#281",
			path: "$.s[*].a.c",
			src: `
s:
  - a:
      b: u1
      c: get1
      d: i1
  - w:
      c: bad
      e:
        - a:
            b: u2
            c: get2
            d: i2
`,
			expected: []interface{}{"get1"},
		},
	}
	for _, test := range tests {
		path, err := yaml.PathString(test.path)
		if err != nil {
			t.Fatal(err)
		}
		t.Run(fmt.Sprintf("path.ReadNode %s path %s", test.name, test.path), func(t *testing.T) {
			file, err := parser.ParseBytes([]byte(test.src), 0)
			if err != nil {
				t.Fatal(err)
			}
			n, err := path.ReadNode(file)
			if err != nil {
				t.Fatal("expected error", err)
			}
			var v interface{}
			err = yaml.Unmarshal([]byte(n.String()), &v)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(test.expected, v) {
				t.Fatalf("expected %v(%T) but got %v(%T)", test.expected, test.expected, v, v)
			}
		})
	}
}

func TestPath_Merge(t *testing.T) {
	tests := []struct {
		path     string
		dst      string
		src      string
		expected string
	}{
		{
			path: "$.c",
			dst: `
a: 1
b: 2
c:
  d: 3
  e: 4
`,
			src: `
f: 5
g: 6
`,
			expected: `
a: 1
b: 2
c:
  d: 3
  e: 4
  f: 5
  g: 6
`,
		},
		{
			path: "$.a.b",
			dst: `
a:
  b:
   - 1
   - 2
`,
			src: `
- 3
- map:
   - 4
   - 5
`,
			expected: `
a:
  b:
   - 1
   - 2
   - 3
   - map:
      - 4
      - 5
`,
		},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			path, err := yaml.PathString(test.path)
			if err != nil {
				t.Fatalf("%+v", err)
			}
			t.Run("FromReader", func(t *testing.T) {
				file, err := parser.ParseBytes([]byte(test.dst), 0)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				if err := path.MergeFromReader(file, strings.NewReader(test.src)); err != nil {
					t.Fatalf("%+v", err)
				}
				actual := "\n" + file.String()
				if test.expected != actual {
					t.Fatalf("expected: %q. but got %q", test.expected, actual)
				}
			})
			t.Run("FromFile", func(t *testing.T) {
				file, err := parser.ParseBytes([]byte(test.dst), 0)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				src, err := parser.ParseBytes([]byte(test.src), 0)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				if err := path.MergeFromFile(file, src); err != nil {
					t.Fatalf("%+v", err)
				}
				actual := "\n" + file.String()
				if test.expected != actual {
					t.Fatalf("expected: %q. but got %q", test.expected, actual)
				}
			})
			t.Run("FromNode", func(t *testing.T) {
				file, err := parser.ParseBytes([]byte(test.dst), 0)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				src, err := parser.ParseBytes([]byte(test.src), 0)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				if len(src.Docs) == 0 {
					t.Fatalf("failed to parse")
				}
				if err := path.MergeFromNode(file, src.Docs[0]); err != nil {
					t.Fatalf("%+v", err)
				}
				actual := "\n" + file.String()
				if test.expected != actual {
					t.Fatalf("expected: %q. but got %q", test.expected, actual)
				}
			})
		})
	}
}

func TestPath_Replace(t *testing.T) {
	tests := []struct {
		path     string
		dst      string
		src      string
		expected string
	}{
		{
			path: "$.a",
			dst: `
a: 1
b: 2
`,
			src: `3`,
			expected: `
a: 3
b: 2
`,
		},
		{
			path: "$.b",
			dst: `
b: 1
c: 2
`,
			src: `
d: e
f:
  g: h
  i: j
`,
			expected: `
b:
  d: e
  f:
    g: h
    i: j
c: 2
`,
		},
		{
			path: "$.a.b[0]",
			dst: `
a:
  b:
  - hello
c: 2
`,
			src: `world`,
			expected: `
a:
  b:
  - world
c: 2
`,
		},

		{
			path: "$.books[*].author",
			dst: `
books:
  - name: book_a
    author: none
  - name: book_b
    author: none
pictures:
  - name: picture_a
    author: none
  - name: picture_b
    author: none
building:
  author: none
`,
			src: `ken`,
			expected: `
books:
  - name: book_a
    author: ken
  - name: book_b
    author: ken
pictures:
  - name: picture_a
    author: none
  - name: picture_b
    author: none
building:
  author: none
`,
		},
		{
			path: "$..author",
			dst: `
books:
  - name: book_a
    author: none
  - name: book_b
    author: none
pictures:
  - name: picture_a
    author: none
  - name: picture_b
    author: none
building:
  author: none
`,
			src: `ken`,
			expected: `
books:
  - name: book_a
    author: ken
  - name: book_b
    author: ken
pictures:
  - name: picture_a
    author: ken
  - name: picture_b
    author: ken
building:
  author: ken
`,
		},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			path, err := yaml.PathString(test.path)
			if err != nil {
				t.Fatalf("%+v", err)
			}
			t.Run("WithReader", func(t *testing.T) {
				file, err := parser.ParseBytes([]byte(test.dst), 0)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				if err := path.ReplaceWithReader(file, strings.NewReader(test.src)); err != nil {
					t.Fatalf("%+v", err)
				}
				actual := "\n" + file.String()
				if test.expected != actual {
					t.Fatalf("expected: %q. but got %q", test.expected, actual)
				}
			})
			t.Run("WithFile", func(t *testing.T) {
				file, err := parser.ParseBytes([]byte(test.dst), 0)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				src, err := parser.ParseBytes([]byte(test.src), 0)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				if err := path.ReplaceWithFile(file, src); err != nil {
					t.Fatalf("%+v", err)
				}
				actual := "\n" + file.String()
				if test.expected != actual {
					t.Fatalf("expected: %q. but got %q", test.expected, actual)
				}
			})
			t.Run("WithNode", func(t *testing.T) {
				file, err := parser.ParseBytes([]byte(test.dst), 0)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				src, err := parser.ParseBytes([]byte(test.src), 0)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				if len(src.Docs) == 0 {
					t.Fatalf("failed to parse")
				}
				if err := path.ReplaceWithNode(file, src.Docs[0]); err != nil {
					t.Fatalf("%+v", err)
				}
				actual := "\n" + file.String()
				if test.expected != actual {
					t.Fatalf("expected: %q. but got %q", test.expected, actual)
				}
			})
		})
	}
}

func TestInvalidPath(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "missing root with dot",
			path: ".foo",
		},
		{
			name: "missing root with index",
			path: "foo[0]",
		},
		{
			name: "missing root with recursive",
			path: "..foo",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := yaml.PathString(test.path); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func ExamplePath_AnnotateSource() {
	yml := `
a: 1
b: "hello"
`
	var v struct {
		A int
		B string
	}
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		panic(err)
	}
	if v.A != 2 {
		// output error with YAML source
		path, err := yaml.PathString("$.a")
		if err != nil {
			log.Fatal(err)
		}
		source, err := path.AnnotateSource([]byte(yml), false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("a value expected 2 but actual %d:\n%s\n", v.A, string(source))
	}
	// OUTPUT:
	// a value expected 2 but actual 1:
	// >  2 | a: 1
	//           ^
	//    3 | b: "hello"
}

func ExamplePath_AnnotateSource_withComment() {
	yml := `
# This is my document
doc:
  # This comment should be line 3
  map:
    # And below should be line 5
    - value1
    - value2
  other: value3
`
	path, err := yaml.PathString("$.doc.map[0]")
	if err != nil {
		log.Fatal(err)
	}
	msg, err := path.AnnotateSource([]byte(yml), false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(msg))
	// OUTPUT:
	//    4 |   # This comment should be line 3
	//    5 |   map:
	//    6 |     # And below should be line 5
	// >  7 |     - value1
	//              ^
	//    8 |     - value2
	//    9 |   other: value3
}

func ExamplePathString() {
	yml := `
store:
  book:
    - author: john
      price: 10
    - author: ken
      price: 12
  bicycle:
    color: red
    price: 19.95
`
	path, err := yaml.PathString("$.store.book[*].author")
	if err != nil {
		log.Fatal(err)
	}
	var authors []string
	if err := path.Read(strings.NewReader(yml), &authors); err != nil {
		log.Fatal(err)
	}
	fmt.Println(authors)
	// OUTPUT:
	// [john ken]
}
