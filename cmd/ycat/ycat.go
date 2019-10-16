package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/printer"
	"github.com/mattn/go-colorable"
)

func _main(args []string) error {
	if len(args) < 2 {
		return errors.New("ycat: usage: ycat file.yml")
	}
	filename := args[1]
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	var lexer lexer.Lexer
	tokens := lexer.Tokenize(string(bytes))
	var p printer.Printer
	p.LineNumber = true
	p.Bool = func(text string) string {
		fn := color.New(color.FgHiMagenta).SprintFunc()
		return fn(text)
	}
	p.MapKey = func(text string) string {
		fn := color.New(color.FgHiCyan).SprintFunc()
		return fn(text)
	}
	p.Anchor = func(text string) string {
		fn := color.New(color.FgHiYellow).SprintFunc()
		return fn(text)
	}
	p.Alias = func(text string) string {
		fn := color.New(color.FgHiYellow).SprintFunc()
		return fn(text)
	}
	p.String = func(text string) string {
		fn := color.New(color.FgHiGreen).SprintFunc()
		return fn(text)
	}
	writer := colorable.NewColorableStdout()
	writer.Write([]byte(p.PrintTokens(tokens) + "\n"))
	return nil
}

func main() {
	if err := _main(os.Args); err != nil {
		fmt.Printf("%v\n", err)
	}
}
