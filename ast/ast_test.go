package ast

import "testing"

func TestEscapeSingleQuote(t *testing.T) {
	expected := `'Victor''s victory'`
	got := escapeSingleQuote("Victor's victory")
	if got != expected {
		t.Fatalf("expected:%s\ngot:%s", expected, got)
	}
}
