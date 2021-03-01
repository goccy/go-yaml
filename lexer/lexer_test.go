package lexer_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/token"
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
	for _, src := range sources {
		lexer.Tokenize(src).Dump()
	}
}

type testToken struct {
	line   int
	column int
	value  string
}

func TestSingleLineToken_ValueLineColumnPosition(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		expect map[int]string // Column -> Value map.
	}{
		{
			name: "single quote, single value array",
			src:  "test: ['test']",
			expect: map[int]string{
				1:  "test",
				5:  ":",
				7:  "[",
				8:  "test",
				14: "]",
			},
		},
		{
			name: "double quote, single value array",
			src:  `test: ["test"]`,
			expect: map[int]string{
				1:  "test",
				5:  ":",
				7:  "[",
				8:  "test",
				14: "]",
			},
		},
		{
			name: "no quotes, single value array",
			src:  "test: [somevalue]",
			expect: map[int]string{
				1:  "test",
				5:  ":",
				7:  "[",
				8:  "somevalue",
				17: "]",
			},
		},
		{
			name: "single quote, multi value array",
			src:  "myarr: ['1','2','3', '444' , '55','66' ,  '77'  ]",
			expect: map[int]string{
				1:  "myarr",
				6:  ":",
				8:  "[",
				9:  "1",
				12: ",",
				13: "2",
				16: ",",
				17: "3",
				20: ",",
				22: "444",
				28: ",",
				30: "55",
				34: ",",
				35: "66",
				40: ",",
				43: "77",
				49: "]",
			},
		},
		{
			name: "double quote, multi value array",
			src:  `myarr: ["1","2","3", "444" , "55","66" ,  "77"  ]`,
			expect: map[int]string{
				1:  "myarr",
				6:  ":",
				8:  "[",
				9:  "1",
				12: ",",
				13: "2",
				16: ",",
				17: "3",
				20: ",",
				22: "444",
				28: ",",
				30: "55",
				34: ",",
				35: "66",
				40: ",",
				43: "77",
				49: "]",
			},
		},
		{
			name: "no quote, multi value array",
			src:  "numbers: [1, 5, 99,100, 3, 7 ]",
			expect: map[int]string{
				1:  "numbers",
				8:  ":",
				10: "[",
				11: "1",
				12: ",",
				14: "5",
				15: ",",
				17: "99",
				19: ",",
				20: "100",
				23: ",",
				25: "3",
				26: ",",
				28: "7",
				30: "]",
			},
		},
		{
			name: "double quotes, nested arrays",
			src:  `Strings: ["1",["2",["3"]]]`,
			expect: map[int]string{
				1:  "Strings",
				8:  ":",
				10: "[",
				11: "1",
				14: ",",
				15: "[",
				16: "2",
				19: ",",
				20: "[",
				21: "3",
				24: "]",
				25: "]",
				26: "]",
			},
		},
		{
			name: "mixed quotes, nested arrays",
			src:  `Values: [1,['2',"3",4,["5",6]]]`,
			expect: map[int]string{
				1:  "Values",
				7:  ":",
				9:  "[",
				10: "1",
				11: ",",
				12: "[",
				13: "2",
				16: ",",
				17: "3",
				20: ",",
				21: "4",
				22: ",",
				23: "[",
				24: "5",
				27: ",",
				28: "6",
				29: "]",
				30: "]",
				31: "]",
			},
		},
		{
			name: "double quote, empty array",
			src:  `Empty: ["", ""]`,
			expect: map[int]string{
				1:  "Empty",
				6:  ":",
				8:  "[",
				9:  "",
				11: ",",
				13: "",
				15: "]",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := lexer.Tokenize(tc.src)
			sort.Slice(got, func(i, j int) bool {
				return got[i].Position.Column < got[j].Position.Column
			})
			var expected []testToken
			for k, v := range tc.expect {
				tt := testToken{
					line:   1,
					column: k,
					value:  v,
				}
				expected = append(expected, tt)
			}
			sort.Slice(expected, func(i, j int) bool {
				return expected[i].column < expected[j].column
			})
			if len(got) != len(expected) {
				t.Errorf("Tokenize(%s) token count mismatch, expected:%d got:%d", tc.src, len(expected), len(got))
			}
			for i, tok := range got {
				if !tokenMatches(tok, expected[i]) {
					t.Errorf("Tokenize(%s) expected:%+v got line:%d column:%d value:%s", tc.src, expected[i], tok.Position.Line, tok.Position.Column, tok.Value)
				}
			}
		})
	}
}

