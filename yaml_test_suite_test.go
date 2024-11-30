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
	"aliases-in-implicit-block-mapping",
	"bare-document-after-document-end-marker",
	"block-mapping-with-missing-keys", // no json.
	"block-mapping-with-missing-values",
	"block-mapping-with-multiline-scalars",
	"block-scalar-with-more-spaces-than-first-content-line",
	"block-scalar-with-wrong-indented-line-after-spaces-only",
	"colon-at-the-beginning-of-adjacent-flow-scalar",
	"comment-in-flow-sequence-before-comma",
	"comment-without-whitespace-after-doublequoted-scalar",
	"construct-binary",
	"dash-in-flow-sequence",
	"directive-variants/00",
	"directive-variants/01",
	"double-quoted-scalar-with-escaped-single-quote",
	"duplicate-yaml-directive",
	"escaped-slash-in-double-quotes",
	"explicit-key-and-value-seperated-by-comment", //nolint: misspell
	"extra-words-on-yaml-directive",
	"empty-implicit-key-in-single-pair-flow-sequences", // no json.
	"empty-keys-in-block-and-flow-mapping",             // no json.
	"empty-lines-at-end-of-document",                   // no json.
	"flow-mapping-separate-values",                     // no json.
	"flow-sequence-in-flow-mapping",
	"flow-collections-over-many-lines/01",
	"flow-mapping-colon-on-line-after-key/02",
	"flow-mapping-edge-cases",
	"flow-sequence-with-invalid-comma-at-the-beginning",
	"folded-block-scalar",
	"folded-block-scalar-1-3",
	"implicit-flow-mapping-key-on-one-line", // no json.
	"invalid-comment-after-comma",
	"invalid-comment-after-end-of-flow-sequence",
	"invalid-tag",
	"leading-tabs-in-double-quoted/02",
	"leading-tabs-in-double-quoted/05",
	"legal-tab-after-indentation",
	"literal-block-scalar-with-more-spaces-in-first-line",
	"literal-modifers/00",
	"literal-modifers/01",
	"literal-modifers/02",
	"literal-modifers/03",
	"literal-scalars",
	"mapping-key-and-flow-sequence-item-anchors", // no json.
	"multiline-double-quoted-implicit-keys",
	"multiline-plain-flow-mapping-key",
	"multiline-plain-value-with-tabs-on-empty-lines",
	"multiline-scalar-at-top-level",
	"multiline-scalar-at-top-level-1-3",
	"multiline-single-quoted-implicit-keys",
	"multiline-unidented-double-quoted-block-key",
	"nested-implicit-complex-keys", // no json.
	"node-anchor-not-indented",
	"plain-dashes-in-flow-sequence",
	"plain-url-in-flow-mapping",
	"question-mark-edge-cases/00", // no json.
	"question-mark-edge-cases/01", // no json.
	"scalar-doc-with-in-content/01",
	"scalar-value-with-two-anchors",
	"single-character-streams/01",                      // no json.
	"single-pair-implicit-entries",                     // no json.
	"spec-example-2-11-mapping-between-sequences",      // no json.
	"spec-example-6-12-separation-spaces",              // no json.
	"spec-example-7-16-flow-mapping-entries",           // no json.
	"spec-example-7-3-completely-empty-flow-nodes",     // no json.
	"spec-example-8-18-implicit-block-mapping-entries", // no json.
	"spec-example-8-19-compact-block-mappings",         // no json.
	"spec-example-6-19-secondary-tag-handle",
	"spec-example-6-24-verbatim-tags",
	"spec-example-6-4-line-prefixes",
	"spec-example-6-6-line-folding",
	"spec-example-6-6-line-folding-1-3",
	"spec-example-6-8-flow-folding",
	"spec-example-7-12-plain-lines",
	"spec-example-7-20-single-pair-explicit-entry",
	"spec-example-8-10-folded-lines-8-13-final-empty-lines",
	"spec-example-8-15-block-sequence-entry-types",
	"spec-example-8-17-explicit-block-mapping-entries",
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
	"tabs-that-look-like-indentation/02",
	"tabs-that-look-like-indentation/03",
	"tabs-that-look-like-indentation/04",
	"tabs-that-look-like-indentation/05",
	"tabs-that-look-like-indentation/07",
	"tabs-that-look-like-indentation/08",
	"tag-shorthand-used-in-documents-but-only-defined-in-the-first",
	"tags-for-block-objects",
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
