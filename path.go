package yaml

import (
	"bytes"
	"fmt"
	"io"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/internal/errors"
	"github.com/goccy/go-yaml/parser"
	"golang.org/x/xerrors"
)

var (
	ErrInvalidQuery = xerrors.New("invalid query")
	ErrInvalidPath  = xerrors.New("invalid path instance")
)

type Path struct {
	node pathNode
}

func (p *Path) String() string {
	return p.node.String()
}

func (p *Path) Read(r io.Reader, v interface{}) error {
	node, err := p.ReadNode(r)
	if err != nil {
		return errors.Wrapf(err, "failed to read node")
	}
	if err := Unmarshal([]byte(node.String()), v); err != nil {
		return errors.Wrapf(err, "failed to unmarshal")
	}
	return nil
}

func (p *Path) ReadNode(r io.Reader) (ast.Node, error) {
	if p.node == nil {
		return nil, ErrInvalidPath
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return nil, errors.Wrapf(err, "failed to copy from reader")
	}
	f, err := parser.ParseBytes(buf.Bytes(), 0)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse yaml")
	}
	for _, doc := range f.Docs {
		node, err := p.node.Filter(doc.Body)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to filter node by path ( %s )", p.node)
		}
		if node != nil {
			return node, nil
		}
	}
	return nil, nil
}

func (p *Path) Visit(node ast.Node) ast.Visitor {
	tk := node.GetToken()
	tk.Prev = nil
	tk.Next = nil
	fmt.Println(tk)
	return p
}

type pathNode interface {
	fmt.Stringer
	Chain(pathNode) pathNode
	Filter(ast.Node) (ast.Node, error)
}

type rootNode struct {
	child pathNode
}

func (n *rootNode) String() string {
	s := "$"
	if n.child != nil {
		s += n.child.String()
	}
	return s
}

func (n *rootNode) Chain(node pathNode) pathNode {
	n.child = node
	return node
}

func (n *rootNode) Filter(node ast.Node) (ast.Node, error) {
	if n.child == nil {
		return nil, nil
	}
	filtered, err := n.child.Filter(node)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to filter")
	}
	return filtered, nil
}

type selectorNode struct {
	selector string
	child    pathNode
}

func (n *selectorNode) Chain(node pathNode) pathNode {
	n.child = node
	return node
}

func (n *selectorNode) Filter(node ast.Node) (ast.Node, error) {
	switch node.Type() {
	case ast.MappingType:
		for _, value := range node.(*ast.MappingNode).Values {
			key := value.Key.GetToken().Value
			if key == n.selector {
				if n.child == nil {
					return value.Value, nil
				}
				filtered, err := n.child.Filter(value.Value)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to filter")
				}
				return filtered, nil
			}
		}
	case ast.MappingValueType:
		value := node.(*ast.MappingValueNode)
		key := value.Key.GetToken().Value
		if key == n.selector {
			if n.child == nil {
				return value.Value, nil
			}
			filtered, err := n.child.Filter(value.Value)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to filter")
			}
			return filtered, nil
		}
	default:
		return nil, errors.Wrapf(ErrInvalidQuery, "expected node type is map or map value. but got %s", node.Type())
	}
	return nil, nil
}

func (n *selectorNode) String() string {
	s := fmt.Sprintf(".%s", n.selector)
	if n.child != nil {
		s += n.child.String()
	}
	return s
}

type indexNode struct {
	selector uint
	child    pathNode
}

func (n *indexNode) Chain(node pathNode) pathNode {
	n.child = node
	return node
}

func (n *indexNode) Filter(node ast.Node) (ast.Node, error) {
	if node.Type() != ast.SequenceType {
		return nil, errors.Wrapf(ErrInvalidQuery, "expected sequence type node. but got %s", node.Type())
	}
	sequence := node.(*ast.SequenceNode)
	if n.selector >= uint(len(sequence.Values)) {
		return nil, errors.Wrapf(ErrInvalidQuery, "expected index is %d. but got sequences has %d items", n.selector, sequence.Values)
	}
	value := sequence.Values[n.selector]
	if n.child == nil {
		return value, nil
	}
	filtered, err := n.child.Filter(value)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to filter")
	}
	return filtered, nil
}

func (n *indexNode) String() string {
	s := fmt.Sprintf("[%d]", n.selector)
	if n.child != nil {
		s += n.child.String()
	}
	return s
}

type PathBuilder struct {
	root *rootNode
	node pathNode
}

func (b *PathBuilder) Root() *PathBuilder {
	root := &rootNode{}
	return &PathBuilder{root: root, node: root}
}

func (b *PathBuilder) Child(name string) *PathBuilder {
	b.node = b.node.Chain(&selectorNode{selector: name})
	return b
}

func (b *PathBuilder) Index(idx uint) *PathBuilder {
	b.node = b.node.Chain(&indexNode{selector: idx})
	return b
}

func (b *PathBuilder) Build() *Path {
	return &Path{node: b.root}
}