func tokenMatches(t *token.Token, e testToken) bool {
	return t != nil && t.Position != nil &&
		t.Value == e.value &&
		t.Position.Line == e.line &&
		t.Position.Column == e.column
}

func TestMultiLineToken_ValueLineColumnPosition(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		expect []testToken
	}{
		{
			name: "double quote",
			src: `one: "1 2 3 4 5"
two: "1 2
3 4
5"
three: "1 2 3 4
5"`,
			expect: []testToken{
				{
					line:   1,
					column: 1,
					value:  "one",
				},
				{
					line:   1,
					column: 4,
					value:  ":",
				},
				{
					line:   1,
					column: 6,
					value:  "1 2 3 4 5",
				},
				{
					line:   2,
					column: 1,
					value:  "two",
				},
				{
					line:   2,
					column: 4,
					value:  ":",
				},
				{
					line:   2,
					column: 6,
					value:  "1 2 3 4 5",
				},
				{
					line:   5,
					column: 1,
					value:  "three",
				},
				{
					line:   5,
					column: 6,
					value:  ":",
				},
				{
					line:   5,
					column: 8,
					value:  "1 2 3 4 5",
				},
			},
		},
		{
			name: "single quote in an array",
			src: `arr: ['1', 'and
two']
last: 'hello'`,
			expect: []testToken{
				{
					line:   1,
					column: 1,
					value:  "arr",
				},
				{
					line:   1,
					column: 4,
					value:  ":",
				},
				{
					line:   1,
					column: 6,
					value:  "[",
				},
				{
					line:   1,
					column: 7,
					value:  "1",
				},
				{
					line:   1,
					column: 10,
					value:  ",",
				},
				{
					line:   1,
					column: 12,
					value:  "and two",
				},
				{
					line:   2,
					column: 5,
					value:  "]",
				},
				{
					line:   3,
					column: 1,
					value:  "last",
				},
				{
					line:   3,
					column: 5,
					value:  ":",
				},
				{
					line:   3,
					column: 7,
					value:  "hello",
				},
			},
		},
		{
			name: "single quote and double quote",
			src: `foo: "test




bar"
foo2: 'bar2'`,
			expect: []testToken{
				{
					line:   1,
					column: 1,
					value:  "foo",
				},
				{
					line:   1,
					column: 4,
					value:  ":",
				},
				{
					line:   1,
					column: 6,
					value:  "test     bar",
				},
				{
					line:   7,
					column: 1,
					value:  "foo2",
				},
				{
					line:   7,
					column: 5,
					value:  ":",
				},
				{
					line:   7,
					column: 7,
					value:  "bar2",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := lexer.Tokenize(tc.src)
			sort.Slice(got, func(i, j int) bool {
				// sort by line, then column
				if got[i].Position.Line < got[j].Position.Line {
					return true
				} else if got[i].Position.Line == got[j].Position.Line {
					return got[i].Position.Column < got[j].Position.Column
				}
				return false
			})
			sort.Slice(tc.expect, func(i, j int) bool {
				if tc.expect[i].line < tc.expect[j].line {
					return true
				} else if tc.expect[i].line == tc.expect[j].line {
					return tc.expect[i].column < tc.expect[j].column
				}
				return false
			})
			if len(got) != len(tc.expect) {
				t.Errorf("Tokenize() token count mismatch, expected:%d got:%d", len(tc.expect), len(got))
			}
			for i, tok := range got {
				if !tokenMatches(tok, tc.expect[i]) {
					t.Errorf("Tokenize() expected:%+v got line:%d column:%d value:%s", tc.expect[i], tok.Position.Line, tok.Position.Column, tok.Value)
				}
			}
		})
	}
}
