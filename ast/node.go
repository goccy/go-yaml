package ast

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/goccy/go-yaml/token"
)

// Node represents a generic AST node
type Node interface {
	// String returns the text representation of the node
	String() string

	// Token returns token instance
	Token() *token.Token

	// Type returns type of node
	Type() NodeType

	// Release gives back the allocated structure to the system
	// so it can be reused for efficiency
	Release(recurse bool)
}

type ScalarNode interface {
	Node
	Value() interface{}
}

type AliasNode struct {
	token *token.Token
	value Node
}

type AnchorNode struct {
	token *token.Token
	name  Node
	value Node
}

type BoolNode struct {
	token *token.Token
	value bool
}

type DirectiveNode struct {
	token *token.Token
	value Node
}

type DocumentNode struct {
	start *token.Token // position of DocumentHeader ( `---` )
	end   *token.Token // position of DocumentEnd ( `...` )
	body  Node
}

type FloatNode struct {
	token      *token.Token
	preceision int
	value      float64
}

type InfinityNode struct {
	token *token.Token
	value float64
}

type IntegerNode struct {
	token *token.Token
	value interface{} // int64 or uint64
}

type LiteralNode struct {
	start *token.Token
	value *StringNode
}

type MappingNode struct {
	start       *token.Token
	end         *token.Token
	isFlowStyle bool
	values      []*MappingValueNode
}

type MappingValueNode struct {
	token *token.Token
	key   Node
	value Node
}

type MergeKeyNode struct {
	token *token.Token
}

type NanNode struct {
	token *token.Token
}

type NullNode struct {
	token *token.Token
}

type SequenceNode struct {
	start       *token.Token
	end         *token.Token
	isFlowStyle bool
	values      []Node
}

type StringNode struct {
	token *token.Token
	value string
}

type TagNode struct {
	token *token.Token
	value Node
}

// utility for code that we call a lot
func tokenValue(n Node) string {
	return n.Token().Value
}

var aliasNodePool = sync.Pool{
	New: func() interface{} { return &AliasNode{} },
}

func Alias(tk *token.Token, value Node) *AliasNode {
	n := aliasNodePool.Get().(*AliasNode)
	n.token = tk
	n.value = value
	return n
}

func (n *AliasNode) Release(recurse bool) {
	if n == nil {
		return
	}

	if recurse {
		n.value.Release(recurse)
	}
	n.token = nil
	n.value = nil
	aliasNodePool.Put(n)
}

// String returns the textual representation
func (n *AliasNode) String() string {
	return fmt.Sprintf("*%s", n.value.String())
}

// Token returns token instance
func (n *AliasNode) Token() *token.Token {
	return n.token
}

func (n *AliasNode) Value() Node {
	return n.value
}

var anchorNodePool = sync.Pool{
	New: func() interface{} { return &AnchorNode{} },
}

func Anchor(tk *token.Token, name, value Node) *AnchorNode {
	n := anchorNodePool.Get().(*AnchorNode)
	n.token = tk
	n.name = name
	n.value = value
	return n
}

func (n *AnchorNode) Name() Node {
	return n.name
}

func (n *AnchorNode) Release(recurse bool) {
	if n == nil {
		return
	}

	if recurse {
		n.name.Release(recurse)
		n.value.Release(recurse)
	}

	n.token = nil
	n.name = nil
	n.value = nil
}

// String returns the textual representation
func (n *AnchorNode) String() string {
	value := n.value.String()
	if len(strings.Split(value, "\n")) > 1 {
		return fmt.Sprintf("&%s\n%s", n.name.String(), value)
	}
	return fmt.Sprintf("&%s %s", n.name.String(), value)
}

// Token returns the associated token
func (n *AnchorNode) Token() *token.Token {
	return n.token
}

func (n *AnchorNode) Value() Node {
	return n.value
}

var boolNodePool = sync.Pool{
	New: func() interface{} { return &BoolNode{} },
}

// Bool creates a new Bool node based on the token. The value
// associated with the token will be parsed using strconv.ParseBool
// and will be used as its value
func Bool(tk *token.Token) *BoolNode {
	b, _ := strconv.ParseBool(tk.Value)
	n := boolNodePool.Get().(*BoolNode)
	n.token = tk
	n.value = b
	return n
}

