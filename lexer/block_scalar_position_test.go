package lexer_test

import (
	"strings"
	"testing"

	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
)

func TestTokenizeBlockScalarStartLine(t *testing.T) {
	styles := []struct {
		header string
		value  string
		chomp  string
	}{
		{"|", "first\nsecond", "clip"},
		{"|-", "first\nsecond", "strip"},
		{"|+", "first\nsecond", "keep"},
		{">", "first second", "clip"},
		{">-", "first second", "strip"},
		{">+", "first second", "keep"},
	}
	endings := []struct {
		name     string
		trailing string
		next     string
	}{
		{"eof_without_newline", "", ""},
		{"eof_after_newline", "\n", ""},
		{"eof_after_blank_lines", "\n\n\n", ""},
		{"comment", "\n\n", "# comment\n"},
		{"mapping", "\n", "next: value\n"},
		{"document_start", "\n\n", "---\nnext: value\n"},
		{"document_end", "\n\n", "...\n"},
	}
	for _, newline := range []struct{ name, value string }{{"LF", "\n"}, {"CRLF", "\r\n"}, {"CR", "\r"}} {
		for _, style := range styles {
			for _, ending := range endings {
				t.Run(newline.name+"/"+style.header+"/"+ending.name, func(t *testing.T) {
					body := "  first\n  second" + ending.trailing
					src := strings.ReplaceAll("\nkey: "+style.header+"\n"+body+ending.next, "\n", newline.value)
					tokens := lexer.Tokenize(src)
					var scalar *token.Token
					for i, tk := range tokens {
						if tk.Type == token.LiteralType || tk.Type == token.FoldedType {
							if i+1 < len(tokens) {
								scalar = tokens[i+1]
							}
							break
						}
					}
					if scalar == nil || scalar.Type != token.StringType {
						t.Fatalf("expected block scalar content token, got %v", scalar)
					}
					if scalar.Position.Line != 3 {
						t.Errorf("content starts on line 3, got %d", scalar.Position.Line)
					}
					wantValue := style.value
					switch style.chomp {
					case "clip":
						if ending.trailing != "" {
							wantValue += "\n"
						}
					case "keep":
						wantValue += ending.trailing
					}
					if scalar.Value != wantValue {
						t.Errorf("value: got %q, want %q", scalar.Value, wantValue)
					}
					wantOrigin := strings.ReplaceAll(body, "\n", newline.value)
					if scalar.Origin != wantOrigin {
						t.Errorf("origin: got %q, want %q", scalar.Origin, wantOrigin)
					}
					if _, err := parser.Parse(tokens, parser.ParseComments); err != nil {
						t.Fatalf("parse tokens: %v", err)
					}
				})
			}
		}
	}
}
