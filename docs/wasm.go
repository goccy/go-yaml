package main

import (
	"fmt"
	"syscall/js"

	"github.com/goccy/go-yaml/lexer"
)

func tokenize(this js.Value, args []js.Value) interface{} {
	source := args[0]
	fmt.Println("source = ", source)
	tokens := lexer.Tokenize(source.String())
	jsTokens := []interface{}{}
	for _, token := range tokens {
		jsTokens = append(jsTokens, map[string]interface{}{
			"type":   token.Type.String(),
			"value":  token.Value,
			"origin": token.Origin,
			"position": map[string]interface{}{
				"line":        token.Position.Line,
				"column":      token.Position.Column,
				"offset":      token.Position.Offset,
				"indentNum":   token.Position.IndentNum,
				"indentLevel": token.Position.IndentLevel,
			},
		})
	}
	return jsTokens
}

func main() {
	fmt.Println("CALLED main")
	js.Global().Set("tokenize", js.FuncOf(tokenize))
	done := make(chan struct{}, 0)
	<-done
}
