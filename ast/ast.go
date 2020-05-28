package ast

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml/token"
	"golang.org/x/xerrors"
)

var (
	ErrInvalidTokenType = xerrors.New("invalid token type")
)

// NodeType type identifier of node
type NodeType int

const (
	// UnknownNodeType type identifier for default
	UnknownNodeType NodeType = iota
	// DocumentType type identifier for document node
	DocumentType
	// NullType type identifier for null node
	NullType
	// BoolType type identifier for boolean node
	BoolType
	// IntegerType type identifier for integer node
	IntegerType
	// FloatType type identifier for float node
	FloatType
	// InfinityType type identifier for infinity node
	InfinityType
	// NanType type identifier for nan node
	NanType
	// StringType type identifier for string node
	StringType
	// MergeKeyType type identifier for merge key node
	MergeKeyType
	// LiteralType type identifier for literal node
	LiteralType
	// MappingType type identifier for mapping node
	MappingType
	// MappingValueType type identifier for mapping value node
	MappingValueType
	// SequenceType type identifier for sequence node
	SequenceType
	// AnchorType type identifier for anchor node
	AnchorType
	// AliasType type identifier for alias node
	AliasType
	// DirectiveType type identifier for directive node
	DirectiveType
	// TagType type identifier for tag node
	TagType
	// CommentType type identifier for comment node
	CommentType
)

// String node type identifier to text
func (t NodeType) String() string {
	switch t {
	case UnknownNodeType:
		return "UnknownNode"
	case DocumentType:
		return "Document"
	case NullType:
		return "Null"
	case BoolType:
		return "Bool"
	case IntegerType:
		return "Integer"
	case FloatType:
		return "Float"
	case InfinityType:
		return "Infinity"
	case NanType:
		return "Nan"
	case StringType:
		return "String"
	case MergeKeyType:
		return "MergeKey"
	case LiteralType:
		return "Literal"
	case MappingType:
		return "Mapping"
	case MappingValueType:
		return "MappingValue"
	case SequenceType:
		return "Sequence"
	case AnchorType:
		return "Anchor"
	case AliasType:
		return "Alias"
	case DirectiveType:
		return "Directive"
	case TagType:
		return "Tag"
	case CommentType:
		return "Comment"
	}
	return ""
}

// Node type of node
type Node interface {
	// String node to text
	String() string
	// GetToken returns token instance
	GetToken() *token.Token
	// Type returns type of node
	Type() NodeType
	// AddColumn add column number to child nodes recursively
	AddColumn(int)
	// SetComment set comment token to node
	SetComment(*token.Token) error
	// Comment returns comment token instance
	GetComment() *token.Token
}

// File contains all documents in YAML file
type File struct {
	Name string
	Docs []*Document
}

// String all documents to text
func (f *File) String() string {
	docs := []string{}
	for _, doc := range f.Docs {
		docs = append(docs, doc.String())
	}
	return strings.Join(docs, "\n")
}

// Document type of Document
type Document struct {
	Comment *token.Token // position of Comment ( `#comment` )
	Start   *token.Token // position of DocumentHeader ( `---` )
	End     *token.Token // position of DocumentEnd ( `...` )
	Body    Node
}

// GetToken returns token instance
func (d *Document) GetToken() *token.Token {
	return d.Body.GetToken()
}

// GetComment returns comment token instance
func (d *Document) GetComment() *token.Token {
	return d.Comment
}

// AddColumn add column number to child nodes recursively
func (d *Document) AddColumn(col int) {
	if d.Body != nil {
		d.Body.AddColumn(col)
	}
}

// SetComment set comment token
func (d *Document) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	d.Comment = tk
	return nil
}

// Type returns DocumentType
func (d *Document) Type() NodeType { return DocumentType }

// String document to text
func (d *Document) String() string {
	doc := []string{}
	if d.Start != nil {
		doc = append(doc, d.Start.Value)
	}
	doc = append(doc, d.Body.String())
	if d.End != nil {
		doc = append(doc, d.End.Value)
	}
	return strings.Join(doc, "\n")
}

// ScalarNode type for scalar node
type ScalarNode interface {
	Node
	GetValue() interface{}
}

