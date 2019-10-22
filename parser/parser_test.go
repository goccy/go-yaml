package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/printer"
)

func TestParser(t *testing.T) {
	sources := []string{
		"null\n",
		"{}\n",
		"v: hi\n",
		"v: \"true\"\n",
		"v: \"false\"\n",
		"v: true\n",
		"v: false\n",
		"v: 10\n",
		"v: -10\n",
		"v: 42\n",
		"v: 4294967296\n",
		"v: \"10\"\n",
		"v: 0.1\n",
		"v: 0.99\n",
		"v: -0.1\n",
		"v: .inf\n",
		"v: -.inf\n",
		"v: .nan\n",
		"v: null\n",
		"v: \"\"\n",
		"v:\n- A\n- B\n",
		"a: '-'\n",
		"123\n",
		"hello: world\n",
		"a: null\n",
		"v:\n- A\n- 1\n- B:\n  - 2\n  - 3\n",
		"a:\n  b: c\n",
		"a: {x: 1}\n",
		"t2: 2018-01-09T10:40:47Z\nt4: 2098-01-09T10:40:47Z\n",
		"a: [1, 2]\n",
		"a: {b: c, d: e}\n",
		"a: 3s\n",
		"a: <foo>\n",
		"a: \"1:1\"\n",
		"a: 1.2.3.4\n",
		"a: \"2015-02-24T18:19:39Z\"\n",
		"a: 'b: c'\n",
		"a: 'Hello #comment'\n",
		"a: 100.5\n",
		"a: bogus\n",
		"a: \"\\0\"\n",
		"b: 2\na: 1\nd: 4\nc: 3\nsub:\n  e: 5\n",
		"       a       :          b        \n",
		"a: b # comment\nb: c\n",
		"---\na: b\n",
		"a: b\n...\n",
		"%YAML 1.2\n---\n",
		"a: !!binary gIGC\n",
		"a: !!binary |\n  " + strings.Repeat("kJCQ", 17) + "kJ\n  CQ\n",
		"- !tag\n  a: b\n  c: d\n",
		"v:\n- A\n- |-\n  B\n  C\n",
		"v:\n- A\n- >-\n  B\n  C\n",
	}
	var (
		p parser.Parser
	)
	for _, src := range sources {
		fmt.Printf(src)
		tokens := lexer.Tokenize(src)
		var printer printer.Printer
		fmt.Println(printer.PrintTokens(tokens))
		ast, err := p.Parse(tokens)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		fmt.Printf("%+v\n", ast)
	}
}

func TestParseComplicatedDocument(t *testing.T) {
	tests := []struct {
		source string
		expect string
	}{
		{
			`
american:
  - Boston Red Sox
  - Detroit Tigers
  - New York Yankees
national:
  - New York Mets
  - Chicago Cubs
  - Atlanta Braves
`, `
american:
  - Boston Red Sox
  - Detroit Tigers
  - New York Yankees
national:
  - New York Mets
  - Chicago Cubs
  - Atlanta Braves
`,
		},
		{
			`
a:
  b: c
  d: e
  f: g
h:
  i: j
  k:
    l: m
    n: o
  p: q
r: s
`, `
a:
  b: c
  d: e
  f: g
h:
  i: j
  k:
    l: m
    n: o
  p: q
r: s
`,
		},
		{
			`
- a:
  - b
  - c
- d
`, `
- a:
  - b
  - c
- d
`,
		},
		{
			`
- a
- b
- c
 - d
 - e
- f
`, `
- a
- b
- c - d - e
- f
`,
		},
		{`
- a:
   b: c
   d: e
- f:
  g: h
`,
			`
- a:
   b: c
   d: e
- f:
  g: h
`,
		},
		{
			`
a:
 b
 c
d: e
`, `
a: b c
d: e
`,
		},
		{
			`
a
b
c
`, `
a b c
`,
		},
		{
			`
a:
 - b
 - c
`, `
a:
 - b
 - c
`,
		},
		{
			`
-     a     :
      b: c
`, `
- a:
  b: c
`,
		},
		{
			`
- a:
   b
   c
   d
  hoge: fuga
`, `
- a: b c d
  hoge: fuga
`,
		},
		{
			`
- a # ' " # - : %
- b # " # - : % '
- c # # - : % ' "
- d # - : % ' " #
- e # : % ' " # -
- f # % ' : # - :
`,
			`
- a
- b
- c
- d
- e
- f
`,
		},
		{
			`
# comment
a: # comment
# comment
 b: c # comment
 # comment
d: e # comment
# comment
`,
			`
a:
 b: c
d: e
`,
		},
		{
			`
anchored: &anchor foo
aliased: *anchor
`,
			`
anchored: &anchor foo
aliased: *anchor
`,
		},
		{
			`
---
- &CENTER { x: 1, y: 2 }
- &LEFT { x: 0, y: 2 }
- &BIG { r: 10 }
- &SMALL { r: 1 }

# All the following maps are equal:

- # Explicit keys
  x: 1
  y: 2
  r: 10
  label: center/big

- # Merge one map
  << : *CENTER
  r: 10
  label: center/big

- # Merge multiple maps
  << : [ *CENTER, *BIG ]
  label: center/big

- # Override
  << : [ *BIG, *LEFT, *SMALL ]
  x: 1
  label: center/big
`,
			`
- &CENTER {x: 1, y: 2}
- &LEFT {x: 0, y: 2}
- &BIG {r: 10}
- &SMALL {r: 1}
- x: 1
  y: 2
  r: 10
  label: center/big
- <<: *CENTER
  r: 10
  label: center/big
- <<: [*CENTER, *BIG]
  label: center/big
- <<: [*BIG, *LEFT, *SMALL]
  x: 1
  label: center/big
`,
		},
	}
	for _, test := range tests {
		var (
			p parser.Parser
		)
		tokens := lexer.Tokenize(test.source)
		tokens.Dump()
		doc, err := p.Parse(tokens)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		var v Visitor
		for _, node := range doc.Nodes {
			ast.Walk(&v, node)
		}
		expect := fmt.Sprintf("\n%+v\n", doc)
		if test.expect != expect {
			t.Fatalf("unexpected output: [%s] != [%s]", test.expect, expect)
		}
	}
}

func TestSyntaxError(t *testing.T) {
	sources := []string{
		"a:\n- b\n  c: d\n  e: f\n  g: h",
	}
	for _, source := range sources {
		var (
			p parser.Parser
		)
		tokens := lexer.Tokenize(source)
		_, err := p.Parse(tokens)
		if err == nil {
			t.Fatal("cannot catch syntax error")
		}
		fmt.Printf("%+v\n", err)
	}
}

type Visitor struct {
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FlowMappingNode:
		n.Start.Next = nil
		n.Start.Prev = nil
		n.End.Next = nil
		n.End.Prev = nil
	case *ast.FlowSequenceNode:
		n.Start.Next = nil
		n.Start.Prev = nil
		n.End.Next = nil
		n.End.Prev = nil
	default:
		tk := n.GetToken()
		tk.Prev = nil
		tk.Next = nil
	}
	return v
}
