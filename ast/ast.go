package ast

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml/token"
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
	// FlowMappingType type identifier for flow mapping node
	FlowMappingType
	// MappingCollectionType type identifier for mapping collection node
	MappingCollectionType
	// MappingValueType type identifier for mapping value node
	MappingValueType
	// FlowSequenceType type identifier for flow sequence node
	FlowSequenceType
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
	case FlowMappingType:
		return "FlowMapping"
	case MappingCollectionType:
		return "MappingCollection"
	case MappingValueType:
		return "MappingValue"
	case FlowSequenceType:
		return "FlowSequence"
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
}

// Document type of Document
type Document struct {
	// Nodes all nodes in document
	Nodes []Node
}

// Type returns DocumentType
func (d *Document) Type() NodeType { return DocumentType }

// String document to text
func (d *Document) String() string {
	values := []string{}
	for _, node := range d.Nodes {
		values = append(values, strings.TrimLeft(node.String(), " "))
	}
	return strings.Join(values, "\n")
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

// Integer create node for integer value
func Integer(tk *token.Token) Node {
	i, _ := strconv.ParseInt(tk.Value, 10, 64)
	return &IntegerNode{
		Token: tk,
		Value: i,
	}
}

// Float create node for float value
func Float(tk *token.Token) Node {
	f, _ := strconv.ParseFloat(tk.Value, 64)
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
	case ".inf":
		node.Value = math.Inf(0)
	case "-.inf":
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

// MergeKey create node for merge key ( << )
func MergeKey(tk *token.Token) Node {
	return &MergeKeyNode{
		Token: tk,
	}
}

// NullNode type of null node
type NullNode struct {
	ScalarNode
	Token *token.Token
}

// Type returns NullType
func (n *NullNode) Type() NodeType { return NullType }

// GetToken returns token instance
func (n *NullNode) GetToken() *token.Token {
	return n.Token
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
	Token *token.Token
	Value int64
}

// Type returns IntegerType
func (n *IntegerNode) Type() NodeType { return IntegerType }

// GetToken returns token instance
func (n *IntegerNode) GetToken() *token.Token {
	return n.Token
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
	Token *token.Token
	Value string
}

// Type returns StringType
func (n *StringNode) Type() NodeType { return StringType }

// GetToken returns token instance
func (n *StringNode) GetToken() *token.Token {
	return n.Token
}

// GetValue returns string value
func (n *StringNode) GetValue() interface{} {
	return n.Value
}

// String string value to text with quote if required
func (n *StringNode) String() string {
	switch n.Token.Type {
	case token.SingleQuoteType:
		return fmt.Sprintf(`'%s'`, n.Value)
	case token.DoubleQuoteType:
		return fmt.Sprintf(`"%s"`, n.Value)
	}
	return n.Value
}

// LiteralNode type of literal node
type LiteralNode struct {
	ScalarNode
	Start *token.Token
	Value *StringNode
}

// Type returns LiteralType
func (n *LiteralNode) Type() NodeType { return LiteralType }

// GetToken returns token instance
func (n *LiteralNode) GetToken() *token.Token {
	return n.Start
}

// GetValue returns string value
func (n *LiteralNode) GetValue() interface{} {
	return n.Value.GetValue()
}

// String literal to text
func (n *LiteralNode) String() string {
	return n.Value.String()
}

// MergeKeyNode type of merge key node
type MergeKeyNode struct {
	ScalarNode
	Token *token.Token
}

// Type returns MergeKeyType
func (n *MergeKeyNode) Type() NodeType { return MergeKeyType }

// GetToken returns token instance
func (n *MergeKeyNode) GetToken() *token.Token {
	return n.Token
}

// GetValue returns '<<' value
func (n *MergeKeyNode) GetValue() interface{} {
	return n.Token.Value
}

// String returns '<<' value
func (n *MergeKeyNode) String() string {
	return n.Token.Value
}

// BoolNode type of boolean node
type BoolNode struct {
	ScalarNode
	Token *token.Token
	Value bool
}

// Type returns BoolType
func (n *BoolNode) Type() NodeType { return BoolType }

// GetToken returns token instance
func (n *BoolNode) GetToken() *token.Token {
	return n.Token
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
	Token *token.Token
	Value float64
}

// Type returns InfinityType
func (n *InfinityNode) Type() NodeType { return InfinityType }

// GetToken returns token instance
func (n *InfinityNode) GetToken() *token.Token {
	return n.Token
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
	Token *token.Token
}

// Type returns NanType
func (n *NanNode) Type() NodeType { return NanType }

// GetToken returns token instance
func (n *NanNode) GetToken() *token.Token {
	return n.Token
}

// GetValue returns math.NaN()
func (n *NanNode) GetValue() interface{} {
	return math.NaN()
}

// String returns .nan
func (n *NanNode) String() string {
	return n.Token.Value
}

// FlowMappingNode type of flow mapping node
type FlowMappingNode struct {
	Start  *token.Token
	End    *token.Token
	Values []*MappingValueNode
}

// Type returns FlowMappingType
func (n *FlowMappingNode) Type() NodeType { return FlowMappingType }

// GetToken returns token instance
func (n *FlowMappingNode) GetToken() *token.Token {
	return n.Start
}

// String flow mapping to text
func (n *FlowMappingNode) String() string {
	values := []string{}
	for _, value := range n.Values {
		values = append(values, strings.TrimLeft(value.String(), " "))
	}
	return fmt.Sprintf("{%s}", strings.Join(values, ", "))
}

// MappingCollectionNode type of mapping collection node
type MappingCollectionNode struct {
	Start  *token.Token
	Values []Node
}

// Type returns MappingCollectionType
func (n *MappingCollectionNode) Type() NodeType { return MappingCollectionType }

// GetToken returns token instance
func (n *MappingCollectionNode) GetToken() *token.Token {
	return n.Start
}

// String mapping collection to text
func (n *MappingCollectionNode) String() string {
	values := []string{}
	for _, value := range n.Values {
		values = append(values, value.String())
	}
	return strings.Join(values, "\n")
}

// MappingValueNode type of mapping value
type MappingValueNode struct {
	Start *token.Token
	Key   Node
	Value Node
}

// Type returns MappingValueType
func (n *MappingValueNode) Type() NodeType { return MappingValueType }

// GetToken returns token instance
func (n *MappingValueNode) GetToken() *token.Token {
	return n.Start
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
	} else if _, ok := n.Value.(*FlowSequenceNode); ok {
		return fmt.Sprintf("%s%s: %s", space, n.Key.String(), n.Value.String())
	} else if _, ok := n.Value.(*AnchorNode); ok {
		return fmt.Sprintf("%s%s: %s", space, n.Key.String(), n.Value.String())
	} else if _, ok := n.Value.(*AliasNode); ok {
		return fmt.Sprintf("%s%s: %s", space, n.Key.String(), n.Value.String())
	}
	return fmt.Sprintf("%s%s:\n%s", space, n.Key.String(), n.Value.String())
}

// FlowSequenceNode type of sequence node
type FlowSequenceNode struct {
	Start  *token.Token
	End    *token.Token
	Values []Node
}

// Type returns FlowSequenceType
func (n *FlowSequenceNode) Type() NodeType { return FlowSequenceType }

// GetToken returns token instance
func (n *FlowSequenceNode) GetToken() *token.Token {
	return n.Start
}

// String flow sequence to text
func (n *FlowSequenceNode) String() string {
	values := []string{}
	for _, value := range n.Values {
		values = append(values, value.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(values, ", "))
}

// SequenceNode type of sequence node
type SequenceNode struct {
	Start  *token.Token
	Values []Node
}

// Type returns SequenceType
func (n *SequenceNode) Type() NodeType { return SequenceType }

// GetToken returns token instance
func (n *SequenceNode) GetToken() *token.Token {
	return n.Start
}

// String sequence to text
func (n *SequenceNode) String() string {
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

// AnchorNode type of anchor node
type AnchorNode struct {
	Start *token.Token
	Name  Node
	Value Node
}

// Type returns AnchorType
func (n *AnchorNode) Type() NodeType { return AnchorType }

// GetToken returns token instance
func (n *AnchorNode) GetToken() *token.Token {
	return n.Start
}

// String anchor to text
func (n *AnchorNode) String() string {
	value := n.Value.String()
	if len(strings.Split(value, "\n")) > 1 {
		return fmt.Sprintf("&%s\n%s", n.Name.String(), value)
	}
	return fmt.Sprintf("&%s %s", n.Name.String(), value)
}

// AliasNode type of alias node
type AliasNode struct {
	Start *token.Token
	Value Node
}

// Type returns AliasType
func (n *AliasNode) Type() NodeType { return AliasType }

// GetToken returns token instance
func (n *AliasNode) GetToken() *token.Token {
	return n.Start
}

// String alias to text
func (n *AliasNode) String() string {
	return fmt.Sprintf("*%s", n.Value.String())
}

// DirectiveNode type of directive node
type DirectiveNode struct {
	Start *token.Token
	Value Node
}

// Type returns DirectiveType
func (n *DirectiveNode) Type() NodeType { return DirectiveType }

// GetToken returns token instance
func (n *DirectiveNode) GetToken() *token.Token {
	return n.Start
}

// String directive to text
func (n *DirectiveNode) String() string {
	return fmt.Sprintf("%s%s", n.Start.Value, n.Value.String())
}

// TagNode type of tag node
type TagNode struct {
	Start *token.Token
	Value Node
}

// Type returns TagType
func (n *TagNode) Type() NodeType { return TagType }

// GetToken returns token instance
func (n *TagNode) GetToken() *token.Token {
	return n.Start
}

// String tag to text
func (n *TagNode) String() string {
	return fmt.Sprintf("%s %s", n.Start.Value, n.Value.String())
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
	case *FlowMappingNode:
		for _, value := range n.Values {
			Walk(v, value)
		}
	case *MappingCollectionNode:
		for _, value := range n.Values {
			Walk(v, value)
		}
	case *MappingValueNode:
		Walk(v, n.Key)
		Walk(v, n.Value)
	case *FlowSequenceNode:
		for _, value := range n.Values {
			Walk(v, value)
		}
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
