package ast

import (
	"testing"

	"github.com/goccy/go-yaml/token"
)

func TestEscapeSingleQuote(t *testing.T) {
	expected := `'Victor''s victory'`
	got := escapeSingleQuote("Victor's victory")
	if got != expected {
		t.Fatalf("expected:%s\ngot:%s", expected, got)
	}
}

func TestStringNodeStringMultiLineZeroColumn(t *testing.T) {
	node := &StringNode{
		BaseNode: &BaseNode{},
		Token: &token.Token{
			Type:     token.StringType,
			Position: &token.Position{Column: 0, IndentNum: 0},
		},
		Value: "Some content\nSome more content\n",
	}
	if got := node.String(); got == "" {
		t.Fatalf("expected non-empty result, got empty")
	}
}

func TestReadNode(t *testing.T) {
	t.Run("utf-8", func(t *testing.T) {
		value := "éɛทᛞ⠻チ▓🦄"
		node := &StringNode{
			BaseNode: &BaseNode{},
			Token:    &token.Token{},
			Value:    value,
		}
		expectedSize := len(value)
		gotBuffer := make([]byte, expectedSize)
		expectedBuffer := []byte(value)
		gotSize, _ := readNode(gotBuffer, node)
		if gotSize != expectedSize {
			t.Fatalf("expected size:%d\ngot:%d", expectedSize, gotSize)
		}
		if string(gotBuffer) != string(expectedBuffer) {
			t.Fatalf("expected buffer:%s\ngot:%s", expectedBuffer, gotBuffer)
		}
	})
}