// Null create node for null value
func Null(tk *token.Token) Node {
	return &NullNode{
		Token: tk,
	}
}

// Bool create node for boolean value
func Bool(tk *token.Token) Node {
	b, _ := strconv.ParseBool(tk.Value)
	return &BoolNode{
		Token: tk,
		Value: b,
	}
}

func removeUnderScoreFromNumber(num string) string {
	return strings.ReplaceAll(num, "_", "")
}

// Integer create node for integer value
func Integer(tk *token.Token) Node {
	value := removeUnderScoreFromNumber(tk.Value)
	switch tk.Type {
	case token.BinaryIntegerType:
		// skip two characters because binary token starts with '0b'
		skipCharacterNum := 2
		negativePrefix := ""
		if value[0] == '-' {
			skipCharacterNum++
			negativePrefix = "-"
		}
		if len(negativePrefix) > 0 {
			i, _ := strconv.ParseInt(negativePrefix+value[skipCharacterNum:], 2, 64)
			return &IntegerNode{Token: tk, Value: i}
		}
		i, _ := strconv.ParseUint(negativePrefix+value[skipCharacterNum:], 2, 64)
		return &IntegerNode{Token: tk, Value: i}
	case token.OctetIntegerType:
		// octet token starts with '0o' or '-0o' or '0' or '-0'
		skipCharacterNum := 1
		negativePrefix := ""
		if value[0] == '-' {
			skipCharacterNum++
			if value[2] == 'o' {
				skipCharacterNum++
			}
			negativePrefix = "-"
		} else {
			if value[1] == 'o' {
				skipCharacterNum++
			}
		}
		if len(negativePrefix) > 0 {
			i, _ := strconv.ParseInt(negativePrefix+value[skipCharacterNum:], 8, 64)
			return &IntegerNode{Token: tk, Value: i}
		}
		i, _ := strconv.ParseUint(value[skipCharacterNum:], 8, 64)
		return &IntegerNode{Token: tk, Value: i}
	case token.HexIntegerType:
		// hex token starts with '0x' or '-0x'
		skipCharacterNum := 2
		negativePrefix := ""
		if value[0] == '-' {
			skipCharacterNum++
			negativePrefix = "-"
		}
		if len(negativePrefix) > 0 {
			i, _ := strconv.ParseInt(negativePrefix+value[skipCharacterNum:], 16, 64)
			return &IntegerNode{Token: tk, Value: i}
		}
		i, _ := strconv.ParseUint(value[skipCharacterNum:], 16, 64)
		return &IntegerNode{Token: tk, Value: i}
	}
	if value[0] == '-' || value[0] == '+' {
		i, _ := strconv.ParseInt(value, 10, 64)
		return &IntegerNode{Token: tk, Value: i}
	}
	i, _ := strconv.ParseUint(value, 10, 64)
	return &IntegerNode{Token: tk, Value: i}
}

// Float create node for float value
func Float(tk *token.Token) Node {
	f, _ := strconv.ParseFloat(removeUnderScoreFromNumber(tk.Value), 64)
	return &FloatNode{
		Token: tk,
		Value: f,
	}
}

// Infinity create node for .inf or -.inf value
func Infinity(tk *token.Token) Node {
	node := &InfinityNode{
		Token: tk,
	}
	switch tk.Value {
	case ".inf", ".Inf", ".INF":
		node.Value = math.Inf(0)
	case "-.inf", "-.Inf", "-.INF":
		node.Value = math.Inf(-1)
	}
	return node
}

// Nan create node for .nan value
func Nan(tk *token.Token) Node {
	return &NanNode{Token: tk}
}

// String create node for string value
func String(tk *token.Token) Node {
	return &StringNode{
		Token: tk,
		Value: tk.Value,
	}
}

// Comment create node for comment
func Comment(tk *token.Token) Node {
	return &CommentNode{Comment: tk}
}

// MergeKey create node for merge key ( << )
func MergeKey(tk *token.Token) Node {
	return &MergeKeyNode{
		Token: tk,
	}
}

