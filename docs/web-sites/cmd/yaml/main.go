package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
)

func Tokenize(v string) ([]byte, error) {
	tks := lexer.Tokenize(v)
	b, err := json.Marshal(tks)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Parse(ctx context.Context, v string) ([]byte, error) {
	gv, err := graphviz.New(ctx)
	if err != nil {
		return nil, err
	}
	file, err := parser.ParseBytes([]byte(v), parser.ParseComments)
	if err != nil {
		return nil, err
	}

	graph, err := gv.Graph()
	if err != nil {
		return nil, err
	}
	graph.SetCompound(true)
	defer func() {
		if err := graph.Close(); err != nil {
			panic(err)
		}
		gv.Close()
	}()
	renderer := &NodeRenderer{}
	if _, err := renderer.renderFile(graph, file); err != nil {
		return nil, err
	}
	var xdot bytes.Buffer
	if err := gv.Render(ctx, graph, graphviz.XDOT, &xdot); err != nil {
		return nil, err
	}
	fmt.Println(xdot.String())

	var out bytes.Buffer
	if err := gv.Render(ctx, graph, graphviz.SVG, &out); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

type NodeRenderer struct {
	id    int
	edges []*Edge
}

type Edge struct {
	start *Node
	end   *Node
}

type Node struct {
	graphName string
	node      *graphviz.Node
}

func (r *NodeRenderer) createID() string {
	r.id++
	return fmt.Sprint(r.id)
}

func (r *NodeRenderer) createNodeGraph(parent *graphviz.Graph, node any, name string) (*graphviz.Graph, error) {
	id := r.createID()
	sub, err := parent.CreateSubGraphByName("cluster" + id)
	if err != nil {
		return nil, err
	}
	sub.SetCompound(true)
	sub.SetLabel(name)
	sub.SetStyle(graphviz.FilledGraphStyle)
	sub.SetBackgroundColor("white")
	return sub, nil
}

func (r *NodeRenderer) createNode(graph *graphviz.Graph, name string) (*graphviz.Node, error) {
	node, err := graph.CreateNodeByName(r.createID())
	if err != nil {
		return nil, err
	}
	node.SetLabel(name)
	node.SetStyle(graphviz.FilledNodeStyle)
	node.SetFillColor("white")
	return node, nil
}

func (r *NodeRenderer) createEdge(fromGraph *graphviz.Graph, fromNode *graphviz.Node, toGraph *graphviz.Graph) error {
	to, err := toGraph.FirstNode()
	if err != nil {
		return err
	}
	edge, err := fromGraph.CreateEdgeByName("", fromNode, to)
	if err != nil {
		return err
	}
	fromGraphName, err := fromGraph.Name()
	if err != nil {
		return err
	}
	toGraphName, err := toGraph.Name()
	if err != nil {
		return err
	}
	edge.SetLogicalTail(fromGraphName)
	edge.SetLogicalHead(toGraphName)
	edge.SetMinLen(2)
	return nil
}

func (r *NodeRenderer) renderFile(graph *graphviz.Graph, file *ast.File) (*graphviz.Graph, error) {
	fileGraph, err := r.createNodeGraph(graph, file, "FileNode")
	if err != nil {
		return nil, err
	}
	fileGraph.SetBackgroundColor("ivory")
	for idx, doc := range file.Docs {
		node, err := r.createNode(fileGraph, fmt.Sprintf("docs[%d]", idx))
		if err != nil {
			return nil, err
		}
		docGraph, err := r.renderDocument(fileGraph, doc)
		if err != nil {
			return nil, err
		}
		if err := r.createEdge(fileGraph, node, docGraph); err != nil {
			return nil, err
		}
	}
	return fileGraph, nil
}

func (r *NodeRenderer) renderDocument(graph *graphviz.Graph, doc *ast.DocumentNode) (*graphviz.Graph, error) {
	docGraph, err := r.createNodeGraph(graph, doc, "DocumentNode")
	if err != nil {
		return nil, err
	}
	docGraph.SetBackgroundColor("mintcream")
	if err := r.renderToken(docGraph, doc.Start); err != nil {
		return nil, err
	}
	if err := r.renderToken(docGraph, doc.End); err != nil {
		return nil, err
	}
	body, err := r.createNode(docGraph, "body")
	if err != nil {
		return nil, err
	}
	bodyGraph, err := r.renderNode(docGraph, doc.Body)
	if err != nil {
		return nil, err
	}
	if err := r.createEdge(docGraph, body, bodyGraph); err != nil {
		return nil, err
	}
	return docGraph, nil
}

func (r *NodeRenderer) renderNode(graph *graphviz.Graph, node ast.Node) (*graphviz.Graph, error) {
	switch n := node.(type) {
	case *ast.MappingNode:
		return r.renderMappingNode(graph, n)
	case *ast.NullNode:
		return r.renderNullNode(graph, n)
	case *ast.IntegerNode:
		return r.renderIntegerNode(graph, n)
	case *ast.FloatNode:
		return r.renderFloatNode(graph, n)
	case *ast.StringNode:
		return r.renderStringNode(graph, n)
	case *ast.LiteralNode:
		return r.renderLiteralNode(graph, n)
	case *ast.MergeKeyNode:
		return r.renderMergeKeyNode(graph, n)
	case *ast.BoolNode:
		return r.renderBoolNode(graph, n)
	case *ast.InfinityNode:
		return r.renderInfinityNode(graph, n)
	case *ast.NanNode:
		return r.renderNaNNode(graph, n)
	case *ast.MappingKeyNode:
		return r.renderMappingKeyNode(graph, n)
	case *ast.SequenceNode:
		return r.renderSequenceNode(graph, n)
	case *ast.AnchorNode:
		return r.renderAnchorNode(graph, n)
	case *ast.AliasNode:
		return r.renderAliasNode(graph, n)
	case *ast.DirectiveNode:
		return r.renderDirectiveNode(graph, n)
	case *ast.TagNode:
		return r.renderTagNode(graph, n)
	case *ast.CommentNode:
	case *ast.CommentGroupNode:
	}
	return nil, fmt.Errorf("unexpected node type %T", node)
}

func (r *NodeRenderer) renderMappingNode(graph *graphviz.Graph, node *ast.MappingNode) (*graphviz.Graph, error) {
	mapGraph, err := r.createNodeGraph(graph, node, "MappingNode")
	if err != nil {
		return nil, err
	}
	mapGraph.SetBackgroundColor("honeydew")
	if err := r.renderPath(mapGraph, node.GetPath()); err != nil {
		return nil, err
	}
	if err := r.renderToken(mapGraph, node.Start); err != nil {
		return nil, err
	}
	if err := r.renderToken(mapGraph, node.End); err != nil {
		return nil, err
	}
	if _, err := r.renderCommentGroupNode(mapGraph, node.GetComment(), yaml.CommentHeadPosition); err != nil {
		return nil, err
	}
	if _, err := r.renderCommentGroupNode(mapGraph, node.FootComment, yaml.CommentFootPosition); err != nil {
		return nil, err
	}
	for idx, value := range node.Values {
		node, err := r.createNode(mapGraph, fmt.Sprintf("values[%d]", idx))
		if err != nil {
			return nil, err
		}
		valueGraph, err := r.renderMappingValueNode(mapGraph, value)
		if err != nil {
			return nil, err
		}
		if err := r.createEdge(mapGraph, node, valueGraph); err != nil {
			return nil, err
		}
	}
	return mapGraph, nil
}

func (r *NodeRenderer) renderMappingValueNode(graph *graphviz.Graph, node *ast.MappingValueNode) (*graphviz.Graph, error) {
	valueGraph, err := r.createNodeGraph(graph, node, "MappingValueNode")
	if err != nil {
		return nil, err
	}
	valueGraph.SetBackgroundColor("seashell")
	if err := r.renderPath(valueGraph, node.GetPath()); err != nil {
		return nil, err
	}
	if err := r.renderToken(valueGraph, node.Start); err != nil {
		return nil, err
	}
	if _, err := r.renderCommentGroupNode(valueGraph, node.GetComment(), yaml.CommentHeadPosition); err != nil {
		return nil, err
	}
	if _, err := r.renderCommentGroupNode(valueGraph, node.FootComment, yaml.CommentFootPosition); err != nil {
		return nil, err
	}
	keyNode, err := r.createNode(valueGraph, "key")
	if err != nil {
		return nil, err
	}
	valueNode, err := r.createNode(valueGraph, "value")
	if err != nil {
		return nil, err
	}
	keyGraph, err := r.renderNode(valueGraph, node.Key)
	if err != nil {
		return nil, err
	}
	valueContentGraph, err := r.renderNode(valueGraph, node.Value)
	if err != nil {
		return nil, err
	}
	if err := r.createEdge(valueGraph, keyNode, keyGraph); err != nil {
		return nil, err
	}
	if err := r.createEdge(valueGraph, valueNode, valueContentGraph); err != nil {
		return nil, err
	}
	return valueGraph, nil
}

func (r *NodeRenderer) renderStringNode(graph *graphviz.Graph, n *ast.StringNode) (*graphviz.Graph, error) {
	subGraph, err := r.createNodeGraph(graph, n, "StringNode")
	if err != nil {
		return nil, err
	}
	subGraph.SetBackgroundColor("lavenderblush")
	if err := r.renderPath(subGraph, n.GetPath()); err != nil {
		return nil, err
	}
	if err := r.renderToken(subGraph, n.Token); err != nil {
		return nil, err
	}
	if _, err := r.renderCommentGroupNode(subGraph, n.GetComment(), yaml.CommentLinePosition); err != nil {
		return nil, err
	}
	return subGraph, nil
}

func (r *NodeRenderer) renderNullNode(graph *graphviz.Graph, n *ast.NullNode) (*graphviz.Graph, error) {
	subGraph, err := r.createNodeGraph(graph, n, "NullNode")
	if err != nil {
		return nil, err
	}
	subGraph.SetBackgroundColor("lavenderblush")
	if err := r.renderPath(subGraph, n.GetPath()); err != nil {
		return nil, err
	}
	if _, err := r.renderCommentGroupNode(subGraph, n.GetComment(), yaml.CommentLinePosition); err != nil {
		return nil, err
	}
	if _, err := r.createNode(subGraph, "null"); err != nil {
		return nil, err
	}
	return subGraph, nil
}

func (r *NodeRenderer) renderIntegerNode(graph *graphviz.Graph, n *ast.IntegerNode) (*graphviz.Graph, error) {
	subGraph, err := r.createNodeGraph(graph, n, "IntegerNode")
	if err != nil {
		return nil, err
	}
	subGraph.SetBackgroundColor("lavenderblush")
	if err := r.renderPath(subGraph, n.GetPath()); err != nil {
		return nil, err
	}
	if _, err := r.renderCommentGroupNode(subGraph, n.GetComment(), yaml.CommentLinePosition); err != nil {
		return nil, err
	}
	if _, err := r.createNode(subGraph, n.Token.Value); err != nil {
		return nil, err
	}
	return subGraph, nil
}

func (r *NodeRenderer) renderFloatNode(graph *graphviz.Graph, n *ast.FloatNode) (*graphviz.Graph, error) {
	subGraph, err := r.createNodeGraph(graph, n, "FloatNode")
	if err != nil {
		return nil, err
	}
	subGraph.SetBackgroundColor("lavenderblush")
	if err := r.renderPath(subGraph, n.GetPath()); err != nil {
		return nil, err
	}
	if _, err := r.createNode(subGraph, n.Token.Value); err != nil {
		return nil, err
	}
	return subGraph, nil
}

func (r *NodeRenderer) renderLiteralNode(graph *graphviz.Graph, n *ast.LiteralNode) (*graphviz.Graph, error) {
	subGraph, err := r.createNodeGraph(graph, n, "LiteralNode")
	if err != nil {
		return nil, err
	}
	subGraph.SetBackgroundColor("beige")
	if err := r.renderPath(subGraph, n.GetPath()); err != nil {
		return nil, err
	}
	if _, err := r.renderCommentGroupNode(subGraph, n.GetComment(), yaml.CommentLinePosition); err != nil {
		return nil, err
	}
	value, err := r.createNode(subGraph, "value")
	if err != nil {
		return nil, err
	}
	strGraph, err := r.renderStringNode(subGraph, n.Value)
	if err != nil {
		return nil, err
	}
	if err := r.createEdge(subGraph, value, strGraph); err != nil {
		return nil, err
	}
	return subGraph, nil
}

func (r *NodeRenderer) renderMergeKeyNode(graph *graphviz.Graph, n *ast.MergeKeyNode) (*graphviz.Graph, error) {
	subGraph, err := r.createNodeGraph(graph, n, "MergeKeyNode")
	if err != nil {
		return nil, err
	}
	subGraph.SetBackgroundColor("cornsilk")
	if err := r.renderPath(subGraph, n.GetPath()); err != nil {
		return nil, err
	}
	if _, err := r.renderCommentGroupNode(subGraph, n.GetComment(), yaml.CommentHeadPosition); err != nil {
		return nil, err
	}
	if _, err := r.createNode(subGraph, n.Token.Value); err != nil {
		return nil, err
	}
	return subGraph, nil
}

func (r *NodeRenderer) renderBoolNode(graph *graphviz.Graph, n *ast.BoolNode) (*graphviz.Graph, error) {
	subGraph, err := r.createNodeGraph(graph, n, "BoolNode")
	if err != nil {
		return nil, err
	}
	subGraph.SetBackgroundColor("lavenderblush")
	if err := r.renderPath(subGraph, n.GetPath()); err != nil {
		return nil, err
	}
	if _, err := r.renderCommentGroupNode(subGraph, n.GetComment(), yaml.CommentLinePosition); err != nil {
		return nil, err
	}
	if _, err := r.createNode(subGraph, n.Token.Value); err != nil {
		return nil, err
	}
	return subGraph, nil
}

func (r *NodeRenderer) renderInfinityNode(graph *graphviz.Graph, n *ast.InfinityNode) (*graphviz.Graph, error) {
	subGraph, err := r.createNodeGraph(graph, n, "InfinityNode")
	if err != nil {
		return nil, err
	}
	subGraph.SetBackgroundColor("lavenderblush")
	if err := r.renderPath(subGraph, n.GetPath()); err != nil {
		return nil, err
	}
	if _, err := r.renderCommentGroupNode(subGraph, n.GetComment(), yaml.CommentLinePosition); err != nil {
		return nil, err
	}
	if _, err := r.createNode(subGraph, n.Token.Value); err != nil {
		return nil, err
	}
	return subGraph, nil
}

func (r *NodeRenderer) renderNaNNode(graph *graphviz.Graph, n *ast.NanNode) (*graphviz.Graph, error) {
	subGraph, err := r.createNodeGraph(graph, n, "NaNNode")
	if err != nil {
		return nil, err
	}
	subGraph.SetBackgroundColor("lavenderblush")
	if err := r.renderPath(subGraph, n.GetPath()); err != nil {
		return nil, err
	}
	if _, err := r.renderCommentGroupNode(subGraph, n.GetComment(), yaml.CommentLinePosition); err != nil {
		return nil, err
	}
	if _, err := r.createNode(subGraph, n.Token.Value); err != nil {
		return nil, err
	}
	return subGraph, nil
}

func (r *NodeRenderer) renderMappingKeyNode(graph *graphviz.Graph, n *ast.MappingKeyNode) (*graphviz.Graph, error) {
	subGraph, err := r.createNodeGraph(graph, n, "MappingKeyNode")
	if err != nil {
		return nil, err
	}
	subGraph.SetBackgroundColor("cornsilk")
	if err := r.renderPath(subGraph, n.GetPath()); err != nil {
		return nil, err
	}
	if _, err := r.renderCommentGroupNode(subGraph, n.GetComment(), yaml.CommentHeadPosition); err != nil {
		return nil, err
	}
	value, err := r.createNode(subGraph, "value")
	if err != nil {
		return nil, err
	}
	valueGraph, err := r.renderNode(subGraph, n.Value)
	if err != nil {
		return nil, err
	}
	if err := r.createEdge(subGraph, value, valueGraph); err != nil {
		return nil, err
	}
	return subGraph, nil
}

func (r *NodeRenderer) renderSequenceNode(graph *graphviz.Graph, n *ast.SequenceNode) (*graphviz.Graph, error) {
	subGraph, err := r.createNodeGraph(graph, n, "SequenceNode")
	if err != nil {
		return nil, err
	}
	subGraph.SetBackgroundColor("honeydew")
	if err := r.renderPath(subGraph, n.GetPath()); err != nil {
		return nil, err
	}
	if err := r.renderToken(subGraph, n.Start); err != nil {
		return nil, err
	}
	if err := r.renderToken(subGraph, n.End); err != nil {
		return nil, err
	}
	for _, head := range n.ValueHeadComments {
		if _, err := r.renderCommentGroupNode(subGraph, head, yaml.CommentHeadPosition); err != nil {
			return nil, err
		}
	}
	if _, err := r.renderCommentGroupNode(subGraph, n.FootComment, yaml.CommentFootPosition); err != nil {
		return nil, err
	}
	for idx, value := range n.Values {
		node, err := r.createNode(subGraph, fmt.Sprintf("values[%d]", idx))
		if err != nil {
			return nil, err
		}
		valueGraph, err := r.renderNode(subGraph, value)
		if err != nil {
			return nil, err
		}
		if err := r.createEdge(subGraph, node, valueGraph); err != nil {
			return nil, err
		}
	}
	return subGraph, nil
}

func (r *NodeRenderer) renderAnchorNode(graph *graphviz.Graph, n *ast.AnchorNode) (*graphviz.Graph, error) {
	subGraph, err := r.createNodeGraph(graph, n, "AnchorNode")
	if err != nil {
		return nil, err
	}
	subGraph.SetBackgroundColor("oldlace")
	if err := r.renderPath(subGraph, n.GetPath()); err != nil {
		return nil, err
	}
	if _, err := r.renderCommentGroupNode(subGraph, n.GetComment(), yaml.CommentLinePosition); err != nil {
		return nil, err
	}
	if err := r.renderToken(subGraph, n.Start); err != nil {
		return nil, err
	}
	name, err := r.createNode(subGraph, "name")
	if err != nil {
		return nil, err
	}
	value, err := r.createNode(subGraph, "value")
	if err != nil {
		return nil, err
	}
	nameGraph, err := r.renderNode(subGraph, n.Name)
	if err != nil {
		return nil, err
	}
	valueGraph, err := r.renderNode(subGraph, n.Value)
	if err != nil {
		return nil, err
	}
	if err := r.createEdge(subGraph, name, nameGraph); err != nil {
		return nil, err
	}
	if err := r.createEdge(subGraph, value, valueGraph); err != nil {
		return nil, err
	}
	return subGraph, nil
}

func (r *NodeRenderer) renderAliasNode(graph *graphviz.Graph, n *ast.AliasNode) (*graphviz.Graph, error) {
	subGraph, err := r.createNodeGraph(graph, n, "AliasNode")
	if err != nil {
		return nil, err
	}
	subGraph.SetBackgroundColor("lavenderblush")
	if err := r.renderPath(subGraph, n.GetPath()); err != nil {
		return nil, err
	}
	if _, err := r.renderCommentGroupNode(subGraph, n.GetComment(), yaml.CommentLinePosition); err != nil {
		return nil, err
	}
	if err := r.renderToken(subGraph, n.Start); err != nil {
		return nil, err
	}
	value, err := r.createNode(subGraph, "value")
	if err != nil {
		return nil, err
	}
	valueGraph, err := r.renderNode(subGraph, n.Value)
	if err != nil {
		return nil, err
	}
	if err := r.createEdge(subGraph, value, valueGraph); err != nil {
		return nil, err
	}
	return subGraph, nil
}

func (r *NodeRenderer) renderDirectiveNode(graph *graphviz.Graph, n *ast.DirectiveNode) (*graphviz.Graph, error) {
	subGraph, err := r.createNodeGraph(graph, n, "DirectiveNode")
	if err != nil {
		return nil, err
	}
	subGraph.SetBackgroundColor("snow")
	if err := r.renderToken(subGraph, n.Start); err != nil {
		return nil, err
	}
	if _, err := r.renderCommentGroupNode(subGraph, n.GetComment(), yaml.CommentLinePosition); err != nil {
		return nil, err
	}
	name, err := r.createNode(subGraph, "name")
	if err != nil {
		return nil, err
	}
	nameGraph, err := r.renderNode(subGraph, n.Name)
	if err != nil {
		return nil, err
	}
	if err := r.createEdge(subGraph, name, nameGraph); err != nil {
		return nil, err
	}
	for idx, value := range n.Values {
		val, err := r.createNode(subGraph, fmt.Sprintf("values[%d]", idx))
		if err != nil {
			return nil, err
		}
		valGraph, err := r.renderNode(subGraph, value)
		if err != nil {
			return nil, err
		}
		if err := r.createEdge(subGraph, val, valGraph); err != nil {
			return nil, err
		}
	}
	return subGraph, nil
}

func (r *NodeRenderer) renderTagNode(graph *graphviz.Graph, n *ast.TagNode) (*graphviz.Graph, error) {
	subGraph, err := r.createNodeGraph(graph, n, "TagNode")
	if err != nil {
		return nil, err
	}
	subGraph.SetBackgroundColor("ghostwhite")
	if err := r.renderPath(subGraph, n.GetPath()); err != nil {
		return nil, err
	}
	if _, err := r.renderCommentGroupNode(subGraph, n.GetComment(), yaml.CommentLinePosition); err != nil {
		return nil, err
	}
	if err := r.renderToken(subGraph, n.Start); err != nil {
		return nil, err
	}
	value, err := r.createNode(subGraph, "value")
	if err != nil {
		return nil, err
	}
	valueGraph, err := r.renderNode(subGraph, n.Value)
	if err != nil {
		return nil, err
	}
	if err := r.createEdge(subGraph, value, valueGraph); err != nil {
		return nil, err
	}
	return subGraph, nil
}

func (r *NodeRenderer) renderCommentGroupNode(graph *graphviz.Graph, n *ast.CommentGroupNode, pos yaml.CommentPosition) (*graphviz.Graph, error) {
	if n == nil {
		return nil, nil
	}
	subGraph, err := r.createNodeGraph(graph, n, fmt.Sprintf("CommentGroupNode (%s)", pos))
	if err != nil {
		return nil, err
	}
	subGraph.SetBackgroundColor("whitesmoke")
	for _, cm := range n.Comments {
		if err := r.renderToken(subGraph, cm.Token); err != nil {
			return nil, err
		}
	}
	return subGraph, nil
}

func (r *NodeRenderer) renderPath(graph *graphviz.Graph, p string) error {
	node, err := graph.CreateNodeByName(r.createID())
	if err != nil {
		return err
	}
	node.SetLabel(fmt.Sprintf("{path|%s}", p))
	node.SetShape("record")
	node.SetStyle(graphviz.FilledNodeStyle)
	node.SetFillColor("white")
	return nil
}

func (r *NodeRenderer) renderToken(graph *graphviz.Graph, tk *token.Token) error {
	if tk == nil {
		return nil
	}
	node, err := graph.CreateNodeByName(r.createID())
	if err != nil {
		return err
	}
	pos := tk.Position
	node.SetLabel(fmt.Sprintf("{pos|%d:%d}|{value|%s}", pos.Line, pos.Column, tk.Value))
	node.SetShape("record")
	node.SetStyle(graphviz.FilledNodeStyle)
	node.SetFillColor("white")
	return nil
}

func main() {
	b, err := Parse(context.Background(), `
a: b # comment
`)
	if err != nil {
		panic(err)
	}
	os.WriteFile("file.svg", b, 0o600)
}