func (n *BoolNode) Release(_ bool) {
	if n == nil {
		return
	}
	n.token = nil
	boolNodePool.Put(n)
}

// String returns the textual representation
func (n *BoolNode) String() string {
	return tokenValue(n)
}

// Token returns the associated token
func (n *BoolNode) Token() *token.Token {
	return n.token
}

// Value returns the associated value
func (n *BoolNode) Value() interface{} {
	return n.value
}

var directiveNodePool = sync.Pool{
	New: func() interface{} { return &DirectiveNode{} },
}

func Directive(tk *token.Token, value Node) *DirectiveNode {
	n := directiveNodePool.Get().(*DirectiveNode)
	n.token = tk
	n.value = value
	return n
}

func (n *DirectiveNode) Release(recurse bool) {
	if n == nil {
		return
	}

	if recurse {
		n.value.Release(recurse)
	}

	n.token = nil
	n.value = nil
}

// Token returns the associated token
func (n *DirectiveNode) Token() *token.Token {
	return n.token
}

// String returns the textual representation
func (n *DirectiveNode) String() string {
	return fmt.Sprintf("%s%s", n.token.Value, n.value.String())
}

var documentPool = sync.Pool{
	New: func() interface{} { return &DocumentNode{} },
}

func Document(start, end *token.Token, body Node) *DocumentNode {
	n := documentPool.Get().(*DocumentNode)
	n.start = start
	n.end = end
	n.body = body
	return n
}

func (d *DocumentNode) Body() Node {
	return d.body
}

func (d *DocumentNode) Release(recurse bool) {
	if d == nil {
		return
	}

	if recurse {
		d.body.Release(recurse)
	}
	d.start = nil
	d.end = nil
	d.body = nil
	documentPool.Put(d)
}

// Token returns token instance
func (d *DocumentNode) Token() *token.Token {
	return d.body.Token()
}

var floatNodePool = sync.Pool{
	New: func() interface{} { return &FloatNode{} },
}

// Float creates aFloat node. It parses the value associated with the
// token (note: currently no parsing error is checked)
func Float(tk *token.Token) *FloatNode {
	f, _ := strconv.ParseFloat(removeUnderScoreFromNumber(tk.Value), 64)
	n := floatNodePool.Get().(*FloatNode)
	n.token = tk
	n.value = f
	return n
}

func (n *FloatNode) Release(_ bool) {
	if n == nil {
		return
	}

	n.token = nil
}

// String returns the textual representation
func (n *FloatNode) String() string {
	return tokenValue(n)
}

// Token returns the assocaited token
func (n *FloatNode) Token() *token.Token {
	return n.token
}

// Value returns the associated value
func (n *FloatNode) Value() interface{} {
	return n.value
}

var infinityNodePool = sync.Pool{
	New: func() interface{} { return &InfinityNode{} },
}

// Infinity creates a node for .inf or -.inf value
func Infinity(tk *token.Token) *InfinityNode {
	var v float64
	switch tk.Value {
	case ".inf", ".Inf", ".INF":
		v = math.Inf(0)
	case "-.inf", "-.Inf", "-.INF":
		v = math.Inf(-1)
		// TODO: should there be a default value?
	}

	n := infinityNodePool.Get().(*InfinityNode)
	n.token = tk
	n.value = v
	return n
}

func (n *InfinityNode) Release(_ bool) {
	if n == nil {
		return
	}

	n.token = nil
}

// Token returns the associated token
func (n *InfinityNode) Token() *token.Token {
	return n.token
}

// Value returns the associated value
func (n *InfinityNode) Value() interface{} {
	return n.value
}

// String returns the textual representation
func (n *InfinityNode) String() string {
	return tokenValue(n)
}

var integerNodePool = sync.Pool{
	New: func() interface{} { return &IntegerNode{} },
}

