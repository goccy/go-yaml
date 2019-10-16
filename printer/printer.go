package printer

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml/token"
)

type Printer struct {
	LineNumber       bool
	LineNumberFormat func(num int) string
}

func defaultLineNumberFormat(num int) string {
	return fmt.Sprintf("%d | ", num)
}

func (p *Printer) PrintTokens(tokens token.Tokens) string {
	source := ""
	for _, token := range tokens {
		source += token.Origin
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
