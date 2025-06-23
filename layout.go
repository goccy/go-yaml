package yaml

import (
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/internal/format"
	"github.com/goccy/go-yaml/parser"
)

// ParseWithComments parses YAML text and keeps comment tokens in the AST.
func ParseWithComments(b []byte) (*ast.File, error) {
	return parser.ParseBytes(b, parser.ParseComments)
}

// FormatFile reconstructs the YAML text from the given AST.
// The order of keys and blank lines are preserved.
func FormatFile(f *ast.File) []byte {
	return []byte(format.FormatFile(f))
}

// MergeYAML merges the YAML source of src into base while preserving layout.
func MergeYAML(base, src []byte) ([]byte, error) {
	baseFile, err := parser.ParseBytes(base, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	srcFile, err := parser.ParseBytes(src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	for i, doc := range srcFile.Docs {
		if i >= len(baseFile.Docs) {
			break
		}
		if err := ast.Merge(baseFile.Docs[i].Body, doc.Body); err != nil {
			return nil, err
		}
	}
	return []byte(format.FormatFile(baseFile)), nil
}