// Mapping create node for map
func Mapping(tk *token.Token, isFlowStyle bool) *MappingNode {
	return &MappingNode{
		Start:       tk,
		IsFlowStyle: isFlowStyle,
		Values:      []*MappingValueNode{},
	}
}

// Sequence create node for sequence
func Sequence(tk *token.Token, isFlowStyle bool) *SequenceNode {
	return &SequenceNode{
		Start:       tk,
		IsFlowStyle: isFlowStyle,
		Values:      []Node{},
	}
}

// NullNode type of null node
type NullNode struct {
	ScalarNode
	Comment *token.Token // position of Comment ( `#comment` )
	Token   *token.Token
}

// Type returns NullType
func (n *NullNode) Type() NodeType { return NullType }

// GetToken returns token instance
func (n *NullNode) GetToken() *token.Token {
	return n.Token
}

// GetComment returns comment token instance
func (n *NullNode) GetComment() *token.Token {
	return n.Comment
}

// AddColumn add column number to child nodes recursively
func (n *NullNode) AddColumn(col int) {
	n.Token.AddColumn(col)
}

// SetComment set comment token
func (n *NullNode) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	n.Comment = tk
	return nil
}

// GetValue returns nil value
func (n *NullNode) GetValue() interface{} {
	return nil
}

// String returns `null` text
func (n *NullNode) String() string {
	return "null"
}

// IntegerNode type of integer node
type IntegerNode struct {
	ScalarNode
	Comment *token.Token // position of Comment ( `#comment` )
	Token   *token.Token
	Value   interface{} // int64 or uint64 value
}

// Type returns IntegerType
func (n *IntegerNode) Type() NodeType { return IntegerType }

// GetToken returns token instance
func (n *IntegerNode) GetToken() *token.Token {
	return n.Token
}

// GetComment returns comment token instance
func (n *IntegerNode) GetComment() *token.Token {
	return n.Comment
}

// AddColumn add column number to child nodes recursively
func (n *IntegerNode) AddColumn(col int) {
	n.Token.AddColumn(col)
}

// SetComment set comment token
func (n *IntegerNode) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	n.Comment = tk
	return nil
}

// GetValue returns int64 value
func (n *IntegerNode) GetValue() interface{} {
	return n.Value
}

// String int64 to text
func (n *IntegerNode) String() string {
	return n.Token.Value
}

// FloatNode type of float node
type FloatNode struct {
	ScalarNode
	Comment   *token.Token // position of Comment ( `#comment` )
	Token     *token.Token
	Precision int
	Value     float64
}

// Type returns FloatType
func (n *FloatNode) Type() NodeType { return FloatType }

// GetToken returns token instance
func (n *FloatNode) GetToken() *token.Token {
	return n.Token
}

// GetComment returns comment token instance
func (n *FloatNode) GetComment() *token.Token {
	return n.Comment
}

// AddColumn add column number to child nodes recursively
func (n *FloatNode) AddColumn(col int) {
	n.Token.AddColumn(col)
}

// SetComment set comment token
func (n *FloatNode) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	n.Comment = tk
	return nil
}

// GetValue returns float64 value
func (n *FloatNode) GetValue() interface{} {
	return n.Value
}

// String float64 to text
func (n *FloatNode) String() string {
	return n.Token.Value
}

// StringNode type of string node
type StringNode struct {
	ScalarNode
	Comment *token.Token // position of Comment ( `#comment` )
	Token   *token.Token
	Value   string
}

// Type returns StringType
func (n *StringNode) Type() NodeType { return StringType }

// GetToken returns token instance
func (n *StringNode) GetToken() *token.Token {
	return n.Token
}

// GetComment returns comment token instance
func (n *StringNode) GetComment() *token.Token {
	return n.Comment
}

// AddColumn add column number to child nodes recursively
func (n *StringNode) AddColumn(col int) {
	n.Token.AddColumn(col)
}

// SetComment set comment token
func (n *StringNode) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	n.Comment = tk
	return nil
}

// GetValue returns string value
func (n *StringNode) GetValue() interface{} {
	return n.Value
}

