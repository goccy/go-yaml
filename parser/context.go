package parser

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml/token"
)

// context context at parsing
type context struct {
	tokenRef *tokenRef
	path     string
	isFlow   bool
}

type tokenRef struct {
	tokens []*Token
	size   int
	idx    int
}

var pathSpecialChars = []string{
	"$", "*", ".", "[", "]",
}

func containsPathSpecialChar(path string) bool {
	for _, char := range pathSpecialChars {
		if strings.Contains(path, char) {
			return true
		}
	}
	return false
}

func normalizePath(path string) string {
	if containsPathSpecialChar(path) {
		return fmt.Sprintf("'%s'", path)
	}
	return path
}

func (c *context) currentToken() *Token {
	if c.tokenRef.idx >= c.tokenRef.size {
		return nil
	}
	return c.tokenRef.tokens[c.tokenRef.idx]
}

func (c *context) isComment() bool {
	return c.currentToken().Type() == token.CommentType
}

func (c *context) nextToken() *Token {
	if c.tokenRef.idx+1 >= c.tokenRef.size {
		return nil
	}
	return c.tokenRef.tokens[c.tokenRef.idx+1]
}

func (c *context) nextNotCommentToken() *Token {
	for i := c.tokenRef.idx + 1; i < c.tokenRef.size; i++ {
		tk := c.tokenRef.tokens[i]
		if tk.Type() == token.CommentType {
			continue
		}
		return tk
	}
	return nil
}

func (c *context) withGroup(g *TokenGroup) *context {
	ctx := *c
	ctx.tokenRef = &tokenRef{
		tokens: g.Tokens,
		size:   len(g.Tokens),
	}
	return &ctx
}

func (c *context) withChild(path string) *context {
	ctx := *c
	ctx.path = c.path + "." + normalizePath(path)
	return &ctx
}

func (c *context) withIndex(idx uint) *context {
	ctx := *c
	ctx.path = c.path + "[" + fmt.Sprint(idx) + "]"
	return &ctx
}

func (c *context) withFlow(isFlow bool) *context {
	ctx := *c
	ctx.isFlow = isFlow
	return &ctx
}

func newContext() *context {
	return &context{
		path: "$",
	}
}

func (c *context) goNext() {
	ref := c.tokenRef
	if ref.size <= ref.idx+1 {
		ref.idx = ref.size
	} else {
		ref.idx++
	}
}

func (c *context) next() bool {
	return c.tokenRef.idx < c.tokenRef.size
}

func (c *context) insertToken(tk *Token) {
	idx := c.tokenRef.idx
	if c.tokenRef.size < idx {
		return
	}
	if c.tokenRef.size == idx {
		curToken := c.tokenRef.tokens[c.tokenRef.size-1]
		tk.RawToken().Next = curToken.RawToken()
		curToken.RawToken().Prev = tk.RawToken()

		c.tokenRef.tokens = append(c.tokenRef.tokens, tk)
		c.tokenRef.size = len(c.tokenRef.tokens)
		return
	}

	curToken := c.tokenRef.tokens[idx]
	tk.RawToken().Next = curToken.RawToken()
	curToken.RawToken().Prev = tk.RawToken()

	c.tokenRef.tokens = append(c.tokenRef.tokens[:idx+1], c.tokenRef.tokens[idx:]...)
	c.tokenRef.tokens[idx] = tk
	c.tokenRef.size = len(c.tokenRef.tokens)
}
