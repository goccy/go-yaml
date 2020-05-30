package yaml_test

import (
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestPath_Read(t *testing.T) {
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
			path:     (&yaml.PathBuilder{}).Root().Child("store").Child("book").Index(0).Child("author").Build(),
			expected: "john",
		},
		{
			name:     "$.store.book[1].price",
			path:     (&yaml.PathBuilder{}).Root().Child("store").Child("book").Index(1).Child("price").Build(),
			expected: uint64(12),
		},
		{
			name:     "$.store.bicycle.price",
			path:     (&yaml.PathBuilder{}).Root().Child("store").Child("bicycle").Child("price").Build(),
			expected: float64(19.95),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var v interface{}
			if err := test.path.Read(strings.NewReader(yml), &v); err != nil {
				t.Fatalf("%+v", err)
			}
			if test.expected != v {
				t.Fatalf("expected %v(%T). but actual %v(%T)", test.expected, test.expected, v, v)
			}
		})
	}
}