// String string value to text with quote or literal header if required
func (n *StringNode) String() string {
	switch n.Token.Type {
	case token.SingleQuoteType:
		return fmt.Sprintf(`'%s'`, n.Value)
	case token.DoubleQuoteType:
		return fmt.Sprintf(`"%s"`, n.Value)
	}

	lbc := token.DetectLineBreakCharacter(n.Value)
	if strings.Contains(n.Value, lbc) {
		// This block assumes that the line breaks in this inside scalar content and the Outside scalar content are the same.
		// It works mostly, but inconsistencies occur if line break characters are mixed.
		header := token.LiteralBlockHeader(n.Value)
		space := strings.Repeat(" ", n.Token.Position.Column-1)
		values := []string{}
		for _, v := range strings.Split(n.Value, lbc) {
			values = append(values, fmt.Sprintf("%s  %s", space, v))
		}
		block := strings.TrimSuffix(strings.TrimSuffix(strings.Join(values, lbc), fmt.Sprintf("%s  %s", lbc, space)), fmt.Sprintf("  %s", space))
		return fmt.Sprintf("%s%s%s", header, lbc, block)
	} else if len(n.Value) > 0 && (n.Value[0] == '{' || n.Value[0] == '[') {
		return fmt.Sprintf(`'%s'`, n.Value)
	}
	return n.Value
}

// LiteralNode type of literal node
type LiteralNode struct {
	ScalarNode
	Comment *token.Token // position of Comment ( `#comment` )
	Start   *token.Token
	Value   *StringNode
}

// Type returns LiteralType
func (n *LiteralNode) Type() NodeType { return LiteralType }

// GetToken returns token instance
func (n *LiteralNode) GetToken() *token.Token {
	return n.Start
}

// GetComment returns comment token instance
func (n *LiteralNode) GetComment() *token.Token {
	return n.Comment
}

// AddColumn add column number to child nodes recursively
func (n *LiteralNode) AddColumn(col int) {
	n.Start.AddColumn(col)
	if n.Value != nil {
		n.Value.AddColumn(col)
	}
}

// SetComment set comment token
func (n *LiteralNode) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	n.Comment = tk
	return nil
}

// GetValue returns string value
func (n *LiteralNode) GetValue() interface{} {
	return n.String()
}

// String literal to text
func (n *LiteralNode) String() string {
	origin := n.Value.GetToken().Origin
	return fmt.Sprintf("|\n%s", strings.TrimRight(strings.TrimRight(origin, " "), "\n"))
}

// MergeKeyNode type of merge key node
type MergeKeyNode struct {
	ScalarNode
	Comment *token.Token // position of Comment ( `#comment` )
	Token   *token.Token
}

// Type returns MergeKeyType
func (n *MergeKeyNode) Type() NodeType { return MergeKeyType }

// GetToken returns token instance
func (n *MergeKeyNode) GetToken() *token.Token {
	return n.Token
}

// GetComment returns comment token instance
func (n *MergeKeyNode) GetComment() *token.Token {
	return n.Comment
}

// GetValue returns '<<' value
func (n *MergeKeyNode) GetValue() interface{} {
	return n.Token.Value
}

// String returns '<<' value
func (n *MergeKeyNode) String() string {
	return n.Token.Value
}

// AddColumn add column number to child nodes recursively
func (n *MergeKeyNode) AddColumn(col int) {
	n.Token.AddColumn(col)
}

// SetComment set comment token
func (n *MergeKeyNode) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	n.Comment = tk
	return nil
}

// BoolNode type of boolean node
type BoolNode struct {
	ScalarNode
	Comment *token.Token // position of Comment ( `#comment` )
	Token   *token.Token
	Value   bool
}

// Type returns BoolType
func (n *BoolNode) Type() NodeType { return BoolType }

// GetToken returns token instance
func (n *BoolNode) GetToken() *token.Token {
	return n.Token
}

// GetComment returns comment token instance
func (n *BoolNode) GetComment() *token.Token {
	return n.Comment
}

// AddColumn add column number to child nodes recursively
func (n *BoolNode) AddColumn(col int) {
	n.Token.AddColumn(col)
}

// SetComment set comment token
func (n *BoolNode) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	n.Comment = tk
	return nil
}

// GetValue returns boolean value
func (n *BoolNode) GetValue() interface{} {
	return n.Value
}

