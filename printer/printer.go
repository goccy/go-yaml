package printer

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/token"
)

// Property additional property set for each the token
type Property struct {
	Prefix string
	Suffix string
}

// PrintFunc returns property instance
type PrintFunc func() *Property

// Printer create text from token collection or ast
type Printer struct {
	LineNumber       bool
	LineNumberFormat func(num int) string
	MapKey           PrintFunc
	Anchor           PrintFunc
	Alias            PrintFunc
	Bool             PrintFunc
	String           PrintFunc
	Number           PrintFunc
}

func defaultLineNumberFormat(num int) string {
	return fmt.Sprintf("%2d | ", num)
}

func (p *Printer) lineTexts(tk *token.Token, prop *Property) []string {
	texts := []string{}
	for idx, src := range strings.Split(tk.Origin, "\n") {
		header := ""
		if p.LineNumber {
			header = p.LineNumberFormat(tk.Position.Line + idx)
		}
		lineText := prop.Prefix + src + prop.Suffix
		texts = append(texts, fmt.Sprintf("%s%s", header, lineText))
	}
	return texts
}

func (p *Printer) property(tk *token.Token) *Property {
	prop := &Property{}
	switch tk.PreviousType() {
	case token.AnchorType:
		if p.Anchor != nil {
			return p.Anchor()
		}
		return prop
	case token.AliasType:
		if p.Alias != nil {
			return p.Alias()
		}
		return prop
	}
	switch tk.NextType() {
	case token.MappingValueType:
		if p.MapKey != nil {
			return p.MapKey()
		}
		return prop
	}
	switch tk.Type {
	case token.BoolType:
		if p.Bool != nil {
			return p.Bool()
		}
		return prop
	case token.AnchorType:
		if p.Anchor != nil {
			return p.Anchor()
		}
		return prop
	case token.AliasType:
		if p.Anchor != nil {
			return p.Alias()
		}
		return prop
	case token.StringType, token.SingleQuoteType, token.DoubleQuoteType:
		if p.String != nil {
			return p.String()
		}
		return prop
	case token.IntegerType, token.FloatType:
		if p.Number != nil {
			return p.Number()
		}
		return prop
	default:
	}
	return prop
}

// PrintTokens create text from token collection
func (p *Printer) PrintTokens(tokens token.Tokens) string {
	if p.LineNumber {
		if p.LineNumberFormat == nil {
			p.LineNumberFormat = defaultLineNumberFormat
		}
	}
	texts := []string{}
	lineNumber := 1
	for _, tk := range tokens {
		lines := strings.Split(tk.Origin, "\n")
		prop := p.property(tk)
		header := ""
		if p.LineNumber {
			header = p.LineNumberFormat(lineNumber)
		}
		if len(lines) == 1 {
			line := prop.Prefix + lines[0] + prop.Suffix
			if len(texts) == 0 {
				texts = append(texts, header+line)
				lineNumber++
			} else {
				text := texts[len(texts)-1]
				texts[len(texts)-1] = text + line
			}
		} else {
			for idx, src := range lines {
				if p.LineNumber {
					header = p.LineNumberFormat(lineNumber)
				}
				line := prop.Prefix + src + prop.Suffix
				if idx == 0 {
					if len(texts) == 0 {
						texts = append(texts, header+line)
						lineNumber++
					} else {
						text := texts[len(texts)-1]
						texts[len(texts)-1] = text + line
					}
				} else {
					texts = append(texts, fmt.Sprintf("%s%s", header, line))
					lineNumber++
				}
			}
		}
	}
	return strings.Join(texts, "\n")
}

// PrintNode create text from ast.Node
func (p *Printer) PrintNode(node ast.Node) []byte {
	return []byte(fmt.Sprintf("%+v\n", node))
}
