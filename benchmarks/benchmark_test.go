package benchmarks

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
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
	b.Run("github.com/goccy/go-yaml", func(b *testing.B) {
		var t T
		for i := 0; i < b.N; i++ {
			if err := yaml.Unmarshal([]byte(src), &t); err != nil {
				b.Fatal(err)
			}
		}
	})
}

type Address [20]byte

func (a *Address) UnmarshalText(text []byte) error {
	decoded, err := hex.DecodeString(strings.TrimPrefix(string(text), "0x"))
	if err != nil {
		return err
	}
	copy(a[:], decoded)
	return nil
}

func BenchmarkUnmarshalBigDocument(b *testing.B) {
	type Token struct {
		Address Address `yaml:"address"`
		Name    string  `yaml:"name"`
	}

	type D struct {
		Tokens []Token `yaml:"tokens"`
	}

	expAddr := Address{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef, 0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef, 0x12, 0x34, 0x56, 0x78}

	getDoc := func(n int) string {
		doc := "tokens:\n"
		for i := 0; i < n; i++ {
			doc += "  - address: 0x1234567890abcdef1234567890abcdef12345678\n"
			doc += "    name: token\n"
		}
		return doc
	}

	src := getDoc(2_000)

	b.Run("gopkg.in/yaml.v2", func(b *testing.B) {
		var t D
		for i := 0; i < b.N; i++ {
			if err := goyaml2.Unmarshal([]byte(src), &t); err != nil {
				b.Fatal(err)
			}
		}
		if t.Tokens[0].Address != expAddr {
			b.Fatal("invalid address")
		}
	})

	b.Run("gopkg.in/yaml.v3", func(b *testing.B) {
		var t D
		for i := 0; i < b.N; i++ {
			if err := goyaml3.Unmarshal([]byte(src), &t); err != nil {
				b.Fatal(err)
			}
		}
		if t.Tokens[0].Address != expAddr {
			b.Fatal("invalid address")
		}
	})

	b.Run("github.com/goccy/go-yaml", func(b *testing.B) {
		var t D
		for i := 0; i < b.N; i++ {
			if err := yaml.Unmarshal([]byte(src), &t); err != nil {
				b.Fatal(err)
			}
		}
		if t.Tokens[0].Address != expAddr {
			b.Fatal("invalid address")
		}
	})
}
