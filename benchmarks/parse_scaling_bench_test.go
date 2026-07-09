package benchmarks

import (
	"fmt"
	"strings"
	"testing"

	"github.com/goccy/go-yaml/parser"
)

// BenchmarkParseScaling guards against the O(n^2) parse-time regression in
// parser.parseMap (a block mapping with many sibling keys used to be parsed by
// recursing once per sibling and re-copying the accumulated values at every
// level, which is quadratic in the number of entries).
//
// Each sub-benchmark parses a flat mapping of n sibling keys:
//
//	key0: value0
//	key1: value1
//	...
//
// Because SetBytes is reported, the throughput (MB/s) is the signal: it stays
// roughly flat across sizes when parsing is linear (O(n)) and falls off as the
// document grows when parsing is quadratic (O(n^2)).
func BenchmarkParseScaling(b *testing.B) {
	for _, n := range []int{10_000, 50_000, 100_000, 200_000} {
		src := generateFlatMap(n)
		b.Run(fmt.Sprintf("nodes=%d", n), func(b *testing.B) {
			b.SetBytes(int64(len(src)))
			b.ReportAllocs()
			for b.Loop() {
				if _, err := parser.ParseBytes(src, 0); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// generateFlatMap builds a block mapping of n sibling keys, all at the same
// column — the pathological input for the parser's block-mapping loop.
func generateFlatMap(n int) []byte {
	var b strings.Builder
	b.Grow(n * 20)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "key%d: value%d\n", i, i)
	}
	return []byte(b.String())
}