// Integer creates a new Integer node. It parses the value associated with
// the token (note: currently no parsing error is checked)
func Integer(tk *token.Token) *IntegerNode {
	value := removeUnderScoreFromNumber(tk.Value)
	var negativePrefix string
	var base int
	var skipCharacterNum int

	unsigned := true

	switch tk.Type {
	case token.BinaryIntegerType:
		base = 2

		// skip two characters because binary token starts with '0b'
		skipCharacterNum = 2

		if value[0] == '-' {
			skipCharacterNum++
			negativePrefix = "-"
			unsigned = false
		}
	case token.OctetIntegerType:
		base = 8

		// octet token starts with '0o' or '-0o' or '0' or '-0'
		skipCharacterNum = 1

		if value[0] == '-' {
			skipCharacterNum++
			if value[2] == 'o' {
				skipCharacterNum++
			}
			negativePrefix = "-"
			unsigned = false
		} else {
			if value[1] == 'o' {
				skipCharacterNum++
			}
		}
	case token.HexIntegerType:
		base = 16

		// hex token starts with '0x' or '-0x'
		skipCharacterNum = 2

		if value[0] == '-' {
			skipCharacterNum++
			negativePrefix = "-"
			unsigned = false
		}
	default:
		base = 10
		if value[0] == '-' || value[0] == '+' {
			unsigned = false
		}
	}

	if skipCharacterNum > 0 {
		value = value[skipCharacterNum:]
	}

	var v interface{}
	if unsigned {
		i, _ := strconv.ParseUint(negativePrefix+value, base, 64)
		v = i
	} else {
		i, _ := strconv.ParseInt(negativePrefix+value, base, 64)
		v = i
	}

	n := integerNodePool.Get().(*IntegerNode)
	n.token = tk
	n.value = v
	return n
}

func (n *IntegerNode) Release(_ bool) {
	if n == nil {
		return
	}

	n.token = nil
	n.value = nil
	integerNodePool.Put(n)
}

// String returns the textual representation
func (n *IntegerNode) String() string {
	return tokenValue(n)
}

// Token returns the associated token
func (n *IntegerNode) Token() *token.Token {
	return n.token
}

// Value returns the associated value
func (n *IntegerNode) Value() interface{} {
	return n.value
}

var literalNodePool = sync.Pool{
	New: func() interface{} { return &LiteralNode{} },
}

func Literal(tk *token.Token, v *StringNode) *LiteralNode {
	n := literalNodePool.Get().(*LiteralNode)
	n.start = tk
	n.value = v
	return n
}

func (n *LiteralNode) Release(recurse bool) {
	if n == nil {
		return
	}

	if recurse {
		n.value.Release(recurse)
	}
	n.start = nil
	n.value = nil
	literalNodePool.Put(n)
}

// String returns the textual representation
func (n *LiteralNode) String() string {
	// TODO: check for n.value == nil?
	origin := n.value.Token().Origin
	return fmt.Sprintf("|\n%s", strings.TrimRight(strings.TrimRight(origin, " "), "\n"))
}

// Token returns the associated token
func (n *LiteralNode) Token() *token.Token {
	return n.start
}

// Value returns the associated value
func (n *LiteralNode) Value() interface{} {
	// TODO: check for n.value == nil?
	return n.value.Value()
}

var mappingNodePool = sync.Pool{
	New: func() interface{} { return &MappingNode{} },
}

func Mapping(start, end *token.Token, isFlowStyle bool, values ...*MappingValueNode) *MappingNode {
	n := mappingNodePool.Get().(*MappingNode)
	n.start = start
	n.end = end
	n.isFlowStyle = isFlowStyle
	n.values = values
	return n
}

func (n *MappingNode) SetFlowStyle(b bool) {
	n.isFlowStyle = b
}

func (n *MappingNode) Release(recurse bool) {
	if n == nil {
		return
	}

	if recurse {
		for _, v := range n.values {
			v.Release(recurse)
		}
	}

	n.start = nil
	n.end = nil
	n.values = nil // TODO: optimize
}

// MapRange implements MapNode protocol
func (n *MappingNode) MapRange() *MapNodeIter {
	return &MapNodeIter{
		idx:    startRangeIndex,
		values: n.values,
	}
}

func (n *MappingNode) flowStyleString() string {
	if len(n.values) == 0 {
		return "{}"
	}
	values := []string{}
	for _, value := range n.values {
		values = append(values, strings.TrimLeft(value.String(), " "))
	}
	return fmt.Sprintf("{%s}", strings.Join(values, ", "))
}

