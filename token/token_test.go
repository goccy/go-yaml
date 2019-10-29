package token_test

import (
	"testing"

	"github.com/goccy/go-yaml/token"
)

func TestToken(t *testing.T) {
	pos := &token.Position{}
	tokens := token.Tokens{
		token.SequenceEntry("-", pos),
		token.MappingKey(pos),
		token.MappingValue(pos),
		token.CollectEntry(",", pos),
		token.SequenceStart("[", pos),
		token.SequenceEnd("]", pos),
		token.MappingStart("{", pos),
		token.MappingEnd("}", pos),
		token.Comment("#", "#", pos),
		token.Anchor("&", pos),
		token.Alias("*", pos),
		token.Literal("|", "|", pos),
		token.Folded(">", ">", pos),
		token.SingleQuote("'", "'", pos),
		token.DoubleQuote(`"`, `"`, pos),
		token.Directive(pos),
		token.Space(pos),
		token.MergeKey("<<", pos),
		token.DocumentHeader(pos),
		token.DocumentEnd(pos),
		token.New("1", "1", pos),
		token.New("3.14", "3.14", pos),
		token.New("-0b101010", "-0b101010", pos),
		token.New("0xA", "0xA", pos),
		token.New("685.230_15e+03", "685.230_15e+03", pos),
		token.New("02472256", "02472256", pos),
		token.New("0o2472256", "0o2472256", pos),
		token.New("", "", pos),
		token.New("_1", "_1", pos),
		token.New("1.1.1.1", "1.1.1.1", pos),
		token.New("+", "+", pos),
		token.New("-", "-", pos),
		token.New("_", "_", pos),
		token.New("true", "true", pos),
		token.New("false", "false", pos),
		token.New(".nan", ".nan", pos),
		token.New(".inf", ".inf", pos),
		token.New("-.inf", "-.inf", pos),
		token.New("null", "null", pos),
		token.Tag("!!null", "!!null", pos),
		token.Tag("!!map", "!!map", pos),
		token.Tag("!!str", "!!str", pos),
		token.Tag("!!seq", "!!seq", pos),
		token.Tag("!!binary", "!!binary", pos),
		token.Tag("!!omap", "!!omap", pos),
		token.Tag("!!set", "!!set", pos),
		token.Tag("!!int", "!!int", pos),
		token.Tag("!!float", "!!float", pos),
		token.Tag("!hoge", "!hoge", pos),
	}
	tokens.Dump()
	tokens.Add(token.New("hoge", "hoge", pos))
	if tokens[len(tokens)-1].PreviousType() != token.TagType {
		t.Fatal("invalid previous token type")
	}
	if tokens[0].PreviousType() != token.UnknownType {
		t.Fatal("invalid previous token type")
	}
	if tokens[len(tokens)-2].NextType() != token.StringType {
		t.Fatal("invalid next token type")
	}
	if tokens[len(tokens)-1].NextType() != token.UnknownType {
		t.Fatal("invalid next token type")
	}
}

func TestIsNeedQuoted(t *testing.T) {
	if !token.IsNeedQuoted("") {
		t.Fatal("failed to quoted judge for empty string")
	}
	if !token.IsNeedQuoted("true") {
		t.Fatal("failed to quoted judge for boolean")
	}
	if !token.IsNeedQuoted("1.234") {
		t.Fatal("failed to quoted judge for number")
	}
	if !token.IsNeedQuoted("1:1") {
		t.Fatal("failed to quoted judge for time")
	}
	if !token.IsNeedQuoted("hoge # comment") {
		t.Fatal("failed to quoted judge for comment")
	}
	if !token.IsNeedQuoted("\\0") {
		t.Fatal("failed to quoted judge for escaped token")
	}
	if token.IsNeedQuoted("Hello World") {
		t.Fatal("failed to unquoted judge")
	}
}