// String boolean to text
func (n *BoolNode) String() string {
	return n.Token.Value
}

// InfinityNode type of infinity node
type InfinityNode struct {
	ScalarNode
	Comment *token.Token // position of Comment ( `#comment` )
	Token   *token.Token
	Value   float64
}

// Type returns InfinityType
func (n *InfinityNode) Type() NodeType { return InfinityType }

// GetToken returns token instance
func (n *InfinityNode) GetToken() *token.Token {
	return n.Token
}

// GetComment returns comment token instance
func (n *InfinityNode) GetComment() *token.Token {
	return n.Comment
}

// AddColumn add column number to child nodes recursively
func (n *InfinityNode) AddColumn(col int) {
	n.Token.AddColumn(col)
}

// SetComment set comment token
func (n *InfinityNode) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	n.Comment = tk
	return nil
}

// GetValue returns math.Inf(0) or math.Inf(-1)
func (n *InfinityNode) GetValue() interface{} {
	return n.Value
}

// String infinity to text
func (n *InfinityNode) String() string {
	return n.Token.Value
}

// NanNode type of nan node
type NanNode struct {
	ScalarNode
	Comment *token.Token // position of Comment ( `#comment` )
	Token   *token.Token
}

// Type returns NanType
func (n *NanNode) Type() NodeType { return NanType }

// GetToken returns token instance
func (n *NanNode) GetToken() *token.Token {
	return n.Token
}

// GetComment returns comment token instance
func (n *NanNode) GetComment() *token.Token {
	return n.Comment
}

// AddColumn add column number to child nodes recursively
func (n *NanNode) AddColumn(col int) {
	n.Token.AddColumn(col)
}

// SetComment set comment token
func (n *NanNode) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	n.Comment = tk
	return nil
}

// GetValue returns math.NaN()
func (n *NanNode) GetValue() interface{} {
	return math.NaN()
}

// String returns .nan
func (n *NanNode) String() string {
	return n.Token.Value
}

// MapNode interface of MappingValueNode / MappingNode
type MapNode interface {
	MapRange() *MapNodeIter
}

// MapNodeIter is an iterator for ranging over a MapNode
type MapNodeIter struct {
	values []*MappingValueNode
	idx    int
}

const (
	startRangeIndex = -1
)

// Next advances the map iterator and reports whether there is another entry.
// It returns false when the iterator is exhausted.
func (m *MapNodeIter) Next() bool {
	m.idx++
	next := m.idx < len(m.values)
	return next
}

// Key returns the key of the iterator's current map node entry.
func (m *MapNodeIter) Key() Node {
	return m.values[m.idx].Key
}

// Value returns the value of the iterator's current map node entry.
func (m *MapNodeIter) Value() Node {
	return m.values[m.idx].Value
}

// MappingNode type of mapping node
type MappingNode struct {
	Comment     *token.Token // position of Comment ( `#comment` )
	Start       *token.Token
	End         *token.Token
	IsFlowStyle bool
	Values      []*MappingValueNode
}

// Type returns MappingType
func (n *MappingNode) Type() NodeType { return MappingType }

// GetToken returns token instance
func (n *MappingNode) GetToken() *token.Token {
	return n.Start
}

// GetComment returns comment token instance
func (n *MappingNode) GetComment() *token.Token {
	return n.Comment
}

// AddColumn add column number to child nodes recursively
func (n *MappingNode) AddColumn(col int) {
	n.Start.AddColumn(col)
	n.End.AddColumn(col)
	for _, value := range n.Values {
		value.AddColumn(col)
	}
}

// SetComment set comment token
func (n *MappingNode) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	n.Comment = tk
	return nil
}

func (n *MappingNode) flowStyleString() string {
	if len(n.Values) == 0 {
		return "{}"
	}
	values := []string{}
	for _, value := range n.Values {
		values = append(values, strings.TrimLeft(value.String(), " "))
	}
	return fmt.Sprintf("{%s}", strings.Join(values, ", "))
}

