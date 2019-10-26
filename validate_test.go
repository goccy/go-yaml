package yaml_test

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"gopkg.in/go-playground/validator.v9"
)

type Person struct {
	Name string `validate:"required"`
	Age  int    `validate:"gte=0,lt=120"`
}

func ExampleStructValidator() {
	yml := `---
- name: john
  age: 20
- name: tom
  age: -1
- name: ken
  age: 10
`
	validate := validator.New()
	dec := yaml.NewDecoder(
		strings.NewReader(yml),
		yaml.Validator(validate),
	)
	var v []*Person
	err := dec.Decode(&v)
	if err == nil {
		panic(err)
	}
	fmt.Printf("%v", err)
	// OUTPUT:
	// [5:8] Key: 'Person.Age' Error:Field validation for 'Age' failed on the 'gte' tag
	//        1 | ---
	//        2 | - name: john
	//        3 |   age: 20
	//        4 | - name: tom
	//     >  5 |   age: -1
	//                  ^
	//        6 | - name: ken
	//        7 |   age: 10
}
