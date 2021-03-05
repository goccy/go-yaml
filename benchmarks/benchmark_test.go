package benchmarks

import (
	"testing"

	"github.com/pulumi/go-yaml"
	goyaml2 "gopkg.in/yaml.v2"
	goyaml3 "gopkg.in/yaml.v3"
)

func Benchmark(b *testing.B) {
	const src = `---
id: 1
message: Hello, World
verified: true
elements:
  - one
  - 0.02
  - null
  - -inf
`
	type T struct {
		ID       int    `yaml:"id"`
		Message  string `yaml:"message"`
		Verified bool   `yaml:"verified,omitempty"`
	}

	b.Run("gopkg.in/yaml.v2", func(b *testing.B) {
		var t T
		for i := 0; i < b.N; i++ {
			if err := goyaml2.Unmarshal([]byte(src), &t); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("gopkg.in/yaml.v3", func(b *testing.B) {
		var t T
		for i := 0; i < b.N; i++ {
			if err := goyaml3.Unmarshal([]byte(src), &t); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("github.com/pulumi/go-yaml", func(b *testing.B) {
		var t T
		for i := 0; i < b.N; i++ {
			if err := yaml.Unmarshal([]byte(src), &t); err != nil {
				b.Fatal(err)
			}
		}
	})
}
