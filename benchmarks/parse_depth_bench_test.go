package benchmarks

import (
	"strings"
	"testing"

	"github.com/goccy/go-yaml/parser"
)

// BenchmarkParseDeeplyNested measures parsing a single deeply nested flow
// sequence. Every nesting level records its full path, so an unbounded parser
// scales quadratically with the depth; the parser now caps the depth so such
// input is rejected quickly.
func BenchmarkParseDeeplyNested(b *testing.B) {
	depth := 50000
	src := []byte(strings.Repeat("[", depth) + strings.Repeat("]", depth))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.ParseBytes(src, 0)
	}
}
