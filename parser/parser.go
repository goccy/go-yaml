package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/internal/errors"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/token"
)

type Mode uint

const (
	ParseComments Mode = 1 << iota // parse comments and add them to AST
)

// ParseBytes parse from byte slice, and returns ast.File
func ParseBytes(bytes []byte, mode Mode, opts ...Option) (*ast.File, error) {
	tokens := lexer.Tokenize(string(bytes))
	f, err := Parse(tokens, mode, opts...)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Parse parse from token instances, and returns ast.File
func Parse(tokens token.Tokens, mode Mode, opts ...Option) (*ast.File, error) {
	if tk := tokens.InvalidToken(); tk != nil {
		return nil, errors.ErrSyntax(tk.Error, tk)
	}
	p, err := newParser(tokens, mode, opts)
	if err != nil {
		return nil, err
	}
	f, err := p.parse(newContext())
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Parse parse from filename, and returns ast.File
func ParseFile(filename string, mode Mode, opts ...Option) (*ast.File, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	f, err := ParseBytes(file, mode, opts...)
	if err != nil {
		return nil, err
	}
	f.Name = filename
	return f, nil
}

type parser struct {
	tokens               []*Token
	pathMap              map[string]ast.Node
	allowDuplicateMapKey bool
}

func newParser(tokens token.Tokens, mode Mode, opts []Option) (*parser, error) {
	filteredTokens := []*token.Token{}
	if mode&ParseComments != 0 {
		filteredTokens = tokens
	} else {
		for _, tk := range tokens {
			if tk.Type == token.CommentType {
				continue
			}
			// keep prev/next reference between tokens containing comments
			// https://github.com/goccy/go-yaml/issues/254
			filteredTokens = append(filteredTokens, tk)
		}
	}
	tks, err := createGroupedTokens(token.Tokens(filteredTokens))
	if err != nil {
		return nil, err
	}
	p := &parser{
		tokens:  tks,
		pathMap: make(map[string]ast.Node),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p, nil
}

func (p *parser) parse(ctx *context) (*ast.File, error) {
	file := &ast.File{Docs: []*ast.DocumentNode{}}
	for _, token := range p.tokens {
		doc, err := p.parseDocument(ctx, token.Group)
		if err != nil {
			return nil, err
		}
		file.Docs = append(file.Docs, doc)
	}
	return file, nil
}

func (p *parser) parseDocument(ctx *context, docGroup *TokenGroup) (*ast.DocumentNode, error) {
	if len(docGroup.Tokens) == 0 {
		return ast.Document(docGroup.RawToken(), nil), nil
	}

	p.pathMap = make(map[string]ast.Node)

	var (
		tokens = docGroup.Tokens
		start  *token.Token
		end    *token.Token
	)
	if docGroup.First().Type() == token.DocumentHeaderType {
		start = docGroup.First().RawToken()
		tokens = tokens[1:]
	}
	if docGroup.Last().Type() == token.DocumentEndType {
		end = docGroup.Last().RawToken()
		tokens = tokens[:len(tokens)-1]
	}

	if len(tokens) == 0 {
		return ast.Document(docGroup.RawToken(), nil), nil
	}

	body, err := p.parseDocumentBody(ctx.withGroup(&TokenGroup{
		Type:   TokenGroupDocumentBody,
		Tokens: tokens,
	}))
	if err != nil {
		return nil, err
	}
	node := ast.Document(start, body)
	node.End = end
	return node, nil
}

func (p *parser) parseDocumentBody(ctx *context) (ast.Node, error) {
	node, err := p.parseToken(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	if ctx.next() {
		return nil, errors.ErrSyntax("value is not allowed in this context", ctx.currentToken().RawToken())
	}
	return node, nil
}

func (p *parser) parseToken(ctx *context, tk *Token) (ast.Node, error) {
	switch tk.GroupType() {
	case TokenGroupMapKey, TokenGroupMapKeyValue:
		return p.parseMap(ctx)
	case TokenGroupDirective:
		node, err := p.parseDirective(ctx.withGroup(tk.Group), tk.Group)
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return node, nil
	case TokenGroupAnchor:
		node, err := p.parseAnchor(ctx.withGroup(tk.Group), tk.Group)
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return node, nil
	case TokenGroupAnchorName:
		anchor, err := p.parseAnchorName(ctx.withGroup(tk.Group))
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		if ctx.isTokenNotFound() {
			return nil, errors.ErrSyntax("could not find anchor value", tk.RawToken())
		}
		value, err := p.parseToken(ctx, ctx.currentToken())
		if err != nil {
			return nil, err
		}
		anchor.Value = value
		return anchor, nil
	case TokenGroupAlias:
		node, err := p.parseAlias(ctx.withGroup(tk.Group))
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return node, nil
	case TokenGroupLiteral, TokenGroupFolded:
		node, err := p.parseLiteral(ctx.withGroup(tk.Group))
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return node, nil
	case TokenGroupScalarTag:
		node, err := p.parseTag(ctx.withGroup(tk.Group))
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return node, nil
	}
	switch tk.Type() {
	case token.CommentType:
		return p.parseComment(ctx)
	case token.TagType:
		return p.parseTag(ctx)
	case token.MappingStartType:
		return p.parseFlowMap(ctx)
	case token.SequenceStartType:
		return p.parseFlowSequence(ctx)
	case token.SequenceEntryType:
		return p.parseSequence(ctx)
	case token.SequenceEndType:
		// SequenceEndType is always validated in parseFlowSequence.
		// Therefore, if this is found in other cases, it is treated as a syntax error.
		return nil, errors.ErrSyntax("could not find '[' character corresponding to ']'", tk.RawToken())
	case token.MappingEndType:
		// MappingEndType is always validated in parseFlowMap.
		// Therefore, if this is found in other cases, it is treated as a syntax error.
		return nil, errors.ErrSyntax("could not find '{' character corresponding to '}'", tk.RawToken())
	case token.MappingValueType:
		return nil, errors.ErrSyntax("found an invalid key for this map", tk.RawToken())
	}
	node, err := p.parseScalarValue(ctx, tk)
	if err != nil {
		return nil, err
	}
	ctx.goNext()
	return node, nil
}

func (p *parser) parseScalarValue(ctx *context, tk *Token) (ast.ScalarNode, error) {
	if tk.Group != nil {
		switch tk.GroupType() {
		case TokenGroupAnchor:
			return p.parseAnchor(ctx.withGroup(tk.Group), tk.Group)
		case TokenGroupAnchorName:
			anchor, err := p.parseAnchorName(ctx.withGroup(tk.Group))
			if err != nil {
				return nil, err
			}
			ctx.goNext()
			if ctx.isTokenNotFound() {
				return nil, errors.ErrSyntax("could not find anchor value", tk.RawToken())
			}
			value, err := p.parseToken(ctx, ctx.currentToken())
			if err != nil {
				return nil, err
			}
			anchor.Value = value
			return anchor, nil
		case TokenGroupAlias:
			return p.parseAlias(ctx.withGroup(tk.Group))
		case TokenGroupLiteral, TokenGroupFolded:
			return p.parseLiteral(ctx.withGroup(tk.Group))
		case TokenGroupScalarTag:
			return p.parseTag(ctx.withGroup(tk.Group))
		default:
			return nil, errors.ErrSyntax("unexpected scalar value", tk.RawToken())
		}
	}
	switch tk.Type() {
	case token.MergeKeyType:
		return newMergeKeyNode(ctx, tk)
	case token.NullType:
		return newNullNode(ctx, tk)
	case token.BoolType:
		return newBoolNode(ctx, tk)
	case token.IntegerType, token.BinaryIntegerType, token.OctetIntegerType, token.HexIntegerType:
		return newIntegerNode(ctx, tk)
	case token.FloatType:
		return newFloatNode(ctx, tk)
	case token.InfinityType:
		return newInfinityNode(ctx, tk)
	case token.NanType:
		return newNanNode(ctx, tk)
	case token.StringType, token.SingleQuoteType, token.DoubleQuoteType:
		return newStringNode(ctx, tk)
	case token.TagType:
		// this case applies when it is a scalar tag and its value does not exist.
		// Examples of cases where the value does not exist include cases like `key: !!str,` or `!!str : value`.
		return p.parseScalarTag(ctx)
	}
	return nil, errors.ErrSyntax("unexpected scalar value type", tk.RawToken())
}

func (p *parser) parseFlowMap(ctx *context) (*ast.MappingNode, error) {
	node, err := newMappingNode(ctx, ctx.currentToken(), true)
	if err != nil {
		return nil, err
	}
	ctx.goNext() // skip MappingStart token

	isFirst := true
	for ctx.next() {
		tk := ctx.currentToken()
		if tk.Type() == token.MappingEndType {
			node.End = tk.RawToken()
			break
		}

		if tk.Type() == token.CollectEntryType {
			ctx.goNext()
		} else if !isFirst {
			return nil, errors.ErrSyntax("',' or '}' must be specified", tk.RawToken())
		}

		if tk := ctx.currentToken(); tk.Type() == token.MappingEndType {
			// this case is here: "{ elem, }".
			// In this case, ignore the last element and break mapping parsing.
			node.End = tk.RawToken()
			break
		}

		mapKeyTk := ctx.currentToken()
		switch mapKeyTk.GroupType() {
		case TokenGroupMapKeyValue:
			value, err := p.parseMapKeyValue(ctx.withGroup(mapKeyTk.Group), mapKeyTk.Group)
			if err != nil {
				return nil, err
			}
			node.Values = append(node.Values, value)
			ctx.goNext()
		case TokenGroupMapKey:
			key, err := p.parseMapKey(ctx.withGroup(mapKeyTk.Group), mapKeyTk.Group)
			if err != nil {
				return nil, err
			}
			ctx := ctx.withChild(p.mapKeyText(key))
			colonTk := mapKeyTk.Group.Last()
			if p.isFlowMapDelim(ctx.nextToken()) {
				value, err := newNullNode(ctx, ctx.insertNullToken(colonTk))
				if err != nil {
					return nil, err
				}
				mapValue, err := newMappingValueNode(ctx, colonTk, key, value)
				if err != nil {
					return nil, err
				}
				node.Values = append(node.Values, mapValue)
				ctx.goNext()
			} else {
				ctx.goNext()
				if ctx.isTokenNotFound() {
					return nil, errors.ErrSyntax("could not find map value", colonTk.RawToken())
				}
				value, err := p.parseToken(ctx, ctx.currentToken())
				if err != nil {
					return nil, err
				}
				mapValue, err := newMappingValueNode(ctx, colonTk, key, value)
				if err != nil {
					return nil, err
				}
				node.Values = append(node.Values, mapValue)
			}
		default:
			if !p.isFlowMapDelim(ctx.nextToken()) {
				return nil, errors.ErrSyntax("could not find flow map content", mapKeyTk.RawToken())
			}
			key, err := p.parseScalarValue(ctx, mapKeyTk)
			if err != nil {
				return nil, err
			}
			value, err := newNullNode(ctx, ctx.insertNullToken(mapKeyTk))
			if err != nil {
				return nil, err
			}
			mapValue, err := newMappingValueNode(ctx, mapKeyTk, key, value)
			if err != nil {
				return nil, err
			}
			node.Values = append(node.Values, mapValue)
			ctx.goNext()
		}
		isFirst = false
	}
	if node.End == nil {
		return nil, errors.ErrSyntax("could not find flow mapping end token '}'", node.Start)
	}
	ctx.goNext() // skip mapping end token.
	return node, nil
}

func (p *parser) isFlowMapDelim(tk *Token) bool {
	return tk.Type() == token.MappingEndType || tk.Type() == token.CollectEntryType
}

func (p *parser) parseMap(ctx *context) (*ast.MappingNode, error) {
	keyTk := ctx.currentToken()
	if keyTk.Group == nil {
		return nil, errors.ErrSyntax("unexpected map key", keyTk.RawToken())
	}
	var keyValueNode *ast.MappingValueNode
	if keyTk.GroupType() == TokenGroupMapKeyValue {
		node, err := p.parseMapKeyValue(ctx.withGroup(keyTk.Group), keyTk.Group)
		if err != nil {
			return nil, err
		}
		keyValueNode = node
		ctx.goNext()
		if err := p.validateMapKeyValueNextToken(ctx, keyTk, ctx.currentToken()); err != nil {
			return nil, err
		}
	} else {
		key, err := p.parseMapKey(ctx.withGroup(keyTk.Group), keyTk.Group)
		if err != nil {
			return nil, err
		}
		ctx.goNext()

		valueTk := ctx.currentToken()
		if keyTk.Line() == valueTk.Line() && valueTk.Type() == token.SequenceEntryType {
			return nil, errors.ErrSyntax("block sequence entries are not allowed in this context", valueTk.RawToken())
		}
		ctx := ctx.withChild(p.mapKeyText(key))
		value, err := p.parseMapValue(ctx, key, keyTk.Group.Last())
		if err != nil {
			return nil, err
		}
		node, err := newMappingValueNode(ctx, keyTk.Group.Last(), key, value)
		if err != nil {
			return nil, err
		}
		keyValueNode = node
	}
	mapNode, err := newMappingNode(ctx, &Token{Token: keyValueNode.GetToken()}, false, keyValueNode)
	if err != nil {
		return nil, err
	}
	var tk *Token
	if ctx.isComment() {
		tk = ctx.nextNotCommentToken()
	} else {
		tk = ctx.currentToken()
	}
	for tk.Column() == keyTk.Column() {
		typ := tk.Type()
		if ctx.isFlow && typ == token.SequenceEndType {
			// [
			// key: value
			// ] <=
			break
		}
		if !p.isMapToken(tk) {
			return nil, errors.ErrSyntax("non-map value is specified", tk.RawToken())
		}
		cm := p.parseHeadComment(ctx)
		if typ == token.MappingEndType {
			// a: {
			//  b: c
			// } <=
			ctx.goNext()
			break
		}
		node, err := p.parseMap(ctx)
		if err != nil {
			return nil, err
		}
		if len(node.Values) != 0 {
			if err := setHeadComment(cm, node.Values[0]); err != nil {
				return nil, err
			}
		}
		mapNode.Values = append(mapNode.Values, node.Values...)
		if node.FootComment != nil {
			mapNode.Values[len(mapNode.Values)-1].FootComment = node.FootComment
		}
		tk = ctx.currentToken()
	}
	if ctx.isComment() {
		if keyTk.Column() <= ctx.currentToken().Column() {
			// If the comment is in the same or deeper column as the last element column in map value,
			// treat it as a footer comment for the last element.
			if len(mapNode.Values) == 1 {
				mapNode.Values[0].FootComment = p.parseFootComment(ctx, keyTk.Column())
				mapNode.Values[0].FootComment.SetPath(mapNode.Values[0].Key.GetPath())
			} else {
				mapNode.FootComment = p.parseFootComment(ctx, keyTk.Column())
				mapNode.FootComment.SetPath(mapNode.GetPath())
			}
		}
	}
	return mapNode, nil
}

func (p *parser) validateMapKeyValueNextToken(ctx *context, keyTk, tk *Token) error {
	if tk == nil {
		return nil
	}
	if tk.Column() <= keyTk.Column() {
		return nil
	}
	if ctx.isFlow && tk.Type() == token.CollectEntryType {
		return nil
	}
	// a: b
	//  c <= this token is invalid.
	return errors.ErrSyntax("value is not allowed in this context. map key-value is pre-defined", tk.RawToken())
}

func (p *parser) isMapToken(tk *Token) bool {
	if tk.Group == nil {
		return tk.Type() == token.MappingStartType || tk.Type() == token.MappingEndType
	}
	g := tk.Group
	return g.Type == TokenGroupMapKey || g.Type == TokenGroupMapKeyValue
}

func (p *parser) parseMapKeyValue(ctx *context, g *TokenGroup) (*ast.MappingValueNode, error) {
	if g.Type != TokenGroupMapKeyValue {
		return nil, errors.ErrSyntax("unexpected map key-value pair", g.RawToken())
	}
	if g.First().Group == nil {
		return nil, errors.ErrSyntax("unexpected map key", g.RawToken())
	}
	keyGroup := g.First().Group
	key, err := p.parseMapKey(ctx.withGroup(keyGroup), keyGroup)
	if err != nil {
		return nil, err
	}
	value, err := p.parseToken(ctx.withChild(p.mapKeyText(key)), g.Last())
	if err != nil {
		return nil, err
	}
	return newMappingValueNode(ctx, keyGroup.Last(), key, value)
}

func (p *parser) parseMapKey(ctx *context, g *TokenGroup) (ast.MapKeyNode, error) {
	if g.Type != TokenGroupMapKey {
		return nil, errors.ErrSyntax("unexpected map key", g.RawToken())
	}
	if g.First().Type() == token.MappingKeyType {
		mapKeyTk := g.First()
		if mapKeyTk.Group != nil {
			ctx = ctx.withGroup(mapKeyTk.Group)
		}
		key, err := newMappingKeyNode(ctx, mapKeyTk)
		if err != nil {
			return nil, err
		}
		ctx.goNext() // skip mapping key token
		if ctx.isTokenNotFound() {
			return nil, errors.ErrSyntax("could not find value for mapping key", mapKeyTk.RawToken())
		}

		scalar, err := p.parseScalarValue(ctx, ctx.currentToken())
		if err != nil {
			return nil, err
		}
		key.Value = scalar
		keyText := p.mapKeyText(scalar)
		keyPath := ctx.withChild(keyText).path
		key.SetPath(keyPath)
		if err := p.validateMapKey(key.GetToken(), keyPath); err != nil {
			return nil, err
		}
		p.pathMap[keyPath] = key
		return key, nil
	}
	if g.Last().Type() != token.MappingValueType {
		return nil, errors.ErrSyntax("expected map key-value delimiter ':'", g.Last().RawToken())
	}

	scalar, err := p.parseScalarValue(ctx, g.First())
	if err != nil {
		return nil, err
	}
	key, ok := scalar.(ast.MapKeyNode)
	if !ok {
		return nil, errors.ErrSyntax("cannot take map-key node", scalar.GetToken())
	}
	keyText := p.mapKeyText(key)
	keyPath := ctx.withChild(keyText).path
	key.SetPath(keyPath)
	if err := p.validateMapKey(key.GetToken(), keyPath); err != nil {
		return nil, err
	}
	p.pathMap[keyPath] = key
	return key, nil
}

func (p *parser) validateMapKey(tk *token.Token, keyPath string) error {
	if !p.allowDuplicateMapKey {
		if n, exists := p.pathMap[keyPath]; exists {
			pos := n.GetToken().Position
			return errors.ErrSyntax(
				fmt.Sprintf("mapping key %q already defined at [%d:%d]", tk.Value, pos.Line, pos.Column),
				tk,
			)
		}
	}
	if tk.Type != token.StringType {
		return nil
	}
	origin := p.removeLeftSideNewLineCharacter(tk.Origin)
	if p.existsNewLineCharacter(origin) {
		return errors.ErrSyntax("unexpected key name", tk)
	}
	return nil
}

func (p *parser) removeLeftSideNewLineCharacter(src string) string {
	// CR or LF or CRLF
	return strings.TrimLeft(strings.TrimLeft(strings.TrimLeft(src, "\r"), "\n"), "\r\n")
}

func (p *parser) existsNewLineCharacter(src string) bool {
	if strings.Index(src, "\n") > 0 {
		return true
	}
	if strings.Index(src, "\r") > 0 {
		return true
	}
	return false
}

func (p *parser) mapKeyText(n ast.Node) string {
	if n == nil {
		return ""
	}
	switch nn := n.(type) {
	case *ast.MappingKeyNode:
		return p.mapKeyText(nn.Value)
	case *ast.TagNode:
		return p.mapKeyText(nn.Value)
	case *ast.AnchorNode:
		return p.mapKeyText(nn.Value)
	case *ast.AliasNode:
		return p.mapKeyText(nn.Value)
	}
	return n.GetToken().Value
}

func (p *parser) parseMapValue(ctx *context, key ast.MapKeyNode, colonTk *Token) (ast.Node, error) {
	tk := ctx.currentToken()
	if tk == nil {
		return newNullNode(ctx, ctx.insertNullToken(colonTk))
	}

	if ctx.isComment() {
		tk = ctx.nextNotCommentToken()
	}
	keyCol := key.GetToken().Position.Column
	keyLine := key.GetToken().Position.Line

	if tk.Column() != keyCol && tk.Line() == keyLine && (tk.GroupType() == TokenGroupMapKey || tk.GroupType() == TokenGroupMapKeyValue) {
		// a: b:
		//    ^
		//
		// a: b: c
		//    ^
		return nil, errors.ErrSyntax("mapping value is not allowed in this context", tk.RawToken())
	}

	if tk.Column() == keyCol && p.isMapToken(tk) {
		// in this case,
		// ----
		// key: <value does not defined>
		// next
		return newNullNode(ctx, ctx.insertNullToken(colonTk))
	}

	if tk.Line() == keyLine && tk.GroupType() == TokenGroupAnchorName &&
		ctx.nextToken().Column() == keyCol && p.isMapToken(ctx.nextToken()) {
		// in this case,
		// ----
		// key: &anchor
		// next
		group := &TokenGroup{
			Type:   TokenGroupAnchor,
			Tokens: []*Token{tk, ctx.createNullToken(tk)},
		}
		anchor, err := p.parseAnchor(ctx.withGroup(group), group)
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return anchor, nil
	}

	if tk.Column() <= keyCol && tk.GroupType() == TokenGroupAnchorName {
		// key: <value does not defined>
		// &anchor
		return nil, errors.ErrSyntax("anchor is not allowed in this context", tk.RawToken())
	}

	if tk.Column() < keyCol {
		// in this case,
		// ----
		//   key: <value does not defined>
		// next
		return newNullNode(ctx, ctx.insertNullToken(colonTk))
	}

	if tk.Line() == keyLine && tk.GroupType() == TokenGroupAnchorName &&
		ctx.nextToken().Column() < keyCol {
		// in this case,
		// ----
		//   key: &anchor
		// next
		group := &TokenGroup{
			Type:   TokenGroupAnchor,
			Tokens: []*Token{tk, ctx.createNullToken(tk)},
		}
		anchor, err := p.parseAnchor(ctx.withGroup(group), group)
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return anchor, nil
	}

	value, err := p.parseToken(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (p *parser) parseAnchor(ctx *context, g *TokenGroup) (*ast.AnchorNode, error) {
	anchorNameGroup := g.First().Group
	anchor, err := p.parseAnchorName(ctx.withGroup(anchorNameGroup))
	if err != nil {
		return nil, err
	}
	ctx.goNext()
	if ctx.isTokenNotFound() {
		return nil, errors.ErrSyntax("could not find anchor value", anchor.GetToken())
	}

	value, err := p.parseToken(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	anchor.Value = value
	return anchor, nil
}

func (p *parser) parseAnchorName(ctx *context) (*ast.AnchorNode, error) {
	anchor, err := newAnchorNode(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	ctx.goNext()
	if ctx.isTokenNotFound() {
		return nil, errors.ErrSyntax("could not find anchor value", anchor.GetToken())
	}

	anchorName, err := p.parseScalarValue(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	if anchorName == nil {
		return nil, errors.ErrSyntax("unexpected anchor. anchor name is not scalar value", ctx.currentToken().RawToken())
	}
	anchor.Name = anchorName
	return anchor, nil
}

func (p *parser) parseAlias(ctx *context) (*ast.AliasNode, error) {
	alias, err := newAliasNode(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	ctx.goNext()
	if ctx.isTokenNotFound() {
		return nil, errors.ErrSyntax("could not find alias value", alias.GetToken())
	}

	aliasName, err := p.parseScalarValue(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	if aliasName == nil {
		return nil, errors.ErrSyntax("unexpected alias. alias name is not scalar value", ctx.currentToken().RawToken())
	}
	alias.Value = aliasName
	return alias, nil
}

func (p *parser) parseLiteral(ctx *context) (*ast.LiteralNode, error) {
	node, err := newLiteralNode(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	ctx.goNext() // skip literal/folded token

	tk := ctx.currentToken()
	if tk == nil {
		value, err := newStringNode(ctx, &Token{Token: token.New("", "", node.Start.Position)})
		if err != nil {
			return nil, err
		}
		node.Value = value
		return node, nil
	}
	value, err := p.parseToken(ctx, tk)
	if err != nil {
		return nil, err
	}
	str, ok := value.(*ast.StringNode)
	if !ok {
		return nil, errors.ErrSyntax("unexpected token. required string token", value.GetToken())
	}
	node.Value = str
	return node, nil
}

func (p *parser) parseScalarTag(ctx *context) (*ast.TagNode, error) {
	tag, err := p.parseTag(ctx)
	if err != nil {
		return nil, err
	}
	if tag.Value == nil {
		return nil, errors.ErrSyntax("specified not scalar tag", tag.GetToken())
	}
	if _, ok := tag.Value.(ast.ScalarNode); !ok {
		return nil, errors.ErrSyntax("specified not scalar tag", tag.GetToken())
	}
	return tag, nil
}

func (p *parser) parseTag(ctx *context) (*ast.TagNode, error) {
	tagTk := ctx.currentToken()
	tagRawTk := tagTk.RawToken()
	node, err := newTagNode(ctx, tagTk)
	if err != nil {
		return nil, err
	}
	ctx.goNext()

	comment := p.parseHeadComment(ctx)
	value, err := p.parseTagValue(ctx, tagRawTk, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	if err := setHeadComment(comment, value); err != nil {
		return nil, err
	}
	node.Value = value
	return node, nil
}

func (p *parser) parseTagValue(ctx *context, tagRawTk *token.Token, tk *Token) (ast.Node, error) {
	if tk == nil {
		return newNullNode(ctx, ctx.createNullToken(&Token{Token: tagRawTk}))
	}
	switch token.ReservedTagKeyword(tagRawTk.Value) {
	case token.MappingTag, token.SetTag:
		if !p.isMapToken(tk) {
			return nil, errors.ErrSyntax("could not find map", tk.RawToken())
		}
		if tk.Type() == token.MappingStartType {
			return p.parseFlowMap(ctx)
		}
		return p.parseMap(ctx)
	case token.IntegerTag, token.FloatTag, token.StringTag, token.BinaryTag, token.TimestampTag, token.BooleanTag, token.NullTag:
		if tk.GroupType() == TokenGroupLiteral || tk.GroupType() == TokenGroupFolded {
			return p.parseLiteral(ctx.withGroup(tk.Group))
		} else if tk.Type() == token.CollectEntryType || tk.Type() == token.MappingValueType {
			return newTagDefaultScalarValueNode(ctx, tagRawTk)
		}
		scalar, err := p.parseScalarValue(ctx, tk)
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return scalar, nil
	case token.SequenceTag, token.OrderedMapTag:
		if tk.Type() == token.SequenceStartType {
			return p.parseFlowSequence(ctx)
		}
		return p.parseSequence(ctx)
	}
	if strings.HasPrefix(tagRawTk.Value, "!!") {
		return nil, errors.ErrSyntax(fmt.Sprintf("unknown secondary tag name %q specified", tagRawTk.Value), tagRawTk)
	}
	return p.parseToken(ctx, tk)
}

func (p *parser) parseFlowSequence(ctx *context) (*ast.SequenceNode, error) {
	node, err := newSequenceNode(ctx, ctx.currentToken(), true)
	if err != nil {
		return nil, err
	}
	ctx.goNext() // skip SequenceStart token

	isFirst := true
	for ctx.next() {
		tk := ctx.currentToken()
		if tk.Type() == token.SequenceEndType {
			node.End = tk.RawToken()
			break
		}

		if tk.Type() == token.CollectEntryType {
			ctx.goNext()
		} else if !isFirst {
			return nil, errors.ErrSyntax("',' or ']' must be specified", tk.RawToken())
		}

		if tk := ctx.currentToken(); tk.Type() == token.SequenceEndType {
			// this case is here: "[ elem, ]".
			// In this case, ignore the last element and break sequence parsing.
			node.End = tk.RawToken()
			break
		}

		if ctx.isTokenNotFound() {
			break
		}

		value, err := p.parseToken(ctx.withIndex(uint(len(node.Values))).withFlow(true), ctx.currentToken())
		if err != nil {
			return nil, err
		}
		node.Values = append(node.Values, value)
		isFirst = false
	}
	if node.End == nil {
		return nil, errors.ErrSyntax("sequence end token ']' not found", node.Start)
	}
	ctx.goNext() // skip sequence end token.
	return node, nil
}

func (p *parser) parseSequence(ctx *context) (*ast.SequenceNode, error) {
	seqTk := ctx.currentToken()
	seqNode, err := newSequenceNode(ctx, seqTk, false)
	if err != nil {
		return nil, err
	}

	tk := seqTk
	for tk.Type() == token.SequenceEntryType && tk.Column() == seqTk.Column() {
		seqTk := tk
		comment := p.parseHeadComment(ctx)
		ctx.goNext() // skip sequence entry token

		value, err := p.parseSequenceValue(ctx.withIndex(uint(len(seqNode.Values))), seqTk)
		if err != nil {
			return nil, err
		}
		seqNode.ValueHeadComments = append(seqNode.ValueHeadComments, comment)
		seqNode.Values = append(seqNode.Values, value)

		if ctx.isComment() {
			tk = ctx.nextNotCommentToken()
		} else {
			tk = ctx.currentToken()
		}
	}
	if ctx.isComment() {
		if seqTk.Column() <= ctx.currentToken().Column() {
			// If the comment is in the same or deeper column as the last element column in sequence value,
			// treat it as a footer comment for the last element.
			seqNode.FootComment = p.parseFootComment(ctx, seqTk.Column())
			if len(seqNode.Values) != 0 {
				seqNode.FootComment.SetPath(seqNode.Values[len(seqNode.Values)-1].GetPath())
			}
		}
	}
	return seqNode, nil
}

func (p *parser) parseSequenceValue(ctx *context, seqTk *Token) (ast.Node, error) {
	tk := ctx.currentToken()
	if tk == nil {
		return newNullNode(ctx, ctx.insertNullToken(seqTk))
	}

	if ctx.isComment() {
		tk = ctx.nextNotCommentToken()
	}
	seqCol := seqTk.Column()
	seqLine := seqTk.Line()

	if tk.Column() == seqCol && tk.Type() == token.SequenceEntryType {
		// in this case,
		// ----
		// - <value does not defined>
		// -
		return newNullNode(ctx, ctx.insertNullToken(seqTk))
	}

	if tk.Line() == seqLine && tk.GroupType() == TokenGroupAnchorName &&
		ctx.nextToken().Column() == seqCol && ctx.nextToken().Type() == token.SequenceEntryType {
		// in this case,
		// ----
		// - &anchor
		// -
		group := &TokenGroup{
			Type:   TokenGroupAnchor,
			Tokens: []*Token{tk, ctx.createNullToken(tk)},
		}
		anchor, err := p.parseAnchor(ctx.withGroup(group), group)
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return anchor, nil
	}

	if tk.Column() <= seqCol && tk.GroupType() == TokenGroupAnchorName {
		// - <value does not defined>
		// &anchor
		return nil, errors.ErrSyntax("anchor is not allowed in this sequence context", tk.RawToken())
	}

	if tk.Column() < seqCol {
		// in this case,
		// ----
		//   - <value does not defined>
		// next
		return newNullNode(ctx, ctx.insertNullToken(seqTk))
	}

	if tk.Line() == seqLine && tk.GroupType() == TokenGroupAnchorName &&
		ctx.nextToken().Column() < seqCol {
		// in this case,
		// ----
		//   - &anchor
		// next
		group := &TokenGroup{
			Type:   TokenGroupAnchor,
			Tokens: []*Token{tk, ctx.createNullToken(tk)},
		}
		anchor, err := p.parseAnchor(ctx.withGroup(group), group)
		if err != nil {
			return nil, err
		}
		ctx.goNext()
		return anchor, nil
	}

	value, err := p.parseToken(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (p *parser) parseDirective(ctx *context, g *TokenGroup) (*ast.DirectiveNode, error) {
	node, err := newDirectiveNode(ctx, g.First())
	if err != nil {
		return nil, err
	}
	value, err := p.parseToken(ctx, g.Last())
	if err != nil {
		return nil, err
	}
	node.Value = value
	return node, nil
}

func (p *parser) parseComment(ctx *context) (ast.Node, error) {
	cm := p.parseHeadComment(ctx)
	node, err := p.parseToken(ctx, ctx.currentToken())
	if err != nil {
		return nil, err
	}
	if err := setHeadComment(cm, node); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) parseHeadComment(ctx *context) *ast.CommentGroupNode {
	tks := []*token.Token{}
	for ctx.isComment() {
		tks = append(tks, ctx.currentToken().RawToken())
		ctx.goNext()
	}
	if len(tks) == 0 {
		return nil
	}
	return ast.CommentGroup(tks)
}

func (p *parser) parseFootComment(ctx *context, col int) *ast.CommentGroupNode {
	tks := []*token.Token{}
	for ctx.isComment() && col <= ctx.currentToken().Column() {
		tks = append(tks, ctx.currentToken().RawToken())
		ctx.goNext()
	}
	if len(tks) == 0 {
		return nil
	}
	return ast.CommentGroup(tks)
}