func (n *MappingNode) blockStyleString() string {
	if len(n.Values) == 0 {
		return "{}"
	}
	values := []string{}
	for _, value := range n.Values {
		values = append(values, value.String())
	}
	return strings.Join(values, "\n")
}

// String mapping values to text
func (n *MappingNode) String() string {
	if n.IsFlowStyle || len(n.Values) == 0 {
		return n.flowStyleString()
	}
	return n.blockStyleString()
}

// MapRange implements MapNode protocol
func (n *MappingNode) MapRange() *MapNodeIter {
	return &MapNodeIter{
		idx:    startRangeIndex,
		values: n.Values,
	}
}

// MappingValueNode type of mapping value
type MappingValueNode struct {
	Comment *token.Token // position of Comment ( `#comment` )
	Start   *token.Token
	Key     Node
	Value   Node
}

// Type returns MappingValueType
func (n *MappingValueNode) Type() NodeType { return MappingValueType }

// GetToken returns token instance
func (n *MappingValueNode) GetToken() *token.Token {
	return n.Start
}

// GetComment returns comment token instance
func (n *MappingValueNode) GetComment() *token.Token {
	return n.Comment
}

// AddColumn add column number to child nodes recursively
func (n *MappingValueNode) AddColumn(col int) {
	n.Start.AddColumn(col)
	if n.Key != nil {
		n.Key.AddColumn(col)
	}
	if n.Value != nil {
		n.Value.AddColumn(col)
	}
}

// SetComment set comment token
func (n *MappingValueNode) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	n.Comment = tk
	return nil
}

// String mapping value to text
func (n *MappingValueNode) String() string {
	space := strings.Repeat(" ", n.Key.GetToken().Position.Column-1)
	keyIndentLevel := n.Key.GetToken().Position.IndentLevel
	valueIndentLevel := n.Value.GetToken().Position.IndentLevel
	if _, ok := n.Value.(ScalarNode); ok {
		return fmt.Sprintf("%s%s: %s", space, n.Key.String(), n.Value.String())
	} else if keyIndentLevel < valueIndentLevel {
		return fmt.Sprintf("%s%s:\n%s", space, n.Key.String(), n.Value.String())
	} else if m, ok := n.Value.(*MappingNode); ok && (m.IsFlowStyle || len(m.Values) == 0) {
		return fmt.Sprintf("%s%s: %s", space, n.Key.String(), n.Value.String())
	} else if s, ok := n.Value.(*SequenceNode); ok && (s.IsFlowStyle || len(s.Values) == 0) {
		return fmt.Sprintf("%s%s: %s", space, n.Key.String(), n.Value.String())
	} else if _, ok := n.Value.(*AnchorNode); ok {
		return fmt.Sprintf("%s%s: %s", space, n.Key.String(), n.Value.String())
	} else if _, ok := n.Value.(*AliasNode); ok {
		return fmt.Sprintf("%s%s: %s", space, n.Key.String(), n.Value.String())
	}
	return fmt.Sprintf("%s%s:\n%s", space, n.Key.String(), n.Value.String())
}

// MapRange implements MapNode protocol
func (n *MappingValueNode) MapRange() *MapNodeIter {
	return &MapNodeIter{
		idx:    startRangeIndex,
		values: []*MappingValueNode{n},
	}
}

// ArrayNode interface of SequenceNode
type ArrayNode interface {
	ArrayRange() *ArrayNodeIter
}

// ArrayNodeIter is an iterator for ranging over a ArrayNode
type ArrayNodeIter struct {
	values []Node
	idx    int
}

// Next advances the array iterator and reports whether there is another entry.
// It returns false when the iterator is exhausted.
func (m *ArrayNodeIter) Next() bool {
	m.idx++
	next := m.idx < len(m.values)
	return next
}

// Value returns the value of the iterator's current array entry.
func (m *ArrayNodeIter) Value() Node {
	return m.values[m.idx]
}

// Len returns length of array
func (m *ArrayNodeIter) Len() int {
	return len(m.values)
}

// SequenceNode type of sequence node
type SequenceNode struct {
	Comment     *token.Token // position of Comment ( `#comment` )
	Start       *token.Token
	End         *token.Token
	IsFlowStyle bool
	Values      []Node
}