func (n *MappingNode) blockStyleString() string {
	if len(n.values) == 0 {
		return "{}"
	}
	values := []string{}
	for _, value := range n.values {
		values = append(values, value.String())
	}
	return strings.Join(values, "\n")
}

// String returns the textual representation
func (n *MappingNode) String() string {
	if n.isFlowStyle {
		return n.flowStyleString()
	}
	return n.blockStyleString()
}

// Token returns the associated token
func (n *MappingNode) Token() *token.Token {
	return n.start
}

func (n *MappingNode) Values() []*MappingValueNode {
	return n.values
}

var mappingValueNodePool = sync.Pool{
	New: func() interface{} { return &MappingValueNode{} },
}

func MappingValue(tk *token.Token, key, value Node) *MappingValueNode {
	n := mappingValueNodePool.Get().(*MappingValueNode)
	n.token = tk
	n.key = key
	n.value = value
	return n
}

func (n *MappingValueNode) Release(recurse bool) {
	if n == nil {
		return
	}

	if recurse {
		n.key.Release(recurse)
		n.value.Release(recurse)
	}

	n.token = nil
	n.key = nil
	n.value = nil
	mappingValueNodePool.Put(n)
}

func (n *MappingValueNode) Key() Node {
	return n.key
}

func (n *MappingValueNode) Value() Node {
	return n.value
}

// String returns the textual representation
func (n *MappingValueNode) String() string {
	space := strings.Repeat(" ", n.key.Token().Position.Column-1)
	keyIndentLevel := n.key.Token().Position.IndentLevel
	valueIndentLevel := n.value.Token().Position.IndentLevel
	if _, ok := n.value.(ScalarNode); ok {
		return fmt.Sprintf("%s%s: %s", space, n.key.String(), n.value.String())
	} else if keyIndentLevel < valueIndentLevel {
		return fmt.Sprintf("%s%s:\n%s", space, n.key.String(), n.value.String())
	} else if m, ok := n.value.(*MappingNode); ok && m.isFlowStyle {
		return fmt.Sprintf("%s%s: %s", space, n.key.String(), n.value.String())
	} else if s, ok := n.value.(*SequenceNode); ok && s.isFlowStyle {
		return fmt.Sprintf("%s%s: %s", space, n.key.String(), n.value.String())
	} else if _, ok := n.value.(*AnchorNode); ok {
		return fmt.Sprintf("%s%s: %s", space, n.key.String(), n.value.String())
	} else if _, ok := n.value.(*AliasNode); ok {
		return fmt.Sprintf("%s%s: %s", space, n.key.String(), n.value.String())
	}
	return fmt.Sprintf("%s%s:\n%s", space, n.key.String(), n.value.String())
}

// Token returns the associated token
func (n *MappingValueNode) Token() *token.Token {
	return n.token
}

// Type returns MappingValueType
func (n *MappingValueNode) Type() NodeType { return MappingValueType }

var mergeKeyNodePool = sync.Pool{
	New: func() interface{} { return &MergeKeyNode{} },
}

// MergeKey creates a node for merge key ( << )
func MergeKey(tk *token.Token) *MergeKeyNode {
	n := mergeKeyNodePool.Get().(*MergeKeyNode)
	n.token = tk
	return n
}

func (n *MergeKeyNode) Release(_ bool) {
	if n == nil {
		return
	}
	n.token = nil
	mergeKeyNodePool.Put(n)
}

// String returns the textual representation
func (n *MergeKeyNode) String() string {
	return tokenValue(n)
}

// Token returns the associated token
func (n *MergeKeyNode) Token() *token.Token {
	return n.token
}

// Value returns the associated value
func (n *MergeKeyNode) Value() interface{} {
	return tokenValue(n)
}

var nanNodePool = sync.Pool{
	New: func() interface{} { return &NanNode{} },
}

// Nan creates a node for .nan value
func Nan(tk *token.Token) *NanNode {
	n := nanNodePool.Get().(*NanNode)
	n.token = tk
	return n
}

func (n *NanNode) Release(_ bool) {
	if n == nil {
		return
	}
	n.token = nil
	nanNodePool.Put(n)
}

