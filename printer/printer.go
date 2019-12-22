package printer

import (
	"fmt"
	"math"
	"strings"

	"github.com/fatih/color"
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
	if len(tokens) == 0 {
		return ""
	}
	if p.LineNumber {
		if p.LineNumberFormat == nil {
			p.LineNumberFormat = defaultLineNumberFormat
		}
	}
	texts := []string{}
	lineNumber := tokens[0].Position.Line
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

const escape = "\x1b"

func format(attr color.Attribute) string {
	return fmt.Sprintf("%s[%dm", escape, attr)
}

func (p *Printer) setDefaultColorSet() {
	p.Bool = func() *Property {
		return &Property{
			Prefix: format(color.FgHiMagenta),
			Suffix: format(color.Reset),
		}
	}
	p.Number = func() *Property {
		return &Property{
			Prefix: format(color.FgHiMagenta),
			Suffix: format(color.Reset),
		}
	}
	p.MapKey = func() *Property {
		return &Property{
			Prefix: format(color.FgHiCyan),
			Suffix: format(color.Reset),
		}
	}
	p.Anchor = func() *Property {
		return &Property{
			Prefix: format(color.FgHiYellow),
			Suffix: format(color.Reset),
		}
	}
	p.Alias = func() *Property {
		return &Property{
			Prefix: format(color.FgHiYellow),
			Suffix: format(color.Reset),
		}
	}
	p.String = func() *Property {
		return &Property{
			Prefix: format(color.FgHiGreen),
			Suffix: format(color.Reset),
		}
	}
}

func (p *Printer) PrintErrorMessage(msg string, isColored bool) string {
	if isColored {
		return fmt.Sprintf("%s%s%s",
			format(color.FgHiRed),
			msg,
			format(color.Reset),
		)
	}
	return msg
}

func (p *Printer) removeLeftSideNewLineChar(src string) string {
	return strings.TrimLeft(strings.TrimLeft(src, "\r"), "\n")
}

func (p *Printer) removeRightSideNewLineChar(src string) string {
	return strings.TrimRight(strings.TrimRight(src, "\r"), "\n")
}

func (p *Printer) removeRightSideWhiteSpaceChar(src string) string {
	return p.removeRightSideNewLineChar(strings.TrimRight(src, " "))
}

func (p *Printer) newLineCount(s string) int {
	src := []rune(s)
	size := len(src)
	cnt := 0
	for i := 0; i < size; i++ {
		c := src[i]
		switch c {
		case '\r':
			if i+1 < size && src[i+1] == '\n' {
				i++
			}
			cnt++
		case '\n':
			cnt++
		}
	}
	return cnt
}

func (p *Printer) isNewLineChar(c byte) bool {
	if c == '\n' {
		return true
	}
	if c == '\r' {
		return true
	}
	return false
}

func (p *Printer) PrintErrorToken(tk *token.Token, isColored bool) string {
	errToken := tk
	pos := tk.Position
	curLine := pos.Line
	curExtLine := curLine + p.newLineCount(p.removeLeftSideNewLineChar(tk.Origin))
	if p.isNewLineChar(tk.Origin[len(tk.Origin)-1]) {
		// if last character is new line character, ignore it.
		curExtLine--
	}
	minLine := int(math.Max(float64(curLine-3), 1))
	maxLine := curExtLine + 3
	for {
		if tk.Position.Line < minLine {
			break
		}
		if tk.Prev == nil {
			break
		}
		tk = tk.Prev
	}
	tokens := token.Tokens{}
	lastTk := tk
	for tk.Position.Line <= curExtLine {
		tokens.Add(tk)
		lastTk = tk
		tk = tk.Next
		if tk == nil {
			break
		}
	}
	org := lastTk.Origin
	trimmed := p.removeRightSideWhiteSpaceChar(org)
	lastTk.Origin = trimmed
	if tk != nil {
		tk.Origin = p.removeLeftSideNewLineChar(tk.Origin)
	}
	p.LineNumber = true
	p.LineNumberFormat = func(num int) string {
		if isColored {
			fn := color.New(color.Bold, color.FgHiWhite).SprintFunc()
			if curLine == num {
				return fn(fmt.Sprintf("> %2d | ", num))
			}
			return fn(fmt.Sprintf("  %2d | ", num))
		}
		if curLine == num {
			return fmt.Sprintf("> %2d | ", num)
		}
		return fmt.Sprintf("  %2d | ", num)
	}
	if isColored {
		p.setDefaultColorSet()
	}
	beforeSource := p.PrintTokens(tokens)
	prefixSpaceNum := len(fmt.Sprintf("  %2d | ", 1))
	annotateLine := strings.Repeat(" ", prefixSpaceNum+errToken.Position.Column-2) + "^"
	tokens = token.Tokens{}
	for tk != nil {
		if tk.Position.Line > maxLine {
			break
		}
		tokens.Add(tk)
		tk = tk.Next
	}
	afterSource := p.PrintTokens(tokens)
	return fmt.Sprintf("%s\n%s\n%s", beforeSource, annotateLine, afterSource)
}
