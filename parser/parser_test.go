package parser_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pulumi/go-yaml/ast"
	"github.com/pulumi/go-yaml/lexer"
	"github.com/pulumi/go-yaml/parser"
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
		"a: abc <<def>> ghi",
		"a: <<abcd",
		"a: <<:abcd",
		"a: <<  :abcd",
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
		"v: |-\n  0\n",
		"v: |-\n  0\nx: 0",
		`"a\n1\nb"`,
		`{"a":"b"}`,
		`!!map {
  ? !!str "explicit":!!str "entry",
  ? !!str "implicit" : !!str "entry",
  ? !!null "" : !!null "",
}`,
	}
	for _, src := range sources {
		if _, err := parser.Parse(lexer.Tokenize(src), 0); err != nil {
			t.Fatalf("parse error: source [%s]: %+v", src, err)
		}
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
		{
			`
a: 0 - 1
`,
			`
a: 0 - 1
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
- f: null
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
- a: null
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
a: b#notcomment
`,
			`
a: b#notcomment
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
---
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
		{
			`
a:
- - b
- - c
  - d
`,
			`
a:
- - b
- - c
  - d
`,
		},
		{
			`
a:
  b:
    c: d
  e:
    f: g
    h: i
j: k
`,
			`
a:
  b:
    c: d
  e:
    f: g
    h: i
j: k
`,
		},
		{
			`
---
a: 1
b: 2
...
---
c: 3
d: 4
...
`,
			`
---
a: 1
b: 2
...
---
c: 3
d: 4
...
`,
		},
		{
			`
a:
  b: |
    {
      [ 1, 2 ]
    }
  c: d
`,
			`
a:
  b: |
    {
      [ 1, 2 ]
    }
  c: d
`,
		},
		{
			`
|
    hoge
    fuga
    piyo`,
			`
|
    hoge
    fuga
    piyo
`,
		},
		{
			`
a: |
   bbbbbbb


   ccccccc
d: eeeeeeeeeeeeeeeee
`,
			`
a: |
   bbbbbbb


   ccccccc
d: eeeeeeeeeeeeeeeee
`,
		},
		{
			`
a: b    
  c
`,
			`
a: b c
`,
		},
		{
			`
a:    
  b: c
`,
			`
a:
  b: c
`,
		},
		{
			`
a: b    
c: d
`,
			`
a: b
c: d
`,
		},
		{
			`
- ab - cd
- ef - gh
`,
			`
- ab - cd
- ef - gh
`,
		},
		{
			`
- 0 - 1
 - 2 - 3
`,
			`
- 0 - 1 - 2 - 3
`,
		},
		{
			`
a - b - c: value
`,
			`
a - b - c: value
`,
		},
		{
			`
a:
-
  b: c
  d: e
-
  f: g
  h: i
`,
			`
a:
- b: c
  d: e
- f: g
  h: i
`,
		},
		{
			`
a: |-
  value
b: c
`,
			`
a: |-
  value
b: c
`,
		},
		{
			`
a:  |+
  value
b: c
`,
			`
a: |+
  value
b: c
`,
		},
	}

	for _, test := range tests {
		tokens := lexer.Tokenize(test.source)
		f, err := parser.Parse(tokens, 0)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		var v Visitor
		for _, doc := range f.Docs {
			ast.Walk(&v, doc.Body)
		}
		expect := fmt.Sprintf("\n%+v\n", f)
		if test.expect != expect {
			tokens.Dump()
			t.Fatalf("unexpected output: [%s] != [%s]", test.expect, expect)
		}
	}
}

func TestNewLineChar(t *testing.T) {
	for _, f := range []string{
		"lf.yml",
		"cr.yml",
		"crlf.yml",
	} {
		ast, err := parser.ParseFile(filepath.Join("testdata", f), 0)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		actual := fmt.Sprintf("%v\n", ast)
		expect := `a: "a"
b: 1
`
		if expect != actual {
			t.Fatal("unexpected result")
		}
	}
}

func TestSyntaxError(t *testing.T) {
	tests := []struct {
		source string
		expect string
	}{
		{
			`
a:
- b
  c: d
  e: f
  g: h`,
			`
[3:3] unexpected key name
   2 | a:
>  3 | - b
   4 |   c: d
         ^
   5 |   e: f
   6 |   g: h`,
		},
		{
			`
a
- b: c`,
			`
[2:1] unexpected key name
>  2 | a
   3 | - b: c
       ^
`,
		},
		{
			`%YAML 1.1 {}`,
			`
[1:2] unexpected directive value. document not started
>  1 | %YAML 1.1 {}
        ^
`,
		},
		{
			`{invalid`,
			`
[1:2] unexpected map
>  1 | {invalid
        ^
`,
		},
		{
			`{ "key": "value" `,
			`
[1:1] unterminated flow mapping
>  1 | { "key": "value"
       ^
`,
		},
	}
	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			_, err := parser.ParseBytes([]byte(test.source), 0)
			if err == nil {
				t.Fatal("cannot catch syntax error")
			}
			actual := "\n" + err.Error()
			if test.expect != actual {
				t.Fatalf("expected: [%s] but got [%s]", test.expect, actual)
			}
		})
	}
}

func TestComment(t *testing.T) {
	tests := []struct {
		name string
		yaml string
	}{
		{
			name: "map with comment",
			yaml: `
# commentA
a: #commentB
  # commentC
  b: c # commentD
  # commentE
  d: e # commentF
  # commentG
  f: g # commentH
# commentI
f: g # commentJ
# commentK
`,
		},
		{
			name: "sequence with comment",
			yaml: `
# commentA
- a # commentB
# commentC
- b: # commentD
  # commentE
  - d # commentF
  - e # commentG
# commentH
`,
		},
		{
			name: "anchor and alias",
			yaml: `
a: &x b # commentA
c: *x # commentB
`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := parser.ParseBytes([]byte(test.yaml), parser.ParseComments)
			if err != nil {
				t.Fatalf("%+v", err)
			}
			var v Visitor
			for _, doc := range f.Docs {
				ast.Walk(&v, doc.Body)
			}
		})
	}
}

type Visitor struct {
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	tk := node.GetToken()
	tk.Prev = nil
	tk.Next = nil
	if comment := node.GetComment(); comment != nil {
		comment.Prev = nil
		comment.Next = nil
	}
	return v
}
