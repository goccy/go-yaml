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
	"anchors-on-empty-scalars",          // no json.
	"aliases-in-flow-objects",           // no json.
	"aliases-in-explicit-block-mapping", // no json.
	"block-mapping-with-missing-keys",   // no json.
	"colon-at-the-beginning-of-adjacent-flow-scalar",
	"comment-without-whitespace-after-doublequoted-scalar",
	"construct-binary",
	"dash-in-flow-sequence",
	"empty-implicit-key-in-single-pair-flow-sequences", // no json.
	"empty-keys-in-block-and-flow-mapping",             // no json.
	"empty-lines-at-end-of-document",                   // no json.
	"flow-mapping-separate-values",                     // no json.
	"flow-sequence-in-flow-mapping",
	"flow-collections-over-many-lines/01",
	"flow-mapping-colon-on-line-after-key/02",
	"flow-mapping-edge-cases",
	"folded-block-scalar",                   // pass yamlv3.
	"folded-block-scalar-1-3",               // pass yamlv3.
	"implicit-flow-mapping-key-on-one-line", // no json.
	"invalid-comment-after-comma",
	"invalid-comment-after-end-of-flow-sequence",
	"invalid-comma-in-tag",
	"invalid-tag",                                    // pass yamlv3.
	"legal-tab-after-indentation",                    // pass yamlv3.
	"literal-scalars",                                // pass yamlv3.
	"mapping-key-and-flow-sequence-item-anchors",     // no json.
	"multiline-plain-value-with-tabs-on-empty-lines", // pass yamlv3.
	"multiline-scalar-at-top-level",                  // pass yamlv3.
	"multiline-scalar-at-top-level-1-3",              // pass yamlv3.
	"nested-implicit-complex-keys",                   // no json.
	"plain-dashes-in-flow-sequence",
	"question-mark-edge-cases/00",                           // no json.
	"question-mark-edge-cases/01",                           // no json.
	"single-character-streams/01",                           // no json.
	"single-pair-implicit-entries",                          // no json.
	"spec-example-2-11-mapping-between-sequences",           // no json.
	"spec-example-6-12-separation-spaces",                   // no json.
	"spec-example-7-16-flow-mapping-entries",                // no json.
	"spec-example-7-3-completely-empty-flow-nodes",          // no json.
	"spec-example-8-18-implicit-block-mapping-entries",      // no json.
	"spec-example-8-19-compact-block-mappings",              // no json.
	"spec-example-6-6-line-folding",                         // pass yamlv3.
	"spec-example-6-6-line-folding-1-3",                     // pass yamlv3.
	"spec-example-8-10-folded-lines-8-13-final-empty-lines", // pass yamlv3.
	"spec-example-8-2-block-indentation-indicator",
	"spec-example-9-3-bare-documents",
	"spec-example-9-4-explicit-documents",
	"spec-example-9-6-stream",
	"spec-example-9-6-stream-1-3",
	"syntax-character-edge-cases/00", // no json.
	"tab-at-beginning-of-line-followed-by-a-flow-mapping",
	"tab-indented-top-flow",
	"tabs-in-various-contexts/003",
	"tabs-that-look-like-indentation/00",
	"tabs-that-look-like-indentation/01",
	"tabs-that-look-like-indentation/03",
	"tabs-that-look-like-indentation/04",
	"tabs-that-look-like-indentation/05", // pass yamlv3.
	"tabs-that-look-like-indentation/07",
	"tag-shorthand-used-in-documents-but-only-defined-in-the-first",
	"tags-on-empty-scalars",                            // no json.
	"trailing-line-of-spaces/01",                       // last '\n' character is needed ?
	"various-combinations-of-explicit-block-mappings",  // no json.
	"various-trailing-comments",                        // no json.
	"various-trailing-comments-1-3",                    // no json.
	"wrong-indented-flow-sequence",                     // error ?
	"wrong-indented-multiline-quoted-scalar",           // error ?
	"zero-indented-sequences-in-explicit-mapping-keys", // no json.
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
