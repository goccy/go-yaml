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
	"anchors-on-empty-scalars",
	"aliases-in-flow-objects",
	"aliases-in-explicit-block-mapping",
	"aliases-in-implicit-block-mapping",
	"allowed-characters-in-alias",
	"anchor-before-sequence-entry-on-same-line",
	"anchor-for-empty-node",
	"anchor-plus-alias",
	"anchors-in-mapping",
	"anchors-with-colon-in-name",
	"bare-document-after-document-end-marker",
	"block-mapping-with-missing-keys",
	"block-mapping-with-missing-values",
	"block-mapping-with-multiline-scalars",
	"block-scalar-with-more-spaces-than-first-content-line",
	"block-scalar-with-wrong-indented-line-after-spaces-only",
	"colon-at-the-beginning-of-adjacent-flow-scalar",
	"comment-between-plain-scalar-lines",
	"comment-in-flow-sequence-before-comma",
	"comment-without-whitespace-after-doublequoted-scalar",
	"construct-binary",
	"dash-in-flow-sequence",
	"directive-variants/00",
	"directive-variants/01",
	"double-quoted-scalar-with-escaped-single-quote",
	"duplicate-yaml-directive",
	"escaped-slash-in-double-quotes",
	"explicit-key-and-value-seperated-by-comment",
	"extra-words-on-yaml-directive",
	"empty-implicit-key-in-single-pair-flow-sequences",
	"empty-keys-in-block-and-flow-mapping",
	"empty-lines-at-end-of-document",
	"flow-mapping-separate-values",
	"flow-sequence-in-flow-mapping",
	"flow-collections-over-many-lines/01",
	"flow-mapping-colon-on-line-after-key/02",
	"flow-mapping-edge-cases",
	"flow-sequence-with-invalid-comma-at-the-beginning",
	"flow-sequence-with-invalid-extra-comma",
	"folded-block-scalar",
	"folded-block-scalar-1-3",
	"implicit-flow-mapping-key-on-one-line",
	"invalid-anchor-in-zero-indented-sequence",
	"invalid-comment-after-comma",
	"invalid-comment-after-end-of-flow-sequence",
	"invalid-document-end-marker-in-single-quoted-string",
	"invalid-document-markers-in-flow-style",
	"invalid-document-start-marker-in-doublequoted-tring",
	"invalid-escape-in-double-quoted-string",
	"invalid-item-after-end-of-flow-sequence",
	"invalid-mapping-after-sequence",
	"invalid-mapping-in-plain-single-line-value",
	"invalid-nested-mapping",
	"invalid-scalar-after-sequence",
	"invalid-tag",
	"key-with-anchor-after-missing-explicit-mapping-value",
	"leading-tab-content-in-literals/00",
	"leading-tab-content-in-literals/01",
	"leading-tabs-in-double-quoted/02",
	"leading-tabs-in-double-quoted/05",
	"legal-tab-after-indentation",
	"literal-block-scalar-with-more-spaces-in-first-line",
	"literal-modifers/00",
	"literal-modifers/01",
	"literal-modifers/02",
	"literal-modifers/03",
	"literal-scalars",
	"mapping-key-and-flow-sequence-item-anchors",
	"mapping-starting-at-line",
	"mapping-with-anchor-on-document-start-line",
	"missing-document-end-marker-before-directive",
	"mixed-block-mapping-explicit-to-implicit",
	"multiline-double-quoted-implicit-keys",
	"multiline-plain-flow-mapping-key",
	"multiline-plain-flow-mapping-key-without-value",
	"multiline-plain-value-with-tabs-on-empty-lines",
	"multiline-scalar-at-top-level",
	"multiline-scalar-at-top-level-1-3",
	"multiline-single-quoted-implicit-keys",
	"multiline-unidented-double-quoted-block-key",
	"nested-implicit-complex-keys",
	"need-document-footer-before-directives",
	"node-anchor-in-sequence",
	"node-anchor-not-indented",
	"plain-dashes-in-flow-sequence",
	"plain-url-in-flow-mapping",
	"question-mark-at-start-of-flow-key",
	"question-mark-edge-cases/00",
	"question-mark-edge-cases/01",
	"scalar-doc-with-in-content/01",
	"scalar-value-with-two-anchors",
	"single-character-streams/01",
	"single-pair-implicit-entries",
	"spec-example-2-11-mapping-between-sequences",
	"spec-example-6-12-separation-spaces",
	"spec-example-7-16-flow-mapping-entries",
	"spec-example-7-3-completely-empty-flow-nodes",
	"spec-example-8-18-implicit-block-mapping-entries",
	"spec-example-8-19-compact-block-mappings",
	"spec-example-2-24-global-tags",
	"spec-example-2-25-unordered-sets",
	"spec-example-2-26-ordered-mappings",
	"spec-example-5-12-tabs-and-spaces",
	"spec-example-5-3-block-structure-indicators",
	"spec-example-5-9-directive-indicator",
	"spec-example-6-1-indentation-spaces",
	"spec-example-6-13-reserved-directives",
	"spec-example-6-19-secondary-tag-handle",
	"spec-example-6-2-indentation-indicators",
	"spec-example-6-21-local-tag-prefix",
	"spec-example-6-23-node-properties",
	"spec-example-6-24-verbatim-tags",
	"spec-example-6-28-non-specific-tags",
	"spec-example-6-3-separation-spaces",
	"spec-example-6-4-line-prefixes",
	"spec-example-6-6-line-folding",
	"spec-example-6-6-line-folding-1-3",
	"spec-example-6-7-block-folding",
	"spec-example-6-8-flow-folding",
	"spec-example-7-12-plain-lines",
	"spec-example-7-19-single-pair-flow-mappings",
	"spec-example-7-2-empty-content",
	"spec-example-7-20-single-pair-explicit-entry",
	"spec-example-7-24-flow-nodes",
	"spec-example-7-6-double-quoted-lines",
	"spec-example-7-9-single-quoted-lines",
	"spec-example-8-10-folded-lines-8-13-final-empty-lines",
	"spec-example-8-15-block-sequence-entry-types",
	"spec-example-8-17-explicit-block-mapping-entries",
	"spec-example-8-2-block-indentation-indicator",
	"spec-example-8-22-block-collection-nodes",
	"spec-example-8-7-literal-scalar",
	"spec-example-8-7-literal-scalar-1-3",
	"spec-example-8-8-literal-content",
	"spec-example-9-3-bare-documents",
	"spec-example-9-4-explicit-documents",
	"spec-example-9-5-directives-documents",
	"spec-example-9-6-stream",
	"spec-example-9-6-stream-1-3",
	"syntax-character-edge-cases/00",
	"tab-at-beginning-of-line-followed-by-a-flow-mapping",
	"tab-indented-top-flow",
	"tabs-in-various-contexts/001",
	"tabs-in-various-contexts/002",
	"tabs-in-various-contexts/004",
	"tabs-in-various-contexts/005",
	"tabs-in-various-contexts/006",
	"tabs-in-various-contexts/008",
	"tabs-in-various-contexts/010",
	"tabs-that-look-like-indentation/00",
	"tabs-that-look-like-indentation/01",
	"tabs-that-look-like-indentation/02",
	"tabs-that-look-like-indentation/03",
	"tabs-that-look-like-indentation/04",
	"tabs-that-look-like-indentation/05",
	"tabs-that-look-like-indentation/07",
	"tabs-that-look-like-indentation/08",
	"tags-for-block-objects",
	"tags-for-flow-objects",
	"tags-for-root-objects",
	"tags-in-explicit-mapping",
	"tags-in-implicit-mapping",
	"tags-on-empty-scalars",
	"three-dashes-and-content-without-space",
	"trailing-line-of-spaces/01",                       // last '\n' character is needed ?
	"various-combinations-of-explicit-block-mappings",  // no json
	"various-trailing-comments",                        // no json
	"various-trailing-comments-1-3",                    // no json
	"wrong-indented-flow-sequence",                     // error ?
	"wrong-indented-multiline-quoted-scalar",           // error ?
	"zero-indented-sequences-in-explicit-mapping-keys", // no json
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