// Type returns SequenceType
func (n *SequenceNode) Type() NodeType { return SequenceType }

// GetToken returns token instance
func (n *SequenceNode) GetToken() *token.Token {
	return n.Start
}

// GetComment returns comment token instance
func (n *SequenceNode) GetComment() *token.Token {
	return n.Comment
}

// AddColumn add column number to child nodes recursively
func (n *SequenceNode) AddColumn(col int) {
	n.Start.AddColumn(col)
	n.End.AddColumn(col)
	for _, value := range n.Values {
		value.AddColumn(col)
	}
}

// SetComment set comment token
func (n *SequenceNode) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	n.Comment = tk
	return nil
}

func (n *SequenceNode) flowStyleString() string {
	values := []string{}
	for _, value := range n.Values {
		values = append(values, value.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(values, ", "))
}

func (n *SequenceNode) blockStyleString() string {
	space := strings.Repeat(" ", n.Start.Position.Column-1)
	values := []string{}
	for _, value := range n.Values {
		valueStr := value.String()
		splittedValues := strings.Split(valueStr, "\n")
		trimmedFirstValue := strings.TrimLeft(splittedValues[0], " ")
		diffLength := len(splittedValues[0]) - len(trimmedFirstValue)
		newValues := []string{trimmedFirstValue}
		for i := 1; i < len(splittedValues); i++ {
			trimmed := splittedValues[i][diffLength:]
			newValues = append(newValues, fmt.Sprintf("%s  %s", space, trimmed))
		}
		newValue := strings.Join(newValues, "\n")
		values = append(values, fmt.Sprintf("%s- %s", space, newValue))
	}
	return strings.Join(values, "\n")
}

// String sequence to text
func (n *SequenceNode) String() string {
	if n.IsFlowStyle || len(n.Values) == 0 {
		return n.flowStyleString()
	}
	return n.blockStyleString()
}

// ArrayRange implements ArrayNode protocol
func (n *SequenceNode) ArrayRange() *ArrayNodeIter {
	return &ArrayNodeIter{
		idx:    startRangeIndex,
		values: n.Values,
	}
}

// AnchorNode type of anchor node
type AnchorNode struct {
	Comment *token.Token // position of Comment ( `#comment` )
	Start   *token.Token
	Name    Node
	Value   Node
}

// Type returns AnchorType
func (n *AnchorNode) Type() NodeType { return AnchorType }

// GetToken returns token instance
func (n *AnchorNode) GetToken() *token.Token {
	return n.Start
}

// GetComment returns comment token instance
func (n *AnchorNode) GetComment() *token.Token {
	return n.Comment
}

// AddColumn add column number to child nodes recursively
func (n *AnchorNode) AddColumn(col int) {
	n.Start.AddColumn(col)
	if n.Name != nil {
		n.Name.AddColumn(col)
	}
	if n.Value != nil {
		n.Value.AddColumn(col)
	}
}

// SetComment set comment token
func (n *AnchorNode) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	n.Comment = tk
	return nil
}

// String anchor to text
func (n *AnchorNode) String() string {
	value := n.Value.String()
	if len(strings.Split(value, "\n")) > 1 {
		return fmt.Sprintf("&%s\n%s", n.Name.String(), value)
	} else if s, ok := n.Value.(*SequenceNode); ok && !s.IsFlowStyle {
		return fmt.Sprintf("&%s\n%s", n.Name.String(), value)
	} else if m, ok := n.Value.(*MappingNode); ok && !m.IsFlowStyle {
		return fmt.Sprintf("&%s\n%s", n.Name.String(), value)
	}
	return fmt.Sprintf("&%s %s", n.Name.String(), value)
}

// AliasNode type of alias node
type AliasNode struct {
	Comment *token.Token // position of Comment ( `#comment` )
	Start   *token.Token
	Value   Node
}

// Type returns AliasType
func (n *AliasNode) Type() NodeType { return AliasType }

// GetToken returns token instance
func (n *AliasNode) GetToken() *token.Token {
	return n.Start
}

// GetComment returns comment token instance
func (n *AliasNode) GetComment() *token.Token {
	return n.Comment
}

