package parser_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
)

func TestParser(t *testing.T) {
	sources := []string{
		"",
		"null\n",
		"0_",
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
		"v: !!foo 1",
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
		"\"a\": a\n\"b\": b",
		"'a': a\n'b': b",
		"a: \r\n  b: 1\r\n",
		"a_ok: \r  bc: 2\r",
		"a_mk: \n  bd: 3\n",
		"a: :a",
		"{a: , b: c}",
		"value: >\n",
		"value: >\n\n",
		"value: >\nother:",
		"value: >\n\nother:",
		"a:\n-",
		"a: {foo}",
		"a: {foo,bar}",
		`
{
  a: {
    b: c
  },
  d: e
}
`,
		`
[
  a: {
    b: c
  }]
`,
		`
{
  a: {
    b: c
  }}
`,
		`
- !tag
  a: b
  c: d
`,
		`
a: !tag
  b: c
`,
		`
a: !tag
  b: c
  d: e
`,
		`
a:
  b: c
     
`,
		`
foo: xxx
---
foo: yyy
---
foo: zzz
`,
		`
v:
  a	: 'a'
  bb	: 'a'
`,
		`
v:
  a : 'x'
  b	: 'y'
`,
		`
v:
  a	: 'x'
  b	: 'y'
  c		: 'z'
`,
		`{a: &a c, *a : b}`,
	}
	for idx, src := range sources {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			f, err := parser.Parse(lexer.Tokenize(src), 0)
			if err != nil {
				t.Fatalf("parse error: source [%s]: %+v", src, err)
			}
			_ = f.String() // ensure no panic
		})
	}
}

func TestParseEmptyDocument(t *testing.T) {
	t.Run("empty document", func(t *testing.T) {
		f, err := parser.ParseBytes([]byte(""), parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		got := f.String()
		expected := "\n"
		if got != expected {
			t.Fatalf("failed to parse comment:\nexpected:\n%q\ngot:\n%q", expected, got)
		}
	})

	t.Run("empty document with comment (parse comment = off)", func(t *testing.T) {
		f, err := parser.ParseBytes([]byte("# comment"), 0)
		if err != nil {
			t.Fatal(err)
		}
		got := f.String()
		expected := "\n"
		if got != expected {
			t.Fatalf("failed to parse comment:\nexpected:\n%q\ngot:\n%q", expected, got)
		}
	})

	t.Run("empty document with comment (parse comment = on)", func(t *testing.T) {
		f, err := parser.ParseBytes([]byte("# comment"), parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		got := f.String()
		expected := "# comment\n"
		if got != expected {
			t.Fatalf("failed to parse comment:\nexpected:\n%q\ngot:\n%q", expected, got)
		}
	})
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
v: |
 a
 b
 c`,
			`
v: |
 a
 b
 c
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
		{
			`
- key1: val
  key2:
    (
      foo
      +
      bar
    )
`,
			`
- key1: val
  key2: ( foo + bar )
`,
		},
		{
			`
"a": b
'c': d
"e": "f"
g: "h"
i: 'j'
`,
			`
"a": b
'c': d
"e": "f"
g: "h"
i: 'j'
`,
		},
		{
			`
a:
  - |2
        b
    c: d
`,
			`
a:
  - |2
        b
    c: d
`,
		},
		{
			`
a:
 b: &anchor
 c: &anchor2
d: e
`,
			`
a:
 b: &anchor
 c: &anchor2
d: e
`,
		},
	}

	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			tokens := lexer.Tokenize(test.source)
			f, err := parser.Parse(tokens, 0)
			if err != nil {
				t.Fatalf("%+v", err)
			}
			got := f.String()
			if got != strings.TrimPrefix(test.expect, "\n") {
				t.Fatalf("failed to parse comment:\nexpected:\n%s\ngot:\n%s", strings.TrimPrefix(test.expect, "\n"), got)
			}
			var v Visitor
			for _, doc := range f.Docs {
				ast.Walk(&v, doc.Body)
			}
			expect := fmt.Sprintf("\n%+v", f)
			if test.expect != expect {
				tokens.Dump()
				t.Fatalf("unexpected output: [%s] != [%s]", test.expect, expect)
			}
		})
	}
}

