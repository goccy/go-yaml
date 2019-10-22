package parser

import (
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/errors"
	"github.com/goccy/go-yaml/token"

	"golang.org/x/xerrors"
)

// Parser convert from token instances to ast
type Parser struct {
}

// Context context at parsing
type Context struct {
	idx    int
	size   int
	tokens token.Tokens
}

func (ctx *Context) next() bool {
	return ctx.idx < ctx.size
}

func (ctx *Context) previousToken() *token.Token {
	if ctx.idx > 0 {
		return ctx.tokens[ctx.idx-1]
	}
	return nil
}

func (ctx *Context) currentToken() *token.Token {
	if ctx.idx >= ctx.size {
		return nil
	}
	return ctx.tokens[ctx.idx]
}

func (ctx *Context) nextToken() *token.Token {
	if ctx.size > ctx.idx+1 {
		return ctx.tokens[ctx.idx+1]
	}
	return nil
}

func (ctx *Context) afterNextToken() *token.Token {
	if ctx.size > ctx.idx+2 {
		return ctx.tokens[ctx.idx+2]
	}
	return nil
}

func (ctx *Context) progress(num int) {
	if ctx.size <= ctx.idx+num {
		ctx.idx = ctx.size
	} else {
		ctx.idx += num
	}
}

func newContext(tokens token.Tokens) *Context {
	removedCommentTokens := token.Tokens{}
	for _, tk := range tokens {
		if tk.Type == token.CommentType {
			continue
		}
		removedCommentTokens.Add(tk)
	}
	return &Context{
		idx:    0,
		size:   len(removedCommentTokens),
		tokens: removedCommentTokens,
	}
}

func (p *Parser) parseMapping(ctx *Context) (ast.Node, error) {
	node := &ast.FlowMappingNode{
		Start:  ctx.currentToken(),
		Values: []*ast.MappingValueNode{},
	}
	ctx.progress(1) // skip MappingStart token
	for ctx.next() {
		tk := ctx.currentToken()
		if tk.Type == token.MappingEndType {
			node.End = tk
			break
		} else if tk.Type == token.CollectEntryType {
			ctx.progress(1)
			continue
		}

		value, err := p.parseToken(ctx, tk)
		if err != nil {
			return nil, xerrors.Errorf("failed to parse mapping value in mapping node: %w", err)
		}
		mvnode, ok := value.(*ast.MappingValueNode)
		if !ok {
			return nil, xerrors.New("failed to parse flow mapping value node")
		}
		node.Values = append(node.Values, mvnode)
		ctx.progress(1)
	}
	return node, nil
}

func (p *Parser) parseSequence(ctx *Context) (ast.Node, error) {
	node := &ast.FlowSequenceNode{
		Start:  ctx.currentToken(),
		Values: []ast.Node{},
	}
	ctx.progress(1) // skip SequenceStart token
	for ctx.next() {
		tk := ctx.currentToken()
		if tk.Type == token.SequenceEndType {
			node.End = tk
			break
		} else if tk.Type == token.CollectEntryType {
			ctx.progress(1)
			continue
		}

		value, err := p.parseToken(ctx, tk)
		if err != nil {
			return nil, xerrors.Errorf("failed to parse sequence value in flow sequence node: %w", err)
		}
		node.Values = append(node.Values, value)
		ctx.progress(1)
	}
	return node, nil
}

func (p *Parser) parseTag(ctx *Context) (ast.Node, error) {
	node := &ast.TagNode{Start: ctx.currentToken()}
	ctx.progress(1) // skip tag token
	value, err := p.parseToken(ctx, ctx.currentToken())
	if err != nil {
		return nil, xerrors.Errorf("failed to parse tag value: %w", err)
	}
	node.Value = value
	return node, nil
}

func (p *Parser) isMapNode(node ast.Node) bool {
	if _, ok := node.(*ast.MappingValueNode); ok {
		return true
	}
	return false
}

func (p *Parser) validateMapKey(tk *token.Token) error {
	if tk.Type != token.StringType {
		return nil
	}
	origin := strings.TrimLeft(tk.Origin, "\n")
	if strings.Index(origin, "\n") > 0 {
		return errors.NewSyntaxError("unexpected key name", tk)
	}
	return nil
}

