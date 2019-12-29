package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/printer"
	"github.com/mattn/go-colorable"
)

const escape = "\x1b"

func format(attr color.Attribute) string {
	return fmt.Sprintf("%s[%dm", escape, attr)
}

func _main(args []string) error {
	if len(args) < 2 {
		return errors.New("ycat: usage: ycat file.yml")
	}
	filename := args[1]
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	tokens := lexer.Tokenize(string(bytes))
	var p printer.Printer
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
	writer := colorable.NewColorableStdout()
	writer.Write([]byte(p.PrintTokens(tokens) + "\n"))
	return nil
}

func main() {
	if err := _main(os.Args); err != nil {
		fmt.Printf("%v\n", yaml.FormatError(err, true, true))
	}
}
