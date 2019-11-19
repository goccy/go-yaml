package ast

import (
	"strings"
	"sync"
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
	}
	return ""
}

// FileNode contains all documents in YAML file
type FileNode struct {
	name string
	docs []*DocumentNode
}

var fileNodePool = sync.Pool{
	New: func() interface{} { return &FileNode{} },
}

func File(docs ...*DocumentNode) *FileNode {
	n := fileNodePool.Get().(*FileNode)
	n.docs = docs
	return n
}

func (n *FileNode) Documents() []*DocumentNode {
	return n.docs
}

func (n *FileNode) Release(recurse bool) {
	if n == nil {
		return
	}

	if recurse {
		for _, v := range n.docs {
			v.Release(recurse)
		}
	}
	n.docs = nil
	n.name = ""
	fileNodePool.Put(n)
}

func (f *FileNode) SetName(n string) {
	f.name = n
}

// String all documents to text
func (f *FileNode) String() string {
	docs := []string{}
	for _, doc := range f.docs {
		docs = append(docs, doc.String())
	}
	return strings.Join(docs, "\n")
}

// Type returns DocumentType
func (d *DocumentNode) Type() NodeType { return DocumentType }

// String document to text
func (d *DocumentNode) String() string {
	doc := []string{}
	if d.start != nil {
		doc = append(doc, d.start.Value)
	}
	doc = append(doc, d.body.String())
	if d.end != nil {
		doc = append(doc, d.end.Value)
	}
	return strings.Join(doc, "\n")
}

func removeUnderScoreFromNumber(num string) string {
	return strings.ReplaceAll(num, "_", "")
}

// Type returns NullType
func (n *NullNode) Type() NodeType { return NullType }

// String returns `null` text
func (n *NullNode) String() string {
	return "null"
}

// Type returns IntegerType
func (n *IntegerNode) Type() NodeType { return IntegerType }

// Type returns FloatType
func (n *FloatNode) Type() NodeType { return FloatType }

// Type returns StringType
func (n *StringNode) Type() NodeType { return StringType }

// Type returns LiteralType
func (n *LiteralNode) Type() NodeType { return LiteralType }

// Type returns MergeKeyType
func (n *MergeKeyNode) Type() NodeType { return MergeKeyType }

// Type returns BoolType
func (n *BoolNode) Type() NodeType { return BoolType }

// Type returns InfinityType
func (n *InfinityNode) Type() NodeType { return InfinityType }

// Type returns NanType
func (n *NanNode) Type() NodeType { return NanType }

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
	return m.values[m.idx].Key()
}

// Value returns the value of the iterator's current map node entry.
func (m *MapNodeIter) Value() Node {
	return m.values[m.idx].Value()
}

// Type returns MappingType
func (n *MappingNode) Type() NodeType { return MappingType }

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

// Type returns SequenceType
func (n *SequenceNode) Type() NodeType { return SequenceType }

// Type returns AnchorType
func (n *AnchorNode) Type() NodeType { return AnchorType }

// Type returns AliasType
func (n *AliasNode) Type() NodeType { return AliasType }

// Type returns DirectiveType
func (n *DirectiveNode) Type() NodeType { return DirectiveType }

// Type returns TagType
func (n *TagNode) Type() NodeType { return TagType }

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
		for _, value := range n.Values() {
			Walk(v, value)
		}
	case *MappingValueNode:
		Walk(v, n.Key())
		Walk(v, n.Value())
	case *SequenceNode:
		for _, value := range n.Values() {
			Walk(v, value)
		}
	case *AnchorNode:
		Walk(v, n.Name())
		Walk(v, n.Value())
	case *AliasNode:
		Walk(v, n.Value())
	}
}
