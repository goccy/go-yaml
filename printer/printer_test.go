package printer_test

import (
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
	t.Run("print stargin from tokens[4]", func(t *testing.T) {
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
	t.Run("output with color", func(t *testing.T) {
		t.Run("token6", func(t *testing.T) {
			tokens := lexer.Tokenize(yml)
			var p printer.Printer
			t.Logf("%s", p.PrintErrorToken(tokens[6], true))
		})
		t.Run("token9", func(t *testing.T) {
			tokens := lexer.Tokenize(yml)
			var p printer.Printer
			t.Logf("%s", p.PrintErrorToken(tokens[9], true))
		})
		t.Run("token12", func(t *testing.T) {
			tokens := lexer.Tokenize(yml)
			var p printer.Printer
			t.Logf("%s", p.PrintErrorToken(tokens[12], true))
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
