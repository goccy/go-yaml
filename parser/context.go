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

func (c *context) nextNotCommentToken() *token.Token {
	for i := c.idx; i+1 < c.size; i++ {
		tk := c.tokens[i+1]
		if tk.Type == token.CommentType {
			continue
		}
		return tk
	}
	return nil

}

func (c *context) nextToken() *token.Token {
	if c.idx+1 >= c.size {
		return nil
	}
	return c.tokens[c.idx+1]
}

func (c *context) afterNextNotCommentToken() *token.Token {
	for i := c.idx; i+2 < c.size; i++ {
		tk := c.tokens[i+2]
		if tk.Type == token.CommentType {
			continue
		}
		return tk
	}
	return nil
}

func (c *context) afterNextToken() *token.Token {
	if c.idx+2 >= c.size {
		return nil
	}
	return c.tokens[c.idx+2]
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
