package ast

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func dumpf(w io.Writer, indentLevel int, typ fmt.Stringer, properties ...string) error {
	indent := strings.Repeat("    ", indentLevel)
	if _, err := fmt.Fprintf(w, "%s- *%s*\n", indent, typ); err != nil {
		return err
	}
	for i := 0; i < len(properties); i += 2 {
		key, value := properties[i], ""
		if i+1 < len(properties) {
			value = properties[i+1]
		}
		value = strconv.Quote(value)
		value = value[1 : len(value)-1]
		if _, err := fmt.Fprintf(w, "%s    - %s: `%s`\n", indent, key, value); err != nil {
			return err
		}
	}
	return nil
}

func dump(w io.Writer, indentLevel int, n interface{}) error {
	if n == nil {
		return nil
	}

	var typ fmt.Stringer
	var properties []string
	if node, ok := n.(Node); ok {
		typ = node.Type()
		if c := node.GetComment(); c != nil {
			properties = append(properties, "Comment", c.Value)
		}
		if t := node.GetToken(); t != nil {
			properties = append(properties, "Token", t.Value)
			properties = append(properties, "Position", t.Position.String())
		}
	}

	var children []interface{}
	switch n := n.(type) {
	case *CommentNode:
	case *NullNode:
	case *IntegerNode:
		properties = append(properties, "Value", fmt.Sprintf("%v", n.Value))
	case *FloatNode:
		properties = append(properties, "Precision", fmt.Sprintf("%v", n.Precision))
		properties = append(properties, "Value", fmt.Sprintf("%v", n.Value))
	case *StringNode:
		properties = append(properties, "Value", fmt.Sprintf("%v", n.Value))
	case *MergeKeyNode:
	case *BoolNode:
		properties = append(properties, "Value", fmt.Sprintf("%v", n.Value))
	case *InfinityNode:
		properties = append(properties, "Value", fmt.Sprintf("%v", n.Value))
	case *NanNode:
	case *LiteralNode:
		properties = append(properties, "Value", fmt.Sprintf("%v", n.Value.Value))
	case *DirectiveNode:
		properties = append(properties, "Start", n.Start.Value)
		children = []interface{}{n.Value}
	case *TagNode:
		properties = append(properties, "Start", n.Start.Value)
		children = []interface{}{n.Value}
	case *DocumentNode:
		if n.Start != nil {
			properties = append(properties, "Start", n.Start.Value)
		}
		if n.End != nil {
			properties = append(properties, "End", n.End.Value)
		}
		children = []interface{}{n.Body}
	case *MappingNode:
		if n.Start != nil {
			properties = append(properties, "Start", n.Start.Value)
		}
		if n.End != nil {
			properties = append(properties, "End", n.End.Value)
		}
		properties = append(properties, "IsFlowStyle", fmt.Sprintf("%v", n.IsFlowStyle))
		for _, value := range n.Values {
			children = append(children, value)
		}
	case *MappingKeyNode:
		if n.Start != nil {
			properties = append(properties, "Start", n.Start.Value)
		}
		children = []interface{}{n.Value}
	case *MappingValueNode:
		if n.Start != nil {
			properties = append(properties, "Start", n.Start.Value)
		}
		children = []interface{}{n.Key, n.Value}
	case *SequenceNode:
		if n.Start != nil {
			properties = append(properties, "Start", n.Start.Value)
		}
		if n.End != nil {
			properties = append(properties, "End", n.End.Value)
		}
		properties = append(properties, "IsFlowStyle", fmt.Sprintf("%v", n.IsFlowStyle))
		for _, v := range n.Values {
			children = append(children, v)
		}
	case *AnchorNode:
		properties = append(properties, "Start", n.Start.Value)
		children = []interface{}{n.Name, n.Value}
	case *AliasNode:
		properties = append(properties, "Start", n.Start.Value)
		children = []interface{}{n.Value}
	}

	if err := dumpf(w, indentLevel, typ, properties...); err != nil {
		return err
	}

	for _, c := range children {
		if err := dump(w, indentLevel+1, c); err != nil {
			return err
		}
	}
	return nil
}

// Dump prints a textual representation of the tree rooted at n to the given writer.
func Dump(w io.Writer, n Node) error {
	return dump(w, 0, n)
}
