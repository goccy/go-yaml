package printer

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml/token"
)

type PrintFunc func(string) string

type Printer struct {
	LineNumber       bool
	LineNumberFormat func(num int) string
	MapKey           PrintFunc
	Anchor           PrintFunc
	Alias            PrintFunc
	Bool             PrintFunc
	String           PrintFunc
}

func defaultLineNumberFormat(num int) string {
	return fmt.Sprintf("%2d | ", num)
}

func (p *Printer) PrintTokens(tokens token.Tokens) string {
	source := ""
	for _, tk := range tokens {
		switch tk.PreviousType() {
		case token.AnchorType:
			if p.Anchor != nil {
				source += p.Anchor(tk.Origin)
			} else {
				source += tk.Origin
			}
			continue
		case token.AliasType:
			if p.Alias != nil {
				source += p.Alias(tk.Origin)
			} else {
				source += tk.Origin
			}
			continue
		}
		switch tk.NextType() {
		case token.MappingValueType:
			if p.MapKey != nil {
				source += p.MapKey(tk.Origin)
			} else {
				source += tk.Origin
			}
			continue
		}
		switch tk.Type {
		case token.BoolType:
			if p.Bool != nil {
				source += p.Bool(tk.Origin)
			} else {
				source += tk.Origin
			}
		case token.AnchorType:
			if p.Anchor != nil {
				source += p.Anchor(tk.Origin)
			} else {
				source += tk.Origin
			}
		case token.AliasType:
			if p.Anchor != nil {
				source += p.Alias(tk.Origin)
			} else {
				source += tk.Origin
			}
		case token.StringType:
			if p.String != nil {
				source += p.String(tk.Origin)
			} else {
				source += tk.Origin
			}
		default:
			source += tk.Origin
		}
	}
	if p.LineNumber {
		if p.LineNumberFormat == nil {
			p.LineNumberFormat = defaultLineNumberFormat
		}
		codes := []string{}
		for idx, src := range strings.Split(source, "\n") {
			codes = append(codes, fmt.Sprintf("%s%s", p.LineNumberFormat(idx+1), src))
		}
		return strings.Join(codes, "\n")
	}
	return source
}
