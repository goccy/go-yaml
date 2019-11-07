package scanner

import (
	"strings"

	"github.com/goccy/go-yaml/token"
)

// Context context at scanning
type Context struct {
	idx         int
	size        int
	src         []rune
	buf         []rune
	obuf        []rune
	tokens      token.Tokens
	isRawFolded bool
	isLiteral   bool
	isFolded    bool
	literalOpt  string
}

func newContext(src []rune) *Context {
	return &Context{
		idx:    0,
		size:   len(src),
		src:    src,
		tokens: token.Tokens{},
		buf:    make([]rune, 0, len(src)),
		obuf:   make([]rune, 0, len(src)),
	}
}

func (c *Context) resetBuffer() {
	c.buf = c.buf[:0]
	c.obuf = c.obuf[:0]
}

func (c *Context) isSaveIndentMode() bool {
	return c.isLiteral || c.isFolded || c.isRawFolded
}

func (c *Context) breakLiteral() {
	c.isLiteral = false
	c.isRawFolded = false
	c.isFolded = false
	c.literalOpt = ""
}

func (c *Context) addToken(tk *token.Token) {
	if tk == nil {
		return
	}
	c.tokens = append(c.tokens, tk)
}

func (c *Context) addBuf(r rune) {
	c.buf = append(c.buf, r)
}

func (c *Context) addOriginBuf(r rune) {
	c.obuf = append(c.obuf, r)
}

func (c *Context) isEOS() bool {
	return len(c.src)-1 <= c.idx
}

func (c *Context) isNextEOS() bool {
	return len(c.src)-1 <= c.idx+1
}

func (c *Context) next() bool {
	return c.idx < c.size
}

func (c *Context) source(s, e int) string {
	return string(c.src[s:e])
}

func (c *Context) previousChar() rune {
	if c.idx > 0 {
		return c.src[c.idx-1]
	}
	return rune(0)
}

func (c *Context) currentChar() rune {
	return c.src[c.idx]
}

func (c *Context) nextChar() rune {
	if c.size > c.idx+1 {
		return c.src[c.idx+1]
	}
	return rune(0)
}

func (c *Context) repeatNum(r rune) int {
	cnt := 0
	for i := c.idx; i < c.size; i++ {
		if c.src[i] == r {
			cnt++
		} else {
			break
		}
	}
	return cnt
}

func (c *Context) progress(num int) {
	c.idx += num
}

func (c *Context) nextPos() int {
	return c.idx + 1
}

func (c *Context) bufferedSrc() string {
	src := strings.TrimLeft(string(c.buf), " ")
	src = strings.TrimRight(src, " ")
	return src
}

func (c *Context) bufferedToken(pos *token.Position) *token.Token {
	if c.idx == 0 {
		return nil
	}
	source := c.bufferedSrc()
	if len(source) == 0 {
		return nil
	}
	tk := token.New(source, string(c.obuf), pos)
	c.buf = c.buf[:0]
	c.obuf = c.obuf[:0]
	return tk
}
