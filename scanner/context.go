package scanner

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/goccy/go-yaml/token"
)

// Context context at scanning
type Context struct {
	idx                      int
	size                     int
	notSpaceCharPos          int
	notSpaceOrgCharPos       int
	src                      []rune
	buf                      []rune
	obuf                     []rune
	tokens                   token.Tokens
	isRawFolded              bool
	isLiteral                bool
	isFolded                 bool
	docOpt                   string
	docFirstLineIndentColumn int
	docPrevLineIndentColumn  int
	docLineIndentColumn      int
	docFoldedNewLine         bool
}

var (
	ctxPool = sync.Pool{
		New: func() interface{} {
			return createContext()
		},
	}
)

func createContext() *Context {
	return &Context{
		idx:    0,
		tokens: token.Tokens{},
	}
}

func newContext(src []rune) *Context {
	ctx, _ := ctxPool.Get().(*Context)
	ctx.reset(src)
	return ctx
}

func (c *Context) release() {
	ctxPool.Put(c)
}

func (c *Context) clear() {
	c.resetBuffer()
	c.isRawFolded = false
	c.isLiteral = false
	c.isFolded = false
	c.docOpt = ""
	c.docFirstLineIndentColumn = 0
	c.docLineIndentColumn = 0
	c.docPrevLineIndentColumn = 0
	c.docFoldedNewLine = false
}

func (c *Context) reset(src []rune) {
	c.idx = 0
	c.size = len(src)
	c.src = src
	c.tokens = c.tokens[:0]
	c.resetBuffer()
	c.isRawFolded = false
	c.isLiteral = false
	c.isFolded = false
	c.docOpt = ""
}

func (c *Context) resetBuffer() {
	c.buf = c.buf[:0]
	c.obuf = c.obuf[:0]
	c.notSpaceCharPos = 0
	c.notSpaceOrgCharPos = 0
}

func (c *Context) breakDocument() {
	c.isLiteral = false
	c.isRawFolded = false
	c.isFolded = false
	c.docOpt = ""
	c.docFirstLineIndentColumn = 0
	c.docLineIndentColumn = 0
	c.docPrevLineIndentColumn = 0
	c.docFoldedNewLine = false
}

func (c *Context) updateDocumentIndentColumn() {
	indent := c.docFirstLineIndentColumnByDocOpt()
	if indent > 0 {
		c.docFirstLineIndentColumn = indent + 1
	}
}

func (c *Context) docFirstLineIndentColumnByDocOpt() int {
	trimmed := strings.TrimPrefix(c.docOpt, "-")
	trimmed = strings.TrimPrefix(trimmed, "+")
	i, _ := strconv.ParseInt(trimmed, 10, 64)
	return int(i)
}

func (c *Context) updateDocumentLineIndentColumn(column int) {
	if c.docFirstLineIndentColumn == 0 {
		c.docFirstLineIndentColumn = column
	}
	if c.docLineIndentColumn == 0 {
		c.docLineIndentColumn = column
	}
}

func (c *Context) validateDocumentLineIndentColumn() error {
	if c.docFirstLineIndentColumnByDocOpt() == 0 {
		return nil
	}
	if c.docFirstLineIndentColumn > c.docLineIndentColumn {
		return fmt.Errorf("invalid number of indent is specified in the document header")
	}
	return nil
}

func (c *Context) updateDocumentNewLineState() {
	c.docPrevLineIndentColumn = c.docLineIndentColumn
	c.docFoldedNewLine = true
	c.docLineIndentColumn = 0
}

func (c *Context) addDocumentIndent(column int) {
	if c.docFirstLineIndentColumn == 0 {
		return
	}

	// If the first line of the document has already been evaluated, the number is treated as the threshold, since the `docFirstLineIndentColumn` is a positive number.
	if c.docFirstLineIndentColumn <= column {
		// In the folded state, new-line-char is normally treated as space,
		// but if the number of indents is different from the number of indents in the first line,
		// new-line-char is used as is instead of space.
		// Therefore, it is necessary to replace the space already added to buf.
		// `c.docFoldedNewLine` is a variable that is set to true for every newline.
		if c.isFolded && c.docFoldedNewLine {
			c.buf[len(c.buf)-1] = '\n'
			c.docFoldedNewLine = false
		}
		// Since addBuf ignore space character, add to the buffer directly.
		c.buf = append(c.buf, ' ')
	}
}

func (c *Context) addDocumentNewLineInFolded(column int) {
	if !c.isFolded {
		return
	}
	if !c.docFoldedNewLine {
		return
	}
	if c.docFirstLineIndentColumn == c.docLineIndentColumn &&
		c.docLineIndentColumn == c.docPrevLineIndentColumn {
		// use space as a new line delimiter.
		return
	}
	c.buf[len(c.buf)-1] = '\n'
	c.docFoldedNewLine = false
}

func (c *Context) addToken(tk *token.Token) {
	if tk == nil {
		return
	}
	c.tokens = append(c.tokens, tk)
}

func (c *Context) addBuf(r rune) {
	if len(c.buf) == 0 && r == ' ' {
		return
	}
	c.buf = append(c.buf, r)
	if r != ' ' && r != '\t' {
		c.notSpaceCharPos = len(c.buf)
	}
}

func (c *Context) addOriginBuf(r rune) {
	c.obuf = append(c.obuf, r)
	if r != ' ' && r != '\t' {
		c.notSpaceOrgCharPos = len(c.obuf)
	}
}

func (c *Context) removeRightSpaceFromBuf() int {
	trimmedBuf := c.obuf[:c.notSpaceOrgCharPos]
	buflen := len(trimmedBuf)
	diff := len(c.obuf) - buflen
	if diff > 0 {
		c.obuf = c.obuf[:buflen]
		c.buf = c.bufferedSrc()
	}
	return diff
}

func (c *Context) isDocument() bool {
	return c.isLiteral || c.isFolded || c.isRawFolded
}

func (c *Context) isEOS() bool {
	return len(c.src)-1 <= c.idx
}

func (c *Context) isNextEOS() bool {
	return len(c.src) <= c.idx+1
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
	if c.size > c.idx {
		return c.src[c.idx]
	}
	return rune(0)
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

func (c *Context) existsBuffer() bool {
	return len(c.bufferedSrc()) != 0
}

func (c *Context) bufferedSrc() []rune {
	src := c.buf[:c.notSpaceCharPos]
	if c.isDocument() && strings.HasPrefix(c.docOpt, "-") {
		// remove end '\n' character and trailing empty lines
		// https://yaml.org/spec/1.2.2/#8112-block-chomping-indicator
		for {
			if len(src) > 0 && src[len(src)-1] == '\n' {
				src = src[:len(src)-1]
				continue
			}
			break
		}
		for {
			if len(src) > 0 && src[len(src)-1] == ' ' {
				src = src[:len(src)-1]
				continue
			}
			break
		}
	}
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
	var tk *token.Token
	if c.isDocument() {
		tk = token.String(string(source), string(c.obuf), pos)
	} else {
		tk = token.New(string(source), string(c.obuf), pos)
	}
	c.resetBuffer()
	return tk
}

func (c *Context) lastToken() *token.Token {
	if len(c.tokens) != 0 {
		return c.tokens[len(c.tokens)-1]
	}
	return nil
}
