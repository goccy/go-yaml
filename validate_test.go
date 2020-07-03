package yaml_test

import (
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"gopkg.in/go-playground/validator.v9"
)

func TestStructValidator(t *testing.T) {

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
			TestName: "Test Nested Validation Missing Internal Required",
			YAMLContent: `---
name: john
age: 10
addr:
  number: seven`,
			ExpectedErr: "",
			Instance: &struct {
				Name string `yaml:"name" validate:"required"`
				Age  int    `yaml:"age" validate:"gte=0,lt=120"`
				Addr struct {
					Number string `yaml:"number" validate:"required"`
					State  string `yaml:"state" validate:"required"`
				} `yaml:"addr"`
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