func (p *Parser) parseMappingValue(ctx *Context) (ast.Node, error) {
	key := p.parseMapKey(ctx.currentToken())
	if key == nil {
		return nil, xerrors.New("failed to parse mapping 'key'. key is undefined")
	}
	if err := p.validateMapKey(key.GetToken()); err != nil {
		return nil, err
	}
	if _, ok := key.(ast.ScalarNode); !ok {
		return nil, xerrors.New("failed to parse mapping 'key', key is not scalar node")
	}
	ctx.progress(1)          // progress to mapping value token
	tk := ctx.currentToken() // get mapping value token
	ctx.progress(1)          // progress to value token
	value, err := p.parseToken(ctx, ctx.currentToken())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse mapping 'value' node")
	}
	keyColumn := key.GetToken().Position.Column
	valueColumn := value.GetToken().Position.Column
	if keyColumn < valueColumn && p.isMapNode(value) {
		// sub mapping
		node := &ast.MappingCollectionNode{
			Start:  tk,
			Values: []ast.Node{value},
		}
		startTk := tk
		tk := ctx.afterNextToken()
		for tk != nil && tk.Type == token.MappingValueType {
			ctx.progress(1)
			value, err := p.parseToken(ctx, ctx.currentToken())
			if err != nil {
				return nil, xerrors.Errorf("failed to parse mapping value node: %w", err)
			}
			if !p.isMapNode(value) {
				return nil, xerrors.Errorf("failed to parse mapping value node")
			}
			node.Values = append(node.Values, value)
			nextKeyToken := ctx.nextToken()
			if nextKeyToken == nil {
				break
			}
			if nextKeyToken.Position.Column != valueColumn {
				break
			}
			tk = ctx.afterNextToken()
			if tk == nil {
				break
			}
			if tk.Type != token.MappingValueType {
				break
			}
		}
		return &ast.MappingValueNode{
			Start: startTk,
			Key:   key,
			Value: node,
		}, nil
	}
	if keyColumn == valueColumn {
		if value.Type() == ast.StringType {
			ntk := ctx.nextToken()
			if ntk == nil || (ntk.Type != token.MappingValueType && ntk.Type != token.SequenceEntryType) {
				return nil, errors.NewSyntaxError("could not found expected ':' token", value.GetToken())
			}
		}
	}
	mvnode := &ast.MappingValueNode{
		Start: tk,
		Key:   key,
		Value: value,
	}
	ntk := ctx.nextToken()
	antk := ctx.afterNextToken()
	for antk != nil && antk.Type == token.MappingValueType &&
		ntk.Position.Column == key.GetToken().Position.Column {
		node := &ast.MappingCollectionNode{
			Start:  tk,
			Values: []ast.Node{mvnode},
		}
		ctx.progress(1)
		value, err := p.parseToken(ctx, ctx.currentToken())
		if err != nil {
			return nil, xerrors.Errorf("failed to parse mapping collection node: %w", err)
		}
		if c, ok := value.(*ast.MappingCollectionNode); ok {
			for _, v := range c.Values {
				node.Values = append(node.Values, v)
			}
		} else {
			node.Values = append(node.Values, value)
		}
		ntk = ctx.nextToken()
		antk = ctx.afterNextToken()
		return node, nil
	}
	return mvnode, nil
}

func (p *Parser) parseSequenceEntry(ctx *Context) (ast.Node, error) {
	tk := ctx.currentToken()
	sequenceNode := &ast.SequenceNode{
		Start:  tk,
		Values: []ast.Node{},
	}
	curIndentLevel := tk.Position.IndentLevel
	for tk.Type == token.SequenceEntryType {
		ctx.progress(1) // skip sequence token
		value, err := p.parseToken(ctx, ctx.currentToken())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse sequence")
		}
		sequenceNode.Values = append(sequenceNode.Values, value)
		tk = ctx.nextToken()
		if tk == nil {
			break
		}
		if tk.Type != token.SequenceEntryType {
			break
		}
		if tk.Position.IndentLevel != curIndentLevel {
			break
		}
		ctx.progress(1)
	}
	return sequenceNode, nil
}

func (p *Parser) parseAnchor(ctx *Context) (ast.Node, error) {
	tk := ctx.currentToken()
	anchor := &ast.AnchorNode{Start: tk}
	ntk := ctx.nextToken()
	if ntk == nil {
		return nil, xerrors.New("failed to parse anchor. anchor name is undefined")
	}
	ctx.progress(1) // skip anchor token
	name, err := p.parseToken(ctx, ctx.currentToken())
	if err != nil {
		return nil, xerrors.Errorf("failed to parser anchor name node: %w", err)
	}
	anchor.Name = name
	ntk = ctx.nextToken()
	if ntk == nil {
		return nil, xerrors.New("failed to parse anchor. anchor value is undefined")
	}
	ctx.progress(1)
	value, err := p.parseToken(ctx, ctx.currentToken())
	if err != nil {
		return nil, xerrors.Errorf("failed to parser anchor name node: %w", err)
	}
	anchor.Value = value
	return anchor, nil
}

