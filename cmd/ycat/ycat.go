package main

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/printer"
)

const escape = "\x1b"

func format(attr color.Attribute) string {
	return fmt.Sprintf("%s[%dm", escape, attr)
}

func _main(args []string) error {
	if len(args) < 2 {
		return errors.New("ycat: usage: ycat file.yml [file.yml...]")
	}

	var p printer.Printer
	if term.IsTerminal(int(os.Stdout.Fd())) {
		p.LineNumber = true
		p.LineNumberFormat = func(num int) string {
			fn := color.New(color.Bold, color.FgHiWhite).SprintFunc()
			return fn(fmt.Sprintf("%2d | ", num))
		}
		p.Bool = func() *printer.Property {
			return &printer.Property{
				Prefix: format(color.FgHiMagenta),
				Suffix: format(color.Reset),
			}
		}
		p.Number = func() *printer.Property {
			return &printer.Property{
				Prefix: format(color.FgHiMagenta),
				Suffix: format(color.Reset),
			}
		}
		p.MapKey = func() *printer.Property {
			return &printer.Property{
				Prefix: format(color.FgHiCyan),
				Suffix: format(color.Reset),
			}
		}
		p.Anchor = func() *printer.Property {
			return &printer.Property{
				Prefix: format(color.FgHiYellow),
				Suffix: format(color.Reset),
			}
		}
		p.Alias = func() *printer.Property {
			return &printer.Property{
				Prefix: format(color.FgHiYellow),
				Suffix: format(color.Reset),
			}
		}
		p.String = func() *printer.Property {
			return &printer.Property{
				Prefix: format(color.FgHiGreen),
				Suffix: format(color.Reset),
			}
		}
		p.Comment = func() *printer.Property {
			return &printer.Property{
				Prefix: format(color.FgHiBlack),
				Suffix: format(color.Reset),
			}
		}
	}

	writer := colorable.NewColorableStdout()

	for _, filename := range args[1:] {
		if len(args) > 2 {
			_, _ = writer.Write([]byte("---\n"))
		}
		bytes, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		tokens := lexer.Tokenize(string(bytes))
		_, _ = writer.Write([]byte(p.PrintTokens(tokens) + "\n"))
	}
	return nil
}

func main() {
	if err := _main(os.Args); err != nil {
		fmt.Printf("%v\n", yaml.FormatError(err, true, true))
	}
}
