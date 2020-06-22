package yaml_test

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

func builder() *yaml.PathBuilder { return &yaml.PathBuilder{} }

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

func ExamplePath_PathString() {
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