// String returns the textual representation
func (n *NanNode) String() string {
	return tokenValue(n)
}

// Token returns the associated token
func (n *NanNode) Token() *token.Token {
	return n.token
}

// Value returns math.NaN()
func (n *NanNode) Value() interface{} {
	return math.NaN()
}

var nullNodePool = sync.Pool{
	New: func() interface{} { return &NullNode{} },
}

// Null creates a new Null node based on the token
func Null(tk *token.Token) *NullNode {
	n := nullNodePool.Get().(*NullNode)
	n.token = tk
	return n
}

func (n *NullNode) Release(_ bool) {
	if n == nil {
		return
	}
	n.token = nil
	nullNodePool.Put(n)
}

// Token returns the associated token
func (n *NullNode) Token() *token.Token {
	return n.token
}

// Value returns the associated value
func (n *NullNode) Value() interface{} {
	return nil
}

var sequenceNodePool = sync.Pool{
	New: func() interface{} { return &SequenceNode{} },
}

// Sequence creates a Sequence node
func Sequence(start, end *token.Token, isFlowStyle bool, values ...Node) *SequenceNode {
	n := sequenceNodePool.Get().(*SequenceNode)
	n.start = start
	n.end = end
	n.isFlowStyle = isFlowStyle
	n.values = values
	return n
}

func (n *SequenceNode) SetFlowStyle(b bool) {
	n.isFlowStyle = b
}

func (n *SequenceNode) Release(recurse bool) {
	if n == nil {
		return
	}

	if recurse {
		for _, v := range n.values {
			v.Release(recurse)
		}
	}
	n.start = nil
	n.end = nil
	n.values = nil
}

func (n *SequenceNode) Values() []Node {
	return n.values
}

// Token returns token instance
func (n *SequenceNode) Token() *token.Token {
	return n.start
}

func (n *SequenceNode) flowStyleString() string {
	values := []string{}
	for _, value := range n.values {
		values = append(values, value.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(values, ", "))
}

func (n *SequenceNode) blockStyleString() string {
	space := strings.Repeat(" ", n.start.Position.Column-1)
	values := []string{}
	for _, value := range n.values {
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
	if n.isFlowStyle {
		return n.flowStyleString()
	}
	return n.blockStyleString()
}

// ArrayRange implements ArrayNode protocol
func (n *SequenceNode) ArrayRange() *ArrayNodeIter {
	return &ArrayNodeIter{
		idx:    startRangeIndex,
		values: n.values,
	}
}

var stringNodePool = sync.Pool{
	New: func() interface{} { return &StringNode{} },
}

// String creates a string node.
func String(tk *token.Token) *StringNode {
	n := stringNodePool.Get().(*StringNode)
	n.token = tk
	n.value = tk.Value
	return n
}

func (n *StringNode) Release(_ bool) {
	if n == nil {
		return
	}
	n.token = nil
	n.value = ""
	stringNodePool.Put(n)
}

// String string the textual representation, including quotes if necessary
func (n *StringNode) String() string {
	switch n.token.Type {
	case token.SingleQuoteType:
		return fmt.Sprintf(`'%s'`, n.value)
	case token.DoubleQuoteType:
		return fmt.Sprintf(`"%s"`, n.value)
	}
	return n.value
}

// Token returns the associated token
func (n *StringNode) Token() *token.Token {
	return n.token
}

// Value returns the associated value
func (n *StringNode) Value() interface{} {
	return n.value
}

var tagNodePool = sync.Pool{
	New: func() interface{} { return &TagNode{} },
}

func Tag(tk *token.Token, value Node) *TagNode {
	n := tagNodePool.Get().(*TagNode)
	n.token = tk
	n.value = value
	return n
}

func (n *TagNode) Release(recurse bool) {
	if n == nil {
		return
	}
	if recurse {
		n.value.Release(recurse)
	}

	n.token = nil
	n.value = nil
}

// String returns the tetual representation
func (n *TagNode) String() string {
	return fmt.Sprintf("%s %s", n.token.Value, n.value.String())
}

// Token returns the associated token
func (n *TagNode) Token() *token.Token {
	return n.token
}

func (n *TagNode) Value() Node {
	return n.value
}
