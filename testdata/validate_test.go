package yaml_test

import (
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"

	"github.com/goccy/go-yaml"
)

func TestStructValidator(t *testing.T) {
	type Inner struct {
		Required string `validate:"required"`
		Lt10     int    `validate:"lt=10"`
	}

	cases := []struct {
		TestName    string
		YAMLContent string
		ExpectedErr string
		Instance    interface{}
	}{
		{
			TestName: "Test Simple Validation",
			YAMLContent: `---
- name: john
  age: 20
- name: tom
  age: -1
- name: ken
  age: 10`,
			ExpectedErr: `[5:8] Key: 'Age' Error:Field validation for 'Age' failed on the 'gte' tag
   2 | - name: john
   3 |   age: 20
   4 | - name: tom
>  5 |   age: -1
              ^
   6 | - name: ken
   7 |   age: 10`,
			Instance: &[]struct {
				Name string `yaml:"name" validate:"required"`
				Age  int    `yaml:"age" validate:"gte=0,lt=120"`
			}{},
		},
		{
			TestName: "Test Missing Required Field",
			YAMLContent: `---
- name: john
  age: 20
- age: 10`,
			ExpectedErr: `[4:1] Key: 'Name' Error:Field validation for 'Name' failed on the 'required' tag
   1 | ---
   2 | - name: john
   3 |   age: 20
>  4 | - age: 10
       ^
`,
			Instance: &[]struct {
				Name string `yaml:"name" validate:"required"`
				Age  int    `yaml:"age" validate:"gte=0,lt=120"`
			}{},
		},
		{
			TestName: "Test Nested Validation Missing Internal Required",
			YAMLContent: `---
name: john
age: 10
addr:
  number: seven`,
			ExpectedErr: `[4:5] Key: 'State' Error:Field validation for 'State' failed on the 'required' tag
   1 | ---
   2 | name: john
   3 | age: 10
>  4 | addr:
           ^
   5 |   number: seven`,
			Instance: &struct {
				Name string `yaml:"name" validate:"required"`
				Age  int    `yaml:"age" validate:"gte=0,lt=120"`
				Addr struct {
					Number string `yaml:"number" validate:"required"`
					State  string `yaml:"state" validate:"required"`
				} `yaml:"addr"`
			}{},
		},
		{
			TestName: "Test nested Validation with unknown field",
			YAMLContent: `---
name: john
age: 20
addr:
  number: seven
  state: washington
  error: error
`,
			ExpectedErr: `[7:3] unknown field "error"
   4 | addr:
   5 |   number: seven
   6 |   state: washington
>  7 |   error: error
         ^
`,
			Instance: &struct {
				Name string `yaml:"name" validate:"required"`
				Age  int    `yaml:"age" validate:"gte=0,lt=120"`
				Addr *struct {
					Number string `yaml:"number" validate:"required"`
					State  string `yaml:"state" validate:"required"`
				} `yaml:"addr" validate:"required"`
			}{},
		},
		{
			TestName: "Test Validation with wrong field type",
			YAMLContent: `---
name: myDocument
roles:
  name: myRole
  permissions:
    - hello
    - how
    - are
    - you
`,
			ExpectedErr: `[4:7] mapping was used where sequence is expected
   1 | ---
   2 | name: myDocument
   3 | roles:
>  4 |   name: myRole
             ^
   5 |   permissions:
   6 |     - hello
   7 |     - how
   8 |     `,
			Instance: &struct {
				Name  string `yaml:"name"`
				Roles []struct {
					Name        string   `yaml:"name"`
					Permissions []string `yaml:"permissions"`
				} `yaml:"roles"`
			}{},
		},
		{
			TestName: "Test inline validation missing required",
			YAMLContent: `---
name: john
age: 20
`,
			ExpectedErr: `Key: 'Inner.Required' Error:Field validation for 'Required' failed on the 'required' tag`,
			Instance: &struct {
				Name  string `yaml:"name" validate:"required"`
				Age   int    `yaml:"age" validate:"gte=0,lt=120"`
				Inner `yaml:",inline"`
			}{},
		},
		{
			TestName: "Test inline validation field error",
			YAMLContent: `---
name: john
age: 20
required: present
lt10: 20
`,
			ExpectedErr: `[5:7] Key: 'Inner.Lt10' Error:Field validation for 'Lt10' failed on the 'lt' tag
   2 | name: john
   3 | age: 20
   4 | required: present
>  5 | lt10: 20
             ^
`,
			Instance: &struct {
				Name  string `yaml:"name" validate:"required"`
				Age   int    `yaml:"age" validate:"gte=0,lt=120"`
				Inner `yaml:",inline"`
			}{},
		},
		{
			TestName: "Test map",
			YAMLContent: `
map:
  hello: hello
  hello2: hello`,
			ExpectedErr: `[2:4] Key: 'Map' Error:Field validation for 'Map' failed on the 'eq' tag
>  2 | map:
          ^
   3 |   hello: hello
   4 |   hello2: hello`,
			Instance: &struct {
				// Make sure that this map is not set
				Map map[string]string `yaml:"map" validate:"eq=0"`
			}{},
		},
		{
			TestName: "Test sequence",
			YAMLContent: `
map:
  - value
`,
			ExpectedErr: `[2:4] Key: 'Map' Error:Field validation for 'Map' failed on the 'eq' tag
>  2 | map:
          ^
   3 |   - value`,
			Instance: &struct {
				// Make sure that this map is not set
				Map []string `yaml:"map" validate:"eq=0"`
			}{},
		},
	}

	for _, tc := range cases {
		tc := tc // NOTE: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tc.TestName, func(t *testing.T) {
			validate := validator.New()
			dec := yaml.NewDecoder(
				strings.NewReader(tc.YAMLContent),
				yaml.Validator(validate),
				yaml.Strict(),
			)
			err := dec.Decode(tc.Instance)
			switch {
			case tc.ExpectedErr != "" && err == nil:
				t.Fatal("expected error")
			case tc.ExpectedErr == "" && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case tc.ExpectedErr != "" && tc.ExpectedErr != err.Error():
				t.Fatalf("expected `%s` but actual `%s`", tc.ExpectedErr, err.Error())
			}
		})
	}
}
