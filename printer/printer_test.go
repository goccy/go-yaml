package printer_test

import (
	"fmt"
	"testing"

	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/printer"
)

func Test_Printer(t *testing.T) {
	yml := `---
text: aaaa
text2: aaaa
 bbbb
 cccc
 dddd
 eeee
text3: ffff
 gggg
 hhhh
 iiii
 jjjj
bool: true
number: 10
anchor: &x 1
alias: *x
`
	t.Run("print starting from tokens[3]", func(t *testing.T) {
		tokens := lexer.Tokenize(yml)
		var p printer.Printer
		actual := "\n" + p.PrintErrorToken(tokens[3], false)
		expect := `
   1 | ---
>  2 | text: aaaa
             ^
   3 | text2: aaaa
   4 |  bbbb
   5 |  cccc
   6 |  dddd
   7 |  eeee
   8 | `
		if actual != expect {
			t.Fatalf("unexpected output: expect:[%s]\n actual:[%s]", expect, actual)
		}
	})
	t.Run("print starting from tokens[4]", func(t *testing.T) {
		tokens := lexer.Tokenize(yml)
		var p printer.Printer
		actual := "\n" + p.PrintErrorToken(tokens[4], false)
		expect := `
   1 | ---
   2 | text: aaaa
>  3 | text2: aaaa
   4 |  bbbb
   5 |  cccc
   6 |  dddd
   7 |  eeee
       ^
`
		if actual != expect {
			t.Fatalf("unexpected output: expect:[%s]\n actual:[%s]", expect, actual)
		}
	})
	t.Run("print starting from tokens[6]", func(t *testing.T) {
		tokens := lexer.Tokenize(yml)
		var p printer.Printer
		actual := "\n" + p.PrintErrorToken(tokens[6], false)
		expect := `
   1 | ---
   2 | text: aaaa
>  3 | text2: aaaa
   4 |  bbbb
   5 |  cccc
   6 |  dddd
   7 |  eeee
              ^
   8 | text3: ffff
   9 |  gggg
  10 |  hhhh
  11 |  iiii
  12 |  jjjj
  13 | `
		if actual != expect {
			t.Fatalf("unexpected output: expect:[%s]\n actual:[%s]", expect, actual)
		}
	})
	t.Run("print error token with document header", func(t *testing.T) {
		tokens := lexer.Tokenize(`---
a:
 b:
  c:
   d: e
   f: g
   h: i

---
`)
		expect := `
   3 |  b:
   4 |   c:
   5 |    d: e
>  6 |    f: g
             ^
   7 |    h: i
   8 | 
   9 | ---`
		var p printer.Printer
		actual := "\n" + p.PrintErrorToken(tokens[12], false)
		if actual != expect {
			t.Fatalf("unexpected output: expect:[%s]\n actual:[%s]", expect, actual)
		}
	})
	t.Run("output with color", func(t *testing.T) {
		t.Run("token6", func(t *testing.T) {
			tokens := lexer.Tokenize(yml)
			var p printer.Printer
			t.Logf("\n%s", p.PrintErrorToken(tokens[6], true))
		})
		t.Run("token9", func(t *testing.T) {
			tokens := lexer.Tokenize(yml)
			var p printer.Printer
			t.Logf("\n%s", p.PrintErrorToken(tokens[9], true))
		})
		t.Run("token12", func(t *testing.T) {
			tokens := lexer.Tokenize(yml)
			var p printer.Printer
			t.Logf("\n%s", p.PrintErrorToken(tokens[12], true))
		})
	})
	t.Run("print error message", func(t *testing.T) {
		var p printer.Printer
		src := "message"
		msg := p.PrintErrorMessage(src, false)
		if msg != src {
			t.Fatal("unexpected result")
		}
		p.PrintErrorMessage(src, true)
	})
}

func TestPrinter_Anchor(t *testing.T) {
	expected := `
anchor: &x 1
alias: *x`
	tokens := lexer.Tokenize(expected)
	var p printer.Printer
	got := p.PrintTokens(tokens)
	if expected != got {
		t.Fatalf("unexpected output: expect:[%s]\n actual:[%s]", expected, got)
	}
}

func Test_Printer_Multiline(t *testing.T) {
	yml := `
text1: 'aaaa
 bbbb
 cccc'
text2: "ffff
 gggg
 hhhh"
text3: hello
`
	tc := []struct {
		token int
		want  string
	}{
		{
			token: 2,
			want: `
>  2 | text1: 'aaaa
   3 |  bbbb
   4 |  cccc'
              ^
   5 | text2: "ffff
   6 |  gggg
   7 |  hhhh"`,
		},
		{token: 3,
			want: `
   2 | text1: 'aaaa
   3 |  bbbb
   4 |  cccc'
>  5 | text2: "ffff
   6 |  gggg
   7 |  hhhh"
       ^
   8 | text3: hello`,
		},
		{token: 5,
			want: `
   2 | text1: 'aaaa
   3 |  bbbb
   4 |  cccc'
>  5 | text2: "ffff
   6 |  gggg
   7 |  hhhh"
              ^
   8 | text3: hello`,
		},
		{token: 6,
			want: `
   5 | text2: "ffff
   6 |  gggg
   7 |  hhhh"
>  8 | text3: hello
       ^
`,
		},
	}
	for _, tt := range tc {
		name := fmt.Sprintf("print starting from tokens[%d]", tt.token)
		t.Run(name, func(t *testing.T) {
			tokens := lexer.Tokenize(yml)
			var p printer.Printer
			got := "\n" + p.PrintErrorToken(tokens[tt.token], false)
			want := tt.want
			if got != want {
				t.Fatalf("PrintErrorToken() got: %s\n want:%s\n", want, got)
			}
		})
	}
}