// AddColumn add column number to child nodes recursively
func (n *AliasNode) AddColumn(col int) {
	n.Start.AddColumn(col)
	if n.Value != nil {
		n.Value.AddColumn(col)
	}
}

// SetComment set comment token
func (n *AliasNode) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	n.Comment = tk
	return nil
}

// String alias to text
func (n *AliasNode) String() string {
	return fmt.Sprintf("*%s", n.Value.String())
}

// DirectiveNode type of directive node
type DirectiveNode struct {
	Comment *token.Token // position of Comment ( `#comment` )
	Start   *token.Token
	Value   Node
}

// Type returns DirectiveType
func (n *DirectiveNode) Type() NodeType { return DirectiveType }

// GetToken returns token instance
func (n *DirectiveNode) GetToken() *token.Token {
	return n.Start
}

// GetComment returns comment token instance
func (n *DirectiveNode) GetComment() *token.Token {
	return n.Comment
}

// AddColumn add column number to child nodes recursively
func (n *DirectiveNode) AddColumn(col int) {
	if n.Value != nil {
		n.Value.AddColumn(col)
	}
}

// SetComment set comment token
func (n *DirectiveNode) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	n.Comment = tk
	return nil
}

// String directive to text
func (n *DirectiveNode) String() string {
	return fmt.Sprintf("%s%s", n.Start.Value, n.Value.String())
}

// TagNode type of tag node
type TagNode struct {
	Comment *token.Token // position of Comment ( `#comment` )
	Start   *token.Token
	Value   Node
}

// Type returns TagType
func (n *TagNode) Type() NodeType { return TagType }

// GetToken returns token instance
func (n *TagNode) GetToken() *token.Token {
	return n.Start
}

// GetComment returns comment token instance
func (n *TagNode) GetComment() *token.Token {
	return n.Comment
}

// AddColumn add column number to child nodes recursively
func (n *TagNode) AddColumn(col int) {
	n.Start.AddColumn(col)
	if n.Value != nil {
		n.Value.AddColumn(col)
	}
}

// SetComment set comment token
func (n *TagNode) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	n.Comment = tk
	return nil
}

// String tag to text
func (n *TagNode) String() string {
	return fmt.Sprintf("%s %s", n.Start.Value, n.Value.String())
}

// CommentNode type of comment node
type CommentNode struct {
	Comment *token.Token // position of Comment ( `#comment` )
}

// Type returns TagType
func (n *CommentNode) Type() NodeType { return CommentType }

// GetToken returns token instance
func (n *CommentNode) GetToken() *token.Token { return n.Comment }

// GetComment returns comment token instance
func (n *CommentNode) GetComment() *token.Token { return n.Comment }

// AddColumn add column number to child nodes recursively
func (n *CommentNode) AddColumn(col int) {
	n.Comment.AddColumn(col)
}

// SetComment set comment token
func (n *CommentNode) SetComment(tk *token.Token) error {
	if tk.Type != token.CommentType {
		return ErrInvalidTokenType
	}
	n.Comment = tk
	return nil
}

// String comment to text
func (n *CommentNode) String() string {
	return n.Comment.Value
}

// Visitor has Visit method that is invokded for each node encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children of node with the visitor w,
// followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(Node) Visitor
}

// Walk traverses an AST in depth-first order: It starts by calling v.Visit(node); node must not be nil.
// If the visitor w returned by v.Visit(node) is not nil,
// Walk is invoked recursively with visitor w for each of the non-nil children of node,
// followed by a call of w.Visit(nil).
func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *NullNode:
	case *IntegerNode:
	case *FloatNode:
	case *StringNode:
	case *MergeKeyNode:
	case *BoolNode:
	case *InfinityNode:
	case *NanNode:
	case *MappingNode:
		for _, value := range n.Values {
			Walk(v, value)
		}
	case *MappingValueNode:
		Walk(v, n.Key)
		Walk(v, n.Value)
	case *SequenceNode:
		for _, value := range n.Values {
			Walk(v, value)
		}
	case *AnchorNode:
		Walk(v, n.Name)
		Walk(v, n.Value)
	case *AliasNode:
		Walk(v, n.Value)
	}
}
