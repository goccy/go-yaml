//go:build !windows

package yaml_test

import (
	"bytes"
	"encoding/json"
	"io"
	"runtime/debug"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/testdata/yaml-test-suite"
)

var failureTestNames = []string{
	"anchors-on-empty-scalars",                         // no json.
	"aliases-in-flow-objects",                          // no json.
	"aliases-in-explicit-block-mapping",                // no json.
	"block-mapping-with-missing-keys",                  // no json.
	"empty-implicit-key-in-single-pair-flow-sequences", // no json.
	"empty-keys-in-block-and-flow-mapping",             // no json.
	"empty-lines-at-end-of-document",                   // no json.
	"flow-mapping-separate-values",                     // no json.
	"flow-sequence-in-flow-mapping",                    // no json.
	"implicit-flow-mapping-key-on-one-line",            // no json.
	"mapping-key-and-flow-sequence-item-anchors",       // no json.
	"nested-implicit-complex-keys",                     // no json.
	"question-mark-edge-cases/00",                      // no json.
	"question-mark-edge-cases/01",                      // no json.
	"single-character-streams/01",                      // no json.
	"single-pair-implicit-entries",                     // no json.
	"spec-example-2-11-mapping-between-sequences",      // no json.
	"spec-example-6-12-separation-spaces",              // no json.
	"spec-example-7-16-flow-mapping-entries",           // no json.
	"spec-example-7-3-completely-empty-flow-nodes",     // no json.
	"spec-example-8-18-implicit-block-mapping-entries", // no json.
	"spec-example-8-19-compact-block-mappings",         // no json.
	"syntax-character-edge-cases/00",                   // no json.
	"tags-on-empty-scalars",                            // no json.
	"various-combinations-of-explicit-block-mappings",  // no json.
	"various-trailing-comments",                        // no json.
	"various-trailing-comments-1-3",                    // no json.
	"zero-indented-sequences-in-explicit-mapping-keys", // no json.

	"colon-at-the-beginning-of-adjacent-flow-scalar",
	"comment-without-whitespace-after-doublequoted-scalar",
	"construct-binary",
	"dash-in-flow-sequence",
	"flow-collections-over-many-lines/01",
	"flow-mapping-colon-on-line-after-key/02",
	"invalid-comment-after-comma",
	"invalid-comment-after-end-of-flow-sequence",
	"invalid-comma-in-tag",
	"plain-dashes-in-flow-sequence",
	"spec-example-9-3-bare-documents",
	"spec-example-9-6-stream",
	"spec-example-9-6-stream-1-3",
	"tabs-in-various-contexts/003",
	"tabs-that-look-like-indentation/04",
	"tag-shorthand-used-in-documents-but-only-defined-in-the-first",
	"trailing-line-of-spaces/01",             // last '\n' character is needed ?
	"wrong-indented-flow-sequence",           // error ?
	"wrong-indented-multiline-quoted-scalar", // error ?
}

var failureTestNameMap map[string]struct{}

func init() {
	failureTestNameMap = make(map[string]struct{})
	for _, name := range failureTestNames {
		failureTestNameMap[name] = struct{}{}
	}
}

func TestYAMLTestSuite(t *testing.T) {
	tests, err := yamltestsuite.TestSuites()
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		if _, exists := failureTestNameMap[test.Name]; exists {
			continue
		}
		t.Run(test.Name, func(t *testing.T) {
			defer func() {
				if e := recover(); e != nil {
					t.Fatalf("panic occurred.\n[input]\n%s\nstack[%s]", string(test.InYAML), debug.Stack())
				}
			}()

			if test.Error {
				var v any
				if err := yaml.Unmarshal(test.InYAML, &v); err == nil {
					t.Fatalf("expected error.\n[input]\n%s\n", string(test.InYAML))
				}
				return
			}

			dec := yaml.NewDecoder(bytes.NewReader(test.InYAML))
			var idx int
			for {
				var v any
				if err := dec.Decode(&v); err != nil {
					if err == io.EOF {
						break
					}
					t.Fatal(err)
				}
				if len(test.InJSON) <= idx {
					t.Fatalf("expected json value is nothing but got %v", v)
				}
				expected, err := json.Marshal(test.InJSON[idx])
				if err != nil {
					t.Fatalf("failed to encode json value: %v", err)
				}
				got, err := json.Marshal(v)
				if err != nil {
					t.Fatalf("failed to encode json value: %v", err)
				}
				if !bytes.Equal(expected, got) {
					t.Fatalf("json mismatch [%s]:\n[expected]\n%s\n[got]\n%s\n", test.Name, string(expected), string(got))
				}
				idx++
			}
		})
	}
	total := len(tests)
	failed := len(failureTestNames)
	passed := total - failed
	passedRate := float32(passed) / float32(total) * 100
	t.Logf("total:[%d] passed:[%d] failure:[%d] passedRate:[%f%%]", total, passed, failed, passedRate)
}
