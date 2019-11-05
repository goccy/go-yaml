package parser

import "github.com/goccy/go-yaml/token"

// context context at parsing
type context struct {
	idx    int
	size   int
	tokens token.Tokens
	mode   Mode
}

func (c *context) next() bool {
	return c.idx < c.size
}

func (c *context) previousToken() *token.Token {
	if c.idx > 0 {
		return c.tokens[c.idx-1]
	}
	return nil
}

func (c *context) currentToken() *token.Token {
	if c.idx >= c.size {
		return nil
	}
	return c.tokens[c.idx]
}

func (c *context) nextToken() *token.Token {
	if c.size > c.idx+1 {
		return c.tokens[c.idx+1]
	}
	return nil
}

func (c *context) afterNextToken() *token.Token {
	if c.size > c.idx+2 {
		return c.tokens[c.idx+2]
	}
	return nil
}

func (c *context) enabledComment() bool {
	return c.mode&ParseComments != 0
}

func (c *context) progress(num int) {
	if c.size <= c.idx+num {
		c.idx = c.size
	} else {
		c.idx += num
	}
}

func newContext(tokens token.Tokens, mode Mode) *context {
	filteredTokens := token.Tokens{}
	if mode&ParseComments != 0 {
		filteredTokens = tokens
	} else {
		for _, tk := range tokens {
			if tk.Type == token.CommentType {
				continue
			}
			filteredTokens.Add(tk)
		}
	}
	return &context{
		idx:    0,
		size:   len(filteredTokens),
		tokens: filteredTokens,
		mode:   mode,
	}
}
