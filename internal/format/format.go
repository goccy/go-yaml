package format

import (
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/token"
)

func FormatNode(n ast.Node, existsComment bool) string {
	return newFormatter(n.GetToken(), existsComment).format(n)
}

func FormatFile(file *ast.File, existsComment bool) string {
	if len(file.Docs) == 0 {
		return ""
	}
	return newFormatter(file.Docs[0].GetToken(), existsComment).formatFile(file)
}

type Formatter struct {
	existsComment    bool
	tokenToOriginMap map[*token.Token]string
}

func newFormatter(tk *token.Token, existsComment bool) *Formatter {
	tokenToOriginMap := make(map[*token.Token]string)
	for tk.Prev != nil {
		tk = tk.Prev
	}
	tokenToOriginMap[tk] = tk.Origin

	var origin string
	for tk.Next != nil {
		tk = tk.Next
		if tk.Type == token.CommentType {
			origin += strings.Repeat("\n", strings.Count(normalizeNewLineChars(tk.Origin), "\n"))
			continue
		}
		origin += tk.Origin
		tokenToOriginMap[tk] = origin
		origin = ""
	}

	return &Formatter{
		existsComment:    existsComment,
		tokenToOriginMap: tokenToOriginMap,
	}
}

func (f *Formatter) format(n ast.Node) string {
	return f.trimSpacePrefix(
		f.trimIndentSpace(
			n.GetToken().Position.IndentNum,
			f.trimNewLineCharPrefix(f.formatNode(n)),
		),
	)
}

func (f *Formatter) formatFile(file *ast.File) string {
	if len(file.Docs) == 0 {
		return ""
	}
	var ret string
	for _, doc := range file.Docs {
		ret += f.formatDocument(doc)
	}
	return ret
}

func (f *Formatter) origin(tk *token.Token) string {
	if f.existsComment {
		return tk.Origin
	}
	return f.tokenToOriginMap[tk]
}

func (f *Formatter) formatDocument(n *ast.DocumentNode) string {
	return f.origin(n.Start) + f.formatNode(n.Body) + f.origin(n.End)
}

func (f *Formatter) formatNull(n *ast.NullNode) string {
	return f.origin(n.Token) + f.formatCommentGroup(n.Comment)
}

func (f *Formatter) formatString(n *ast.StringNode) string {
	return f.origin(n.Token) + f.formatCommentGroup(n.Comment)
}

func (f *Formatter) formatInteger(n *ast.IntegerNode) string {
	return f.origin(n.Token) + f.formatCommentGroup(n.Comment)
}

func (f *Formatter) formatFloat(n *ast.FloatNode) string {
	return f.origin(n.Token) + f.formatCommentGroup(n.Comment)
}

func (f *Formatter) formatBool(n *ast.BoolNode) string {
	return f.origin(n.Token) + f.formatCommentGroup(n.Comment)
}

func (f *Formatter) formatInfinity(n *ast.InfinityNode) string {
	return f.origin(n.Token) + f.formatCommentGroup(n.Comment)
}

func (f *Formatter) formatNan(n *ast.NanNode) string {
	return f.origin(n.Token) + f.formatCommentGroup(n.Comment)
}

func (f *Formatter) formatLiteral(n *ast.LiteralNode) string {
	return f.origin(n.Start) + f.formatCommentGroup(n.Comment) + f.origin(n.Value.Token)
}

func (f *Formatter) formatMergeKey(n *ast.MergeKeyNode) string {
	return f.origin(n.Token)
}

func (f *Formatter) formatMappingValue(n *ast.MappingValueNode) string {
	return f.formatCommentGroup(n.Comment) +
		f.formatMapKey(n.Key) + ":" + f.formatNode(n.Value) +
		f.formatCommentGroup(n.FootComment)
}

func (f *Formatter) formatMapKey(n ast.MapKeyNode) string {
	return f.formatCommentGroup(n.GetComment()) + f.formatNode(n)
}

func (f *Formatter) formatDirective(n *ast.DirectiveNode) string {
	ret := f.origin(n.Start) + f.formatNode(n.Name)
	for _, val := range n.Values {
		ret += f.formatNode(val)
	}
	return ret
}

func (f *Formatter) formatMapping(n *ast.MappingNode) string {
	if len(n.Values) == 0 {
		return "{}"
	}

	var ret string
	if n.IsFlowStyle {
		ret = f.origin(n.Start)
	}
	ret += f.formatCommentGroup(n.Comment)
	entry := n.Start
	for _, value := range n.Values {
		if n.IsFlowStyle {
			tk := value.GetToken()
			for tk.Prev != nil && tk != entry {
				tk = tk.Prev
				if tk.Type == token.CollectEntryType {
					ret += f.origin(tk)
					entry = tk
					break
				}
			}
		}
		ret += f.formatMappingValue(value)
	}
	if n.IsFlowStyle {
		ret += f.origin(n.End)
	}
	return ret
}