func TestParseWhitespace(t *testing.T) {
	tests := []struct {
		source string
		expect string
	}{
		{
			`
a: b

c: d


e: f
g: h
`,
			`
a: b

c: d

e: f
g: h
`,
		},
		{
			`
a:
  - b: c
    d: e

  - f: g
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
a:
  - b: c
    d: e

  - f: g
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
a:
- b: c
  d: e

- f: g
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
a:
# comment 1
- b: c
  d: e

# comment 2
- f: g
  h: i
`,
			`
a:
# comment 1
- b: c
  d: e

# comment 2
- f: g
  h: i
`,
		},
		{
			`
a:
  # comment 1
  - b: c
    # comment 2
    d: e

  # comment 3
  # comment 4
  - f: g
    h: i # comment 5
`,
			`
a:
  # comment 1
  - b: c
    # comment 2
    d: e

  # comment 3
  # comment 4
  - f: g
    h: i # comment 5
`,
		},
		{
			`
a:
  # comment 1
  - b: c
    # comment 2
    d: e

  # comment 3
  # comment 4
  - f: |
      g
      g
    h: i # comment 5
`,
			`
a:
  # comment 1
  - b: c
    # comment 2
    d: e

  # comment 3
  # comment 4
  - f: |
      g
      g
    h: i # comment 5
`,
		},
		{
			`
a:
  # comment 1
  - b: c
    # comment 2
    d: e

  # comment 3
  # comment 4
  - f: |
      asd
      def

    h: i # comment 5
`,
			`
a:
  # comment 1
  - b: c
    # comment 2
    d: e

  # comment 3
  # comment 4
  - f: |
      asd
      def

    h: i # comment 5
`,
		},
		{
			`
- b: c
  d: e

- f: g
  h: i # comment 4
`,
			`
- b: c
  d: e

- f: g
  h: i # comment 4
`,
		},
		{
			`
a: null
b: null

d: e
`,
			`
a: null
b: null

d: e
`,
		},
		{
			`
foo:
  bar: null # comment

  baz: 1
`,
			`
foo:
  bar: null # comment

  baz: 1
`,
		},
		{
			`
foo:
  bar: null # comment

baz: 1
`,
			`
foo:
  bar: null # comment

baz: 1
`,
		},
		{
			`
{
	"apiVersion": "apps/v1",
	"kind": "Deployment",
	"metadata": {
		"name": "foo",
		"labels": {
			"app": "bar"
		}
	},
	"spec": {
		"replicas": 3,
		"selector": {
			"matchLabels": {
				"app": "bar"
			}
		},
		"template": {
			"metadata": {
				"labels": {
					"app": "bar"
				}
			}
		}
	}
}
`,
			`
{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "foo", "labels": {"app": "bar"}}, "spec": {"replicas": 3, "selector": {"matchLabels": {"app": "bar"}}, "template": {"metadata": {"labels": {"app": "bar"}}}}}
`,
		},
	}

	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			f, err := parser.ParseBytes([]byte(test.source), parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}
			got := f.String()
			if got != strings.TrimPrefix(test.expect, "\n") {
				t.Fatalf("failed to parse comment:\nexpected:\n%s\ngot:\n%s", strings.TrimPrefix(test.expect, "\n"), got)
			}
		})
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
		actual := fmt.Sprintf("%v", ast)
		expect := `a: "a"

b: 1
`
		if expect != actual {
			t.Fatalf("unexpected result\nexpected:\n%s\ngot:\n%s", expect, actual)
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
[1:1] unexpected directive value. document not started
>  1 | %YAML 1.1 {}
       ^
`,
		},
		{
			`{invalid`,
			`
[1:2] could not find flow map content
>  1 | {invalid
        ^
`,
		},
		{
			`{ "key": "value" `,
			`
[1:1] could not find flow mapping end token '}'
>  1 | { "key": "value"
       ^
`,
		},
		{
			`
a: |invalidopt
  foo
`,
			`
[2:4] invalid header option: "invalidopt"
>  2 | a: |invalidopt
          ^
   3 |   foo`,
		},
		{
			`
a: 1
b
`,
			`
[3:1] non-map value is specified
   2 | a: 1
>  3 | b
       ^
`,
		},
		{
			`
a: 'b'
  c: d
`,
			`
[3:3] value is not allowed in this context. map key-value is pre-defined
   2 | a: 'b'
>  3 |   c: d
         ^
`,
		},
		{
			`
a: 'b'
  - c
`,
			`
[3:3] value is not allowed in this context. map key-value is pre-defined
   2 | a: 'b'
>  3 |   - c
         ^
`,
		},
		{
			`
a: 'b'
  # comment
  - c
`,
			`
[4:3] value is not allowed in this context. map key-value is pre-defined
   2 | a: 'b'
   3 |   # comment
>  4 |   - c
         ^
`,
		},
		{
			`
a: 1
b
- c
`,
			`
[3:1] non-map value is specified
   2 | a: 1
>  3 | b
       ^
   4 | - c`,
		},
		{
			`a: [`,
			`
[1:4] sequence end token ']' not found
>  1 | a: [
          ^
`,
		},
		{
			`a: ]`,
			`
[1:4] could not find '[' character corresponding to ']'
>  1 | a: ]
          ^
`,
		},
		{
			`a: [ [1] [2] [3] ]`,
			`
[1:10] ',' or ']' must be specified
>  1 | a: [ [1] [2] [3] ]
                ^
`,
		},
		{
			`
a: -
b: -
`,
			`
[2:4] block sequence entries are not allowed in this context
>  2 | a: -
          ^
   3 | b: -`,
		},
		{
			`
a: - 1
b: - 2
`,
			`
[2:4] block sequence entries are not allowed in this context
>  2 | a: - 1
          ^
   3 | b: - 2`,
		},
		{
			`a: 'foobarbaz`,
			`
[1:4] could not find end character of single-quoted text
>  1 | a: 'foobarbaz
          ^
`,
		},
		{
			`a: "\"key\": \"value:\"`,
			`
[1:4] could not find end character of double-quoted text
>  1 | a: "\"key\": \"value:\"
          ^
`,
		},
		{
			`foo: [${should not be allowed}]`,
			`
[1:8] ',' or ']' must be specified
>  1 | foo: [${should not be allowed}]
              ^
`,
		},
		{
			`foo: [$[should not be allowed]]`,
			`
[1:8] ',' or ']' must be specified
>  1 | foo: [$[should not be allowed]]
              ^
`,
		},
		{
			">\n>",
			`
[2:1] could not find multi-line content
   1 | >
>  2 | >
       ^
`,
		},
		{
			">\n1",
			`
[2:1] could not find multi-line content
   1 | >
>  2 | 1
       ^
`,
		},
		{
			"|\n1",
			`
[2:1] could not find multi-line content
   1 | |
>  2 | 1
       ^
`,
		},
		{
			"a: >3\n  1",
			`
[2:3] invalid number of indent is specified in the multi-line header
   1 | a: >3
>  2 |   1
         ^
`,
		},
		{
			`
a:
  - |
        b
    c: d
`,
			`
[5:5] value is not allowed in this context
   2 | a:
   3 |   - |
   4 |         b
>  5 |     c: d
           ^
`,
		},
		{
			`
a:
  - |
        b
    c:
      d: e
`,
			`
[5:5] value is not allowed in this context
   2 | a:
   3 |   - |
   4 |         b
>  5 |     c:
           ^
   6 |       d: e`,
		},
		{
			"key: [@val]",
			`
[1:7] '@' is a reserved character
>  1 | key: [@val]
             ^
`,
		},
		{
			"key: [`val]",
			"\n[1:7] '`' is a reserved character\n>  1 | key: [`val]\n             ^\n",
		},
		{
			`{a: b}: v`,
			`
[1:7] found an invalid key for this map
>  1 | {a: b}: v
             ^
`,
		},
		{
			`[a]: v`,
			`
[1:4] found an invalid key for this map
>  1 | [a]: v
          ^
`,
		},
		{
			`
foo:
  bar:
    foo: 2
  baz:
    foo: 3
foo: 2
`,
			`
[7:1] mapping key "foo" already defined at [2:1]
   4 |     foo: 2
   5 |   baz:
   6 |     foo: 3
>  7 | foo: 2
       ^
`,
		},
		{
			`
foo:
  bar:
    foo: 2
  baz:
    foo: 3
    foo: 4
`,
			`
[7:5] mapping key "foo" already defined at [6:5]
   4 |     foo: 2
   5 |   baz:
   6 |     foo: 3
>  7 |     foo: 4
           ^
`,
		},
		{
			`{"000":0000A,`,
			`
[1:13] could not find flow map content
>  1 | {"000":0000A,
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
		name     string
		yaml     string
		expected string
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
		{
			name: "multiline",
			yaml: `
# foo comment
# foo comment2
foo: # map key comment
  # bar above comment
  # bar above comment2
  bar: 10 # comment for bar
  # baz above comment
  # baz above comment2
  baz: bbbb # comment for baz
  piyo: # sequence key comment
  # sequence1 above comment 1
  # sequence1 above comment 2
  - sequence1 # sequence1
  # sequence2 above comment 1
  # sequence2 above comment 2
  - sequence2 # sequence2
  # sequence3 above comment 1
  # sequence3 above comment 2
  - false # sequence3
# foo2 comment
# foo2 comment2
foo2: &anchor text # anchor comment
# foo3 comment
# foo3 comment2
foo3: *anchor # alias comment
`,
		},
		{
			name: "literal",
			yaml: `
foo: | # comment
  x: 42
`,
		},
		{
			name: "folded",
			yaml: `
foo: > # comment
  x: 42
`,
		},
		{
			name: "unattached comment",
			yaml: `
# This comment is in its own document
---
a: b
`,
		},
		{
			name: "map with misaligned indentation in comments",
			yaml: `
 # commentA
a:  #commentB
   # commentC
  b: c  # commentD
    # commentE
  d: e  # commentF
 # commentG
`,
			expected: `
# commentA
a: #commentB
  # commentC
  b: c # commentD
  # commentE
  d: e # commentF
# commentG
`,
		},
		{
			name: "sequence with misaligned indentation in comments",
			yaml: `
 # commentA
- a  # commentB
 # commentC
- b:   # commentD
   # commentE
  - d  # commentF
    # commentG
  - e  # commentG
 # commentH
`,
			expected: `
# commentA
- a # commentB
# commentC
- b: # commentD
  # commentE
  - d # commentF
  # commentG
  - e # commentG
# commentH
`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := parser.ParseBytes([]byte(test.yaml), parser.ParseComments)
			if err != nil {
				t.Fatalf("%+v", err)
			}
			got := "\n" + f.String()
			expected := test.yaml
			if test.expected != "" {
				expected = test.expected
			}
			if expected != got {
				t.Fatalf("expected:%s\ngot:%s", expected, got)
			}
		})
	}
}

func TestCommentWithNull(t *testing.T) {
	t.Run("same line", func(t *testing.T) {
		content := `
foo:
  bar: # comment
  baz: 1
`
		expected := `
foo:
  bar: # comment
  baz: 1`
		f, err := parser.ParseBytes([]byte(content), parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		if len(f.Docs) != 1 {
			t.Fatal("failed to parse content with same line comment")
		}
		got := f.Docs[0].String()
		if got != strings.TrimPrefix(expected, "\n") {
			t.Fatalf("failed to parse comment:\nexpected:\n%s\ngot:\n%s", strings.TrimPrefix(expected, "\n"), got)
		}
	})
	t.Run("next line", func(t *testing.T) {
		content := `
foo:
  bar:
    # comment
  baz: 1
`
		expected := `
foo:
  bar:
  # comment
  baz: 1`
		f, err := parser.ParseBytes([]byte(content), parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		if len(f.Docs) != 1 {
			t.Fatal("failed to parse content with next line comment")
		}
		got := f.Docs[0].String()
		if got != strings.TrimPrefix(expected, "\n") {
			t.Fatalf("failed to parse comment:\nexpected:\n%s\ngot:\n%s", strings.TrimPrefix(expected, "\n"), got)
		}
	})
	t.Run("next line and different indent", func(t *testing.T) {
		content := `
foo:
  bar:
 # comment
baz: 1`
		f, err := parser.ParseBytes([]byte(content), parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		if len(f.Docs) != 1 {
			t.Fatal("failed to parse content with next line comment")
		}
		expected := `
foo:
  bar:
# comment
baz: 1`
		got := f.Docs[0].String()
		if got != strings.TrimPrefix(expected, "\n") {
			t.Fatalf("failed to parse comment:\nexpected:\n%s\ngot:\n%s", strings.TrimPrefix(expected, "\n"), got)
		}
	})
}

func TestSequenceComment(t *testing.T) {
	content := `
foo:
  - # comment
    bar: 1
baz:
  - xxx
`
	f, err := parser.ParseBytes([]byte(content), parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	if len(f.Docs) != 1 {
		t.Fatal("failed to parse content with next line with sequence")
	}
	expected := `
foo:
  # comment
  - bar: 1
baz:
  - xxx`
	got := f.Docs[0].String()
	if got != strings.TrimPrefix(expected, "\n") {
		t.Fatalf("failed to parse comment:\nexpected:\n%s\ngot:\n%s", strings.TrimPrefix(expected, "\n"), got)
	}
	t.Run("foo[0].bar", func(t *testing.T) {
		path, err := yaml.PathString("$.foo[0].bar")
		if err != nil {
			t.Fatal(err)
		}
		v, err := path.FilterFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if v.String() != "1" {
			t.Fatal("failed to get foo[0].bar value")
		}
	})
	t.Run("baz[0]", func(t *testing.T) {
		path, err := yaml.PathString("$.baz[0]")
		if err != nil {
			t.Fatal(err)
		}
		v, err := path.FilterFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if v.String() != "xxx" {
			t.Fatal("failed to get baz[0] value")
		}
	})
}

func TestCommentWithMap(t *testing.T) {
	yml := `
single:
  # foo comment
  foo: bar

multiple:
    # a comment
    a: b
    # c comment
    c: d
`

	file, err := parser.ParseBytes([]byte(yml), parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	if len(file.Docs) == 0 {
		t.Fatal("cannot get file docs")
	}
	if file.Docs[0].Body == nil {
		t.Fatal("cannot get docs body")
	}
	mapNode, ok := file.Docs[0].Body.(*ast.MappingNode)
	if !ok {
		t.Fatalf("failed to get map node. got: %T\n", file.Docs[0].Body)
	}
	if len(mapNode.Values) != 2 {
		t.Fatalf("failed to get map values. got %d", len(mapNode.Values))
	}

	singleNode, ok := mapNode.Values[0].Value.(*ast.MappingNode)
	if !ok {
		t.Fatalf("failed to get single node. got %T", mapNode.Values[0].Value)
	}
	if len(singleNode.Values) != 1 {
		t.Fatalf("failed to get single node values. got %d", len(singleNode.Values))
	}
	if singleNode.Values[0].GetComment().GetToken().Value != " foo comment" {
		t.Fatalf("failed to get comment from single. got %q", singleNode.GetComment().GetToken().Value)
	}

	multiNode, ok := mapNode.Values[1].Value.(*ast.MappingNode)
	if !ok {
		t.Fatalf("failed to get multiple node. got: %T", mapNode.Values[1])
	}
	if multiNode.GetComment() != nil {
		t.Fatalf("found unexpected comment")
	}
	if len(multiNode.Values) != 2 {
		t.Fatalf("failed to get multiple node values. got %d", len(multiNode.Values))
	}
	if multiNode.Values[0].GetComment().GetToken().Value != " a comment" {
		t.Fatalf("failed to get comment from multiple[0]. got %q", multiNode.Values[0].GetComment().GetToken().Value)
	}
	if multiNode.Values[1].GetComment().GetToken().Value != " c comment" {
		t.Fatalf("failed to get comment from multiple[1]. got %q", multiNode.Values[1].GetComment().GetToken().Value)
	}
}

func TestInFlowStyle(t *testing.T) {
	type inFlowStyle interface {
		SetIsFlowStyle(bool)
	}

	tests := []struct {
		source string
		expect string
	}{
		{
			`
  - foo
  - bar
  - baz
`,
			`[foo, bar, baz]
`,
		},
		{
			`
foo: bar
baz: fizz
`,
			`{foo: bar, baz: fizz}
`,
		},
		{
			`
foo:
  - bar
  - baz
  - fizz: buzz
`,
			`{foo: [bar, baz, {fizz: buzz}]}
`,
		},
	}

	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			f, err := parser.ParseBytes([]byte(test.source), parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}

			if len(f.Docs) != 1 {
				t.Fatal("failed to parse content")
			}

			ifs, ok := f.Docs[0].Body.(inFlowStyle)
			if !ok {
				t.Fatalf("failed to get inFlowStyle. got: %T", f.Docs[0].Body)
			}
			ifs.SetIsFlowStyle(true)

			got := f.String()

			if got != test.expect {
				t.Fatalf("failed to parse comment:\nexpected:\n%s\ngot:\n%s", test.expect, got)
			}
		})
	}
}

func TestNodePath(t *testing.T) {
	yml := `
a: # commentA
  b: # commentB
    c: foo # commentC
    d: bar # commentD
    e: baz # commentE
  f: # commentF
    g: hoge # commentG
  h: # commentH
   - list1 # comment list1
   - list2 # comment list2
   - list3 # comment list3
  i: fuga # commentI
j: piyo # commentJ
k.l.m.n: moge # commentKLMN
o#p: hogera # commentOP
q#.r: hogehoge # commentQR
`
	f, err := parser.ParseBytes([]byte(yml), parser.ParseComments)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	var capturer pathCapturer
	for _, doc := range f.Docs {
		ast.Walk(&capturer, doc.Body)
	}
	commentPaths := []string{}
	for i := 0; i < capturer.capturedNum; i++ {
		if capturer.orderedTypes[i] == ast.CommentType {
			commentPaths = append(commentPaths, capturer.orderedPaths[i])
		}
	}
	expectedPaths := []string{
		"$.a",
		"$.a.b",
		"$.a.b.c",
		"$.a.b.d",
		"$.a.b.e",
		"$.a.f",
		"$.a.f.g",
		"$.a.h",
		"$.a.h[0]",
		"$.a.h[1]",
		"$.a.h[2]",
		"$.a.i",
		"$.j",
		"$.'k.l.m.n'",
		"$.o#p",
		"$.'q#.r'",
	}
	if !reflect.DeepEqual(expectedPaths, commentPaths) {
		t.Fatalf("failed to get YAMLPath to the comment node:\nexpected[%s]\ngot     [%s]", expectedPaths, commentPaths)
	}
}

type pathCapturer struct {
	capturedNum   int
	orderedPaths  []string
	orderedTypes  []ast.NodeType
	orderedTokens []*token.Token
}

func (c *pathCapturer) Visit(node ast.Node) ast.Visitor {
	c.capturedNum++
	c.orderedPaths = append(c.orderedPaths, node.GetPath())
	c.orderedTypes = append(c.orderedTypes, node.Type())
	c.orderedTokens = append(c.orderedTokens, node.GetToken())
	return c
}

type Visitor struct{}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	tk := node.GetToken()
	tk.Prev = nil
	tk.Next = nil
	return v
}
