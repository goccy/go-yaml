package yaml

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
	"golang.org/x/xerrors"
)

type Decoder struct {
	reader io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{reader: r}
}

func (d *Decoder) nodeToValue(node ast.Node) interface{} {
	switch n := node.(type) {
	case *ast.NullNode:
		return nil
	case *ast.StringNode:
		return n.GetValue()
	case *ast.IntegerNode:
		return n.GetValue()
	case *ast.FloatNode:
		return n.GetValue()
	case *ast.BoolNode:
		return n.GetValue()
	case *ast.InfinityNode:
		return n.GetValue()
	case *ast.NanNode:
		return n.GetValue()
	case *ast.TagNode:
		switch n.Start.Value {
		case token.BinaryTag:
			b, _ := base64.StdEncoding.DecodeString(d.nodeToValue(n.Value).(string))
			return b
		}
	case *ast.LiteralNode:
		return n.Value.GetValue()
	case *ast.FlowMappingNode:
		m := map[string]interface{}{}
		for _, value := range n.Values {
			key := value.Key.GetToken().Value
			m[key] = d.nodeToValue(value.Value)
		}
		return m
	case *ast.MappingNode:
		m := map[string]interface{}{}
		key := n.Key.GetToken().Value
		subMap := map[string]interface{}{}
		for _, value := range n.Values {
			for k, v := range d.nodeToValue(value).(map[string]interface{}) {
				subMap[k] = v
			}
		}
		m[key] = subMap
		return m
	case *ast.MappingValueNode:
		m := map[string]interface{}{}
		key := n.Key.GetToken().Value
		m[key] = d.nodeToValue(n.Value)
		return m
	case *ast.MappingCollectionNode:
		m := map[string]interface{}{}
		for _, value := range n.Values {
			subMap := d.nodeToValue(value).(map[string]interface{})
			for k, v := range subMap {
				m[k] = v
			}
		}
		return m
	case *ast.FlowSequenceNode:
		v := []interface{}{}
		for _, value := range n.Values {
			v = append(v, d.nodeToValue(value))
		}
		return v
	case *ast.SequenceNode:
		v := []interface{}{}
		for _, value := range n.Values {
			v = append(v, d.nodeToValue(value))
		}
		return v
	}
	return nil
}

func (d *Decoder) docToValue(doc *ast.Document) interface{} {
	for _, node := range doc.Nodes {
		return d.nodeToValue(node)
	}
	return nil
}

func (d *Decoder) decodeValue(valueType reflect.Type, value interface{}) reflect.Value {
	switch valueType.Kind() {
	case reflect.Ptr:
		return reflect.ValueOf(nil)
	case reflect.Interface:
		return reflect.ValueOf(value)
	case reflect.Map:
		return d.decodeMap(valueType, value)
	case reflect.Array, reflect.Slice:
		return d.decodeSlice(valueType, value)
	case reflect.Struct:
		return reflect.ValueOf(nil)
	}
	return reflect.ValueOf(value).Convert(valueType)
}

func (d *Decoder) decodeSlice(sliceType reflect.Type, value interface{}) reflect.Value {
	slice := value.([]interface{})
	sliceValue := reflect.MakeSlice(sliceType, 0, len(slice))
	sliceValueType := sliceValue.Type().Elem()
	for _, v := range slice {
		sliceValue = reflect.Append(sliceValue, d.decodeValue(sliceValueType, v))
	}
	return sliceValue
}

func (d *Decoder) decodeMap(mapType reflect.Type, value interface{}) reflect.Value {
	mapValue := reflect.MakeMap(mapType)
	keyType := mapValue.Type().Key()
	valueType := mapValue.Type().Elem()
	for k, v := range value.(map[string]interface{}) {
		castedKey := reflect.ValueOf(k).Convert(keyType)
		mapValue.SetMapIndex(castedKey, d.decodeValue(valueType, v))
	}
	return mapValue
}

func (d *Decoder) Decode(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Type().Kind() != reflect.Ptr {
		return xerrors.New("required pointer type value")
	}
	bytes, err := ioutil.ReadAll(d.reader)
	if err != nil {
		return xerrors.Errorf("failed to read buffer: %w", err)
	}
	var (
		lex    lexer.Lexer
		parser parser.Parser
	)
	tokens := lex.Tokenize(string(bytes))
	doc, err := parser.Parse(tokens)
	if err != nil {
		return xerrors.Errorf("failed to parse yaml: %w", err)
	}
	value := d.docToValue(doc)
	if value == nil {
		return nil
	}
	decodedValue := d.decodeValue(rv.Elem().Type(), value)
	if decodedValue.IsValid() {
		rv.Elem().Set(decodedValue)
	}
	return nil
}