func (f *Formatter) formatTag(n *ast.TagNode) string {
	return f.origin(n.Start) + f.formatNode(n.Value)
}

func (f *Formatter) formatMappingKey(n *ast.MappingKeyNode) string {
	return f.origin(n.Start) + f.formatNode(n.Value)
}

func (f *Formatter) formatSequence(n *ast.SequenceNode) string {
	if len(n.Values) == 0 {
		return "[]"
	}

	var (
		ret   string
		entry = n.Start
	)
	if n.IsFlowStyle {
		ret = f.origin(n.Start)
	}
	for idx, value := range n.Values {
		tk := value.GetToken()
		for tk.Prev != nil && tk != entry {
			tk = tk.Prev
			if tk.Type == token.SequenceEntryType || tk.Type == token.CollectEntryType {
				ret += f.origin(tk)
				entry = tk
				break
			}
		}
		if len(n.ValueHeadComments) > idx {
			ret += f.formatCommentGroup(n.ValueHeadComments[idx])
		}
		ret += f.formatNode(value)
	}
	if n.IsFlowStyle {
		ret += f.origin(n.End)
	}
	ret += f.formatCommentGroup(n.FootComment)
	return ret
}

func (f *Formatter) formatAnchor(n *ast.AnchorNode) string {
	return f.origin(n.Start) + f.formatNode(n.Name) + f.formatNode(n.Value)
}

func (f *Formatter) formatAlias(n *ast.AliasNode) string {
	return f.origin(n.Start) + f.formatNode(n.Value)
}

func (f *Formatter) formatNode(n ast.Node) string {
	switch nn := n.(type) {
	case *ast.DocumentNode:
		return f.formatDocument(nn)
	case *ast.NullNode:
		return f.formatNull(nn)
	case *ast.BoolNode:
		return f.formatBool(nn)
	case *ast.IntegerNode:
		return f.formatInteger(nn)
	case *ast.FloatNode:
		return f.formatFloat(nn)
	case *ast.StringNode:
		return f.formatString(nn)
	case *ast.InfinityNode:
		return f.formatInfinity(nn)
	case *ast.NanNode:
		return f.formatNan(nn)
	case *ast.LiteralNode:
		return f.formatLiteral(nn)
	case *ast.DirectiveNode:
		return f.formatDirective(nn)
	case *ast.TagNode:
		return f.formatTag(nn)
	case *ast.MappingNode:
		return f.formatMapping(nn)
	case *ast.MappingKeyNode:
		return f.formatMappingKey(nn)
	case *ast.MappingValueNode:
		return f.formatMappingValue(nn)
	case *ast.MergeKeyNode:
		return f.formatMergeKey(nn)
	case *ast.SequenceNode:
		return f.formatSequence(nn)
	case *ast.AnchorNode:
		return f.formatAnchor(nn)
	case *ast.AliasNode:
		return f.formatAlias(nn)
	}
	return ""
}

func (f *Formatter) formatCommentGroup(g *ast.CommentGroupNode) string {
	if g == nil {
		return ""
	}
	var ret string
	for _, cm := range g.Comments {
		ret += f.formatComment(cm)
	}
	return ret
}

func (f *Formatter) formatComment(n *ast.CommentNode) string {
	if n == nil {
		return ""
	}
	return n.Token.Origin
}

// nolint: unused
func (f *Formatter) formatIndent(col int) string {
	if col <= 1 {
		return ""
	}
	return strings.Repeat(" ", col-1)
}

func (f *Formatter) trimNewLineCharPrefix(v string) string {
	return strings.TrimLeftFunc(v, func(r rune) bool {
		return r == '\n' || r == '\r'
	})
}

func (f *Formatter) trimSpacePrefix(v string) string {
	return strings.TrimLeftFunc(v, func(r rune) bool {
		return r == ' '
	})
}

func (f *Formatter) trimIndentSpace(trimIndentNum int, v string) string {
	if trimIndentNum == 0 {
		return v
	}
	lines := strings.Split(normalizeNewLineChars(v), "\n")
	out := make([]string, 0, len(lines))
	for _, line := range strings.Split(v, "\n") {
		var cnt int
		out = append(out, strings.TrimLeftFunc(line, func(r rune) bool {
			cnt++
			return r == ' ' && cnt <= trimIndentNum
		}))
	}
	return strings.Join(out, "\n")
}

// normalizeNewLineChars normalize CRLF and CR to LF.
func normalizeNewLineChars(v string) string {
	return strings.ReplaceAll(strings.ReplaceAll(v, "\r\n", "\n"), "\r", "\n")
}
