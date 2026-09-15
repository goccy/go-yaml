package yaml

import (
	"bytes"
	"io"
	"testing"

	"github.com/goccy/go-yaml/parser"
)

func TestParseConsecutiveEmptyDocuments(t *testing.T) {
	f, err := parser.ParseBytes([]byte("---\n---\nk: v\n"), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(f.Docs) != 2 {
		t.Fatalf("docs=%d want 2", len(f.Docs))
	}
	if f.Docs[0].Body != nil {
		t.Fatalf("first doc should be empty, got %v", f.Docs[0].Body)
	}
	if f.Docs[1].Body == nil {
		t.Fatal("second doc should have a body")
	}
}

func TestDecodeEmptyDocumentsBeforeValue(t *testing.T) {
	d := NewDecoder(bytes.NewReader([]byte("---\n---\nk: v\n")))
	var v struct {
		K string `yaml:"k"`
	}
	if err := d.Decode(&v); err != nil {
		t.Fatalf("first decode: %v", err)
	}
	if v.K != "v" {
		t.Fatalf("got %+v want k=v", v)
	}
	if err := d.Decode(&v); err != io.EOF {
		t.Fatalf("second decode: %v want EOF", err)
	}
}
