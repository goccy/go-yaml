package ast

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml/token"
)

type Document struct {
	Nodes []Node
}

func (d *Document) String() string {
	values := []string{}
	for _, node := range d.Nodes {
		values = append(values, strings.TrimLeft(node.String(), " "))
	}
	return strings.Join(values, "\n")
}

type Node interface {
	String() string
	GetToken() *token.Token
}

type ScalarNode interface {
	Node
	GetValue() interface{}
}

var (
	Null = func(tk *token.Token) Node {
		return &NullNode{
			Token: tk,
		}
	}
	Bool = func(tk *token.Token) Node {
		b, _ := strconv.ParseBool(tk.Value)
		return &BoolNode{
			Token: tk,
			Value: b,
		}
	}
	Integer = func(tk *token.Token) Node {
		i, _ := strconv.ParseInt(tk.Value, 10, 64)
		return &IntegerNode{
			Token: tk,
			Value: i,
		}
	}
	Float = func(tk *token.Token) Node {
		f, _ := strconv.ParseFloat(tk.Value, 64)
		return &FloatNode{
			Token: tk,
			Value: f,
		}
	}
	Infinity = func(tk *token.Token) Node {
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
	Nan = func(tk *token.Token) Node {
		return &NanNode{Token: tk}
	}
	String = func(tk *token.Token) Node {
		return &StringNode{
			Token: tk,
			Value: tk.Value,
		}
	}
	MergeKey = func(tk *token.Token) Node {
		return &MergeKeyNode{
			Token: tk,
		}
	}
)

type NullNode struct {
	ScalarNode
	Token *token.Token
}

func (n *NullNode) GetToken() *token.Token {
	return n.Token
}

func (n *NullNode) GetValue() interface{} {
	return nil
}

func (n *NullNode) String() string {
	return "null"
}

type IntegerNode struct {
	ScalarNode
	Token *token.Token
	Value int64
}

func (n *IntegerNode) GetToken() *token.Token {
	return n.Token
}

func (n *IntegerNode) GetValue() interface{} {
	return n.Value
}

func (n *IntegerNode) String() string {
	return n.Token.Value
}

type FloatNode struct {
	ScalarNode
	Token     *token.Token
	Precision int
	Value     float64
}

func (n *FloatNode) GetToken() *token.Token {
	return n.Token
}

func (n *FloatNode) GetValue() interface{} {
	return n.Value
}

func (n *FloatNode) String() string {
	return n.Token.Value
}

type StringNode struct {
	ScalarNode
	Token *token.Token
	Value string
}

func (n *StringNode) GetToken() *token.Token {
	return n.Token
}

func (n *StringNode) GetValue() interface{} {
	return n.Value
}

func (n *StringNode) String() string {
	switch n.Token.Type {
	case token.SingleQuoteType:
		return fmt.Sprintf(`'%s'`, n.Value)
	case token.DoubleQuoteType:
		return fmt.Sprintf(`"%s"`, n.Value)
	}
	return n.Value
}

type LiteralNode struct {
	ScalarNode
	Start *token.Token
	Value *StringNode
}

func (n *LiteralNode) GetToken() *token.Token {
	return n.Start
}

func (n *LiteralNode) GetValue() interface{} {
	return n.Value.GetValue()
}

func (n *LiteralNode) String() string {
	return n.Value.String()
}

type MergeKeyNode struct {
	ScalarNode
	Token *token.Token
}

func (n *MergeKeyNode) GetToken() *token.Token {
	return n.Token
}

func (n *MergeKeyNode) GetValue() interface{} {
	return n.Token.Value
}

func (n *MergeKeyNode) String() string {
	return n.Token.Value
}

type BoolNode struct {
	ScalarNode
	Token *token.Token
	Value bool
}

func (n *BoolNode) GetToken() *token.Token {
	return n.Token
}

func (n *BoolNode) GetValue() interface{} {
	return n.Value
}

func (n *BoolNode) String() string {
	return n.Token.Value
}

type InfinityNode struct {
	ScalarNode
	Token *token.Token
	Value float64
}

func (n *InfinityNode) GetToken() *token.Token {
	return n.Token
}

func (n *InfinityNode) GetValue() interface{} {
	return n.Value
}

func (n *InfinityNode) String() string {
	return n.Token.Value
}

type NanNode struct {
	ScalarNode
	Token *token.Token
}

func (n *NanNode) GetToken() *token.Token {
	return n.Token
}

func (n *NanNode) GetValue() interface{} {
	return math.NaN()
}

func (n *NanNode) String() string {
	return n.Token.Value
}

type FlowMappingNode struct {
	Start  *token.Token
	End    *token.Token
	Values []*MappingValueNode
}

func (n *FlowMappingNode) GetToken() *token.Token {
	return n.Start
}

func (n *FlowMappingNode) String() string {
	values := []string{}
	for _, value := range n.Values {
		values = append(values, strings.TrimLeft(value.String(), " "))
	}
	return fmt.Sprintf("{%s}", strings.Join(values, ", "))
}

type MappingCollectionNode struct {
	Start  *token.Token
	Values []Node
}

func (n *MappingCollectionNode) GetToken() *token.Token {
	return n.Start
}

func (n *MappingCollectionNode) String() string {
	values := []string{}
	for _, value := range n.Values {
		values = append(values, value.String())
	}
	return strings.Join(values, "\n")
}

type MappingValueNode struct {
	Start *token.Token
	Key   Node
	Value Node
}

func (n *MappingValueNode) GetToken() *token.Token {
	return n.Start
}

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

type FlowSequenceNode struct {
	Start  *token.Token
	End    *token.Token
	Values []Node
}

func (n *FlowSequenceNode) GetToken() *token.Token {
	return n.Start
}

func (n *FlowSequenceNode) String() string {
	values := []string{}
	for _, value := range n.Values {
		values = append(values, value.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(values, ", "))
}

type SequenceNode struct {
	Start  *token.Token
	Values []Node
}

func (n *SequenceNode) GetToken() *token.Token {
	return n.Start
}

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

type AnchorNode struct {
	Start *token.Token
	Name  Node
	Value Node
}

func (n *AnchorNode) GetToken() *token.Token {
	return n.Start
}

func (n *AnchorNode) String() string {
	value := n.Value.String()
	if len(strings.Split(value, "\n")) > 1 {
		return fmt.Sprintf("&%s\n%s", n.Name.String(), value)
	}
	return fmt.Sprintf("&%s %s", n.Name.String(), value)
}

type AliasNode struct {
	Start *token.Token
	Value Node
}

func (n *AliasNode) GetToken() *token.Token {
	return n.Start
}

func (n *AliasNode) String() string {
	return fmt.Sprintf("*%s", n.Value.String())
}

type DirectiveNode struct {
	Start *token.Token
	Value Node
}

func (n *DirectiveNode) GetToken() *token.Token {
	return n.Start
}

func (n *DirectiveNode) String() string {
	return fmt.Sprintf("%s%s", n.Start.Value, n.Value.String())
}

type TagNode struct {
	Start *token.Token
	Value Node
}

func (n *TagNode) GetToken() *token.Token {
	return n.Start
}

func (n *TagNode) String() string {
	return fmt.Sprintf("%s %s", n.Start.Value, n.Value.String())
}

type Visitor interface {
	Visit(Node) Visitor
}

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
