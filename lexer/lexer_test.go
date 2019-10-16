package lexer_test

import (
	"strings"
	"testing"

	"github.com/goccy/go-yaml/lexer"
)

func TestTokenize(t *testing.T) {
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
		"v:\n- A\n- |-\n  B\n  C\n",
		"v:\n- A\n- 1\n- B:\n  - 2\n  - 3\n",
		"a:\n  b: c\n",
		"a: '-'\n",
		"123\n",
		"hello: world\n",
		"a: null\n",
		"a: {x: 1}\n",
		"a: [1, 2]\n",
		"t2: 2018-01-09T10:40:47Z\nt4: 2098-01-09T10:40:47Z\n",
		"a: {b: c, d: e}\n",
		"a: 3s\n",
		"a: <foo>\n",
		"a: \"1:1\"\n",
		"a: \"\\0\"\n",
		"a: !!binary gIGC\n",
		"a: !!binary |\n  " + strings.Repeat("kJCQ", 17) + "kJ\n  CQ\n",
		"b: 2\na: 1\nd: 4\nc: 3\nsub:\n  e: 5\n",
		"a: 1.2.3.4\n",
		"a: \"2015-02-24T18:19:39Z\"\n",
		"a: 'b: c'\n",
		"a: 'Hello #comment'\n",
		"a: 100.5\n",
		"a: bogus\n",
	}
	var lex lexer.Lexer
	for _, src := range sources {
		lex.Tokenize(src).Dump()
	}
}