func (p *Parser) parseAlias(ctx *Context) (ast.Node, error) {
	tk := ctx.currentToken()
	alias := &ast.AliasNode{Start: tk}
	ntk := ctx.nextToken()
	if ntk == nil {
		return nil, xerrors.New("failed to parse alias. alias name is undefined")
	}
	ctx.progress(1) // skip alias token
	name, err := p.parseToken(ctx, ctx.currentToken())
	if err != nil {
		return nil, xerrors.Errorf("failed to parser alias name node: %w", err)
	}
	alias.Value = name
	return alias, nil
}

func (p *Parser) parseMapKey(tk *token.Token) ast.Node {
	if node := p.parseStringValue(tk); node != nil {
		return node
	}
	if tk.Type == token.MergeKeyType {
		return ast.MergeKey(tk)
	}
	return nil
}

func (p *Parser) parseStringValue(tk *token.Token) ast.Node {
	switch tk.Type {
	case token.StringType,
		token.SingleQuoteType,
		token.DoubleQuoteType:
		return ast.String(tk)
	}
	return nil
}

func (p *Parser) parseScalarValue(tk *token.Token) ast.Node {
	if node := p.parseStringValue(tk); node != nil {
		return node
	}
	switch tk.Type {
	case token.NullType:
		return ast.Null(tk)
	case token.BoolType:
		return ast.Bool(tk)
	case token.IntegerType:
		return ast.Integer(tk)
	case token.FloatType:
		return ast.Float(tk)
	case token.InfinityType:
		return ast.Infinity(tk)
	case token.NanType:
		return ast.Nan(tk)
	}
	return nil
}

func (p *Parser) parseDirective(ctx *Context) (ast.Node, error) {
	node := &ast.DirectiveNode{Start: ctx.currentToken()}
	ctx.progress(1) // skip directive token
	value, err := p.parseToken(ctx, ctx.currentToken())
	if err != nil {
		return nil, xerrors.Errorf("failed to parse directive value: %w", err)
	}
	node.Value = value
	ctx.progress(1)
	if ctx.currentToken().Type != token.DocumentHeaderType {
		return nil, xerrors.New("failed to parse directive value. document not started")
	}
	return node, nil
}

func (p *Parser) parseLiteral(ctx *Context) (ast.Node, error) {
	node := &ast.LiteralNode{Start: ctx.currentToken()}
	ctx.progress(1) // skip literal/folded token
	value, err := p.parseToken(ctx, ctx.currentToken())
	if err != nil {
		return nil, xerrors.Errorf("failed to parse literal/folded value: %w", err)
	}
	snode, ok := value.(*ast.StringNode)
	if !ok {
		return nil, xerrors.New("invalid literal type. required string node")
	}
	node.Value = snode
	return node, nil
}

func (p *Parser) parseToken(ctx *Context, tk *token.Token) (ast.Node, error) {
	if tk.NextType() == token.MappingValueType {
		return p.parseMappingValue(ctx)
	}
	if node := p.parseScalarValue(tk); node != nil {
		return node, nil
	}
	switch tk.Type {
	case token.MappingStartType:
		return p.parseMapping(ctx)
	case token.SequenceStartType:
		return p.parseSequence(ctx)
	case token.SequenceEntryType:
		return p.parseSequenceEntry(ctx)
	case token.AnchorType:
		return p.parseAnchor(ctx)
	case token.AliasType:
		return p.parseAlias(ctx)
	case token.DirectiveType:
		return p.parseDirective(ctx)
	case token.TagType:
		return p.parseTag(ctx)
	case token.LiteralType, token.FoldedType:
		return p.parseLiteral(ctx)
	}
	return nil, nil
}

// Parse parse from token instances, and returns ast.Document
func (p *Parser) Parse(tokens token.Tokens) (*ast.Document, error) {
	ctx := newContext(tokens)
	doc := &ast.Document{Nodes: []ast.Node{}}
	for ctx.next() {
		node, err := p.parseToken(ctx, ctx.currentToken())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse")
		}
		ctx.progress(1)
		if node == nil {
			continue
		}
		if len(doc.Nodes) == 0 {
			doc.Nodes = append(doc.Nodes, node)
			continue
		}
		if _, ok := node.(*ast.MappingValueNode); !ok {
			doc.Nodes = append(doc.Nodes, node)
			continue
		}
		lastNode := doc.Nodes[len(doc.Nodes)-1]
		switch n := lastNode.(type) {
		case *ast.MappingValueNode:
			doc.Nodes[len(doc.Nodes)-1] = &ast.MappingCollectionNode{
				Start:  n.GetToken(),
				Values: []ast.Node{lastNode, node},
			}
		case *ast.MappingCollectionNode:
			n.Values = append(n.Values, node)
		}
	}
	return doc, nil
}
