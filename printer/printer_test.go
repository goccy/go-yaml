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
  12 |  jjjj`
		if actual != expect {
			t.Fatalf("unexpected output: expect:[%s]\n actual:[%s]", expect, actual)
		}
	})
	t.Run("output with color", func(t *testing.T) {
		tokens := lexer.Tokenize(yml)
		var p printer.Printer
		t.Logf("%s", p.PrintErrorToken(tokens[6], true))
	})
}
