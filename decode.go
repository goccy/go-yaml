package yaml

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/internal/errors"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
	"golang.org/x/xerrors"
)

// Decoder reads and decodes YAML values from an input stream.
type Decoder struct {
	reader              io.Reader
	referenceReaders    []io.Reader
	anchorMap           map[string]ast.Node
	opts                []DecodeOption
	referenceFiles      []string
	referenceDirs       []string
	isRecursiveDir      bool
	isResolvedReference bool
	validator           StructValidator
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader, opts ...DecodeOption) *Decoder {
	return &Decoder{
		reader:              r,
		anchorMap:           map[string]ast.Node{},
		opts:                opts,
		referenceReaders:    []io.Reader{},
		referenceFiles:      []string{},
		referenceDirs:       []string{},
		isRecursiveDir:      false,
		isResolvedReference: false,
	}
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
	case *ast.AnchorNode:
		anchorName := n.Name.GetToken().Value
		anchorValue := d.nodeToValue(n.Value)
		d.anchorMap[anchorName] = n.Value
		return anchorValue
	case *ast.AliasNode:
		aliasName := n.Value.GetToken().Value
		return d.nodeToValue(d.anchorMap[aliasName])
	case *ast.LiteralNode:
		return n.Value.GetValue()
	case *ast.FlowMappingNode:
		m := map[string]interface{}{}
		for _, value := range n.Values {
			key := value.Key.GetToken().Value
			m[key] = d.nodeToValue(value.Value)
		}
		return m
	case *ast.MappingValueNode:
		m := map[string]interface{}{}
		if n.Key.Type() == ast.MergeKeyType {
			mapValue := d.nodeToValue(n.Value).(map[string]interface{})
			for k, v := range mapValue {
				m[k] = v
			}
		} else {
			key := n.Key.GetToken().Value
			m[key] = d.nodeToValue(n.Value)
		}
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

func (d *Decoder) getMapNode(node ast.Node) (ast.MapNode, error) {
	if anchor, ok := node.(*ast.AnchorNode); ok {
		mapNode, ok := anchor.Value.(ast.MapNode)
		if ok {
			return mapNode, nil
		}
		return nil, xerrors.Errorf("%s node doesn't MapNode", anchor.Value.Type())
	}
	if alias, ok := node.(*ast.AliasNode); ok {
		aliasName := alias.Value.GetToken().Value
		anchorNode := d.anchorMap[aliasName]
		mapNode, ok := anchorNode.(ast.MapNode)
		if ok {
			return mapNode, nil
		}
		return nil, xerrors.Errorf("%s node doesn't MapNode", anchorNode.Type())
	}
	mapNode, ok := node.(ast.MapNode)
	if !ok {
		return nil, xerrors.Errorf("%s node doesn't MapNode", node.Type())
	}
	return mapNode, nil
}

func (d *Decoder) getArrayNode(node ast.Node) (ast.ArrayNode, error) {
	if anchor, ok := node.(*ast.AnchorNode); ok {
		arrayNode, ok := anchor.Value.(ast.ArrayNode)
		if ok {
			return arrayNode, nil
		}
		return nil, xerrors.Errorf("%s node doesn't ArrayNode", anchor.Value.Type())
	}
	if alias, ok := node.(*ast.AliasNode); ok {
		aliasName := alias.Value.GetToken().Value
		anchorNode := d.anchorMap[aliasName]
		arrayNode, ok := anchorNode.(ast.ArrayNode)
		if ok {
			return arrayNode, nil
		}
		return nil, xerrors.Errorf("%s node doesn't ArrayNode", anchorNode.Type())
	}
	arrayNode, ok := node.(ast.ArrayNode)
	if !ok {
		return nil, xerrors.Errorf("%s node doesn't ArrayNode", node.Type())
	}
	return arrayNode, nil
}

func (d *Decoder) docToNode(doc *ast.Document) ast.Node {
	for _, node := range doc.Nodes {
		if v := d.nodeToValue(node); v != nil {
			return node
		}
	}
	return nil
}

func (d *Decoder) decodeValue(dst reflect.Value, src ast.Node) error {
	valueType := dst.Type()
	if unmarshaler, ok := dst.Addr().Interface().(BytesUnmarshaler); ok {
		b := fmt.Sprintf("%v", src)
		if err := unmarshaler.UnmarshalYAML([]byte(b)); err != nil {
			return errors.Wrapf(err, "failed to UnmarshalYAML")
		}
		return nil
	} else if unmarshaler, ok := dst.Addr().Interface().(InterfaceUnmarshaler); ok {
		if err := unmarshaler.UnmarshalYAML(func(v interface{}) error {
			rv := reflect.ValueOf(v)
			if rv.Type().Kind() != reflect.Ptr {
				return errors.ErrDecodeRequiredPointerType
			}
			if err := d.decodeValue(rv.Elem(), src); err != nil {
				return errors.Wrapf(err, "failed to decode value")
			}
			return nil
		}); err != nil {
			return errors.Wrapf(err, "failed to UnmarshalYAML")
		}
		return nil
	}
	switch valueType.Kind() {
	case reflect.Ptr:
		if dst.IsNil() {
			return nil
		}
		v := d.createDecodableValue(dst.Type())
		if err := d.decodeValue(v, src); err != nil {
			return errors.Wrapf(err, "failed to decode ptr value")
		}
		dst.Set(d.castToAssignableValue(v, dst.Type()))
	case reflect.Interface:
		v := reflect.ValueOf(d.nodeToValue(src))
		if v.IsValid() {
			dst.Set(v)
		}
	case reflect.Map:
		return d.decodeMap(dst, src)
	case reflect.Array, reflect.Slice:
		return d.decodeSlice(dst, src)
	case reflect.Struct:
		return d.decodeStruct(dst, src)
	}
	v := reflect.ValueOf(d.nodeToValue(src))
	if v.IsValid() {
		dst.Set(v.Convert(dst.Type()))
	}
	return nil
}

func (d *Decoder) createDecodableValue(typ reflect.Type) reflect.Value {
	for {
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
			continue
		}
		break
	}
	return reflect.New(typ).Elem()
}

func (d *Decoder) castToAssignableValue(value reflect.Value, target reflect.Type) reflect.Value {
	if target.Kind() != reflect.Ptr {
		return value
	}
	maxTryCount := 5
	tryCount := 0
	for {
		if tryCount > maxTryCount {
			return value
		}
		if value.Type().AssignableTo(target) {
			break
		}
		value = value.Addr()
		tryCount++
	}
	return value
}

func (d *Decoder) keyToNodeMap(node ast.Node) (map[string]ast.Node, error) {
	mapNode, err := d.getMapNode(node)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get map node")
	}
	mapIter := mapNode.MapRange()
	keyToNodeMap := map[string]ast.Node{}
	for mapIter.Next() {
		keyNode := mapIter.Key()
		if keyNode.Type() == ast.MergeKeyType {
			mergeMap, err := d.keyToNodeMap(mapIter.Value())
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get keyToNodeMap by MergeKey node")
			}
			for k, v := range mergeMap {
				keyToNodeMap[k] = v
			}
		} else {
			key, ok := d.nodeToValue(keyNode).(string)
			if !ok {
				return nil, errors.Wrapf(err, "failed to decode map key")
			}
			keyToNodeMap[key] = mapIter.Value()
		}
	}
	return keyToNodeMap, nil
}

func (d *Decoder) decodeStruct(dst reflect.Value, src ast.Node) error {
	if src == nil {
		return nil
	}
	structType := dst.Type()
	structValue := reflect.New(structType)
	structFieldMap, err := structFieldMap(structType)
	if err != nil {
		return errors.Wrapf(err, "failed to create struct field map")
	}
	keyToNodeMap, err := d.keyToNodeMap(src)
	if err != nil {
		return errors.Wrapf(err, "failed to get keyToNodeMap")
	}
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if isIgnoredStructField(field) {
			continue
		}
		structField := structFieldMap[field.Name]
		if structField.IsInline {
			fieldValue := structValue.Elem().FieldByName(field.Name)
			if !fieldValue.CanSet() {
				return xerrors.Errorf("cannot set embedded type as unexported field %s.%s", field.PkgPath, field.Name)
			}
			newFieldValue := d.createDecodableValue(fieldValue.Type())
			if err := d.decodeValue(newFieldValue, src); err != nil {
				return errors.Wrapf(err, "failed to decode value")
			}
			fieldValue.Set(d.castToAssignableValue(newFieldValue, fieldValue.Type()))
			continue
		}
		v, exists := keyToNodeMap[structField.RenderName]
		if !exists {
			continue
		}
		fieldValue := structValue.Elem().FieldByName(field.Name)
		newFieldValue := d.createDecodableValue(fieldValue.Type())
		if err := d.decodeValue(newFieldValue, v); err != nil {
			return errors.Wrapf(err, "failed to decode value")
		}
		fieldValue.Set(d.castToAssignableValue(newFieldValue, fieldValue.Type()))
	}
	if d.validator != nil {
		if err := d.validator.Struct(structValue.Interface()); err != nil {
			ev := reflect.ValueOf(err)
			if ev.Type().Kind() == reflect.Slice {
				for i := 0; i < ev.Len(); i++ {
					fieldErr, ok := ev.Index(i).Interface().(FieldError)
					if !ok {
						continue
					}
					fieldName := fieldErr.StructField()
					structField := structFieldMap[fieldName]
					node, exists := keyToNodeMap[structField.RenderName]
					if exists {
						// TODO: to make FieldError message cutomizable
						return errors.ErrSyntax(fmt.Sprintf("%s", err), node.GetToken())
					}
				}
			}
		}
	}
	dst.Set(structValue.Elem())
	return nil
}

func (d *Decoder) decodeSlice(dst reflect.Value, src ast.Node) error {
	arrayNode, err := d.getArrayNode(src)
	if err != nil {
		return errors.Wrapf(err, "failed to get array node")
	}
	iter := arrayNode.ArrayRange()
	sliceType := dst.Type()
	sliceValue := reflect.MakeSlice(sliceType, 0, iter.Len())
	elemType := sliceType.Elem()
	for iter.Next() {
		v := iter.Value()
		dstValue := d.createDecodableValue(elemType)
		if err := d.decodeValue(dstValue, v); err != nil {
			return errors.Wrapf(err, "failed to decode value")
		}
		sliceValue = reflect.Append(sliceValue, d.castToAssignableValue(dstValue, elemType))
	}
	dst.Set(sliceValue)
	return nil
}

func (d *Decoder) decodeMap(dst reflect.Value, src ast.Node) error {
	mapNode, err := d.getMapNode(src)
	if err != nil {
		return errors.Wrapf(err, "failed to get map node")
	}
	mapType := dst.Type()
	mapValue := reflect.MakeMap(mapType)
	keyType := mapValue.Type().Key()
	valueType := mapValue.Type().Elem()
	mapIter := mapNode.MapRange()
	for mapIter.Next() {
		key := mapIter.Key()
		value := mapIter.Value()
		dstValue := d.createDecodableValue(valueType)
		if err := d.decodeValue(dstValue, value); err != nil {
			return errors.Wrapf(err, "failed to decode value")
		}
		castedKey := reflect.ValueOf(d.nodeToValue(key)).Convert(keyType)
		mapValue.SetMapIndex(castedKey, d.castToAssignableValue(dstValue, valueType))
	}
	dst.Set(mapValue)
	return nil
}

func (d *Decoder) fileToReader(file string) (io.Reader, error) {
	reader, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file")
	}
	return reader, nil
}

func (d *Decoder) isYAMLFile(file string) bool {
	ext := filepath.Ext(file)
	if ext == ".yml" {
		return true
	}
	if ext == ".yaml" {
		return true
	}
	return false
}

func (d *Decoder) readersUnderDir(dir string) ([]io.Reader, error) {
	pattern := fmt.Sprintf("%s/*", dir)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get files by %s", pattern)
	}
	readers := []io.Reader{}
	for _, match := range matches {
		if !d.isYAMLFile(match) {
			continue
		}
		reader, err := d.fileToReader(match)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get reader")
		}
		readers = append(readers, reader)
	}
	return readers, nil
}

func (d *Decoder) readersUnderDirRecursive(dir string) ([]io.Reader, error) {
	readers := []io.Reader{}
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !d.isYAMLFile(path) {
			return nil
		}
		reader, err := d.fileToReader(path)
		if err != nil {
			return errors.Wrapf(err, "failed to get reader")
		}
		readers = append(readers, reader)
		return nil
	}); err != nil {
		return nil, errors.Wrapf(err, "interrupt walk in %s", dir)
	}
	return readers, nil
}

func (d *Decoder) resolveReference() error {
	for _, opt := range d.opts {
		if err := opt(d); err != nil {
			return errors.Wrapf(err, "failed to exec option")
		}
	}
	for _, file := range d.referenceFiles {
		reader, err := d.fileToReader(file)
		if err != nil {
			return errors.Wrapf(err, "failed to get reader")
		}
		d.referenceReaders = append(d.referenceReaders, reader)
	}
	for _, dir := range d.referenceDirs {
		if !d.isRecursiveDir {
			readers, err := d.readersUnderDir(dir)
			if err != nil {
				return errors.Wrapf(err, "failed to get readers from under the %s", dir)
			}
			d.referenceReaders = append(d.referenceReaders, readers...)
		} else {
			readers, err := d.readersUnderDirRecursive(dir)
			if err != nil {
				return errors.Wrapf(err, "failed to get readers from under the %s", dir)
			}
			d.referenceReaders = append(d.referenceReaders, readers...)
		}
	}
	for _, reader := range d.referenceReaders {
		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			return errors.Wrapf(err, "failed to read buffer")
		}

		// assign new anchor definition to anchorMap
		if _, err := d.decode(bytes); err != nil {
			return errors.Wrapf(err, "failed to decode")
		}
	}
	d.isResolvedReference = true
	return nil
}

func (d *Decoder) decode(bytes []byte) (ast.Node, error) {
	var (
		parser parser.Parser
	)
	tokens := lexer.Tokenize(string(bytes))
	doc, err := parser.Parse(tokens)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse yaml")
	}
	return d.docToNode(doc), nil
}

// Decode reads the next YAML-encoded value from its input
// and stores it in the value pointed to by v.
//
// See the documentation for Unmarshal for details about the
// conversion of YAML into a Go value.
func (d *Decoder) Decode(v interface{}) error {
	if !d.isResolvedReference {
		if err := d.resolveReference(); err != nil {
			return errors.Wrapf(err, "failed to resolve reference")
		}
	}
	rv := reflect.ValueOf(v)
	if rv.Type().Kind() != reflect.Ptr {
		return errors.ErrDecodeRequiredPointerType
	}
	bytes, err := ioutil.ReadAll(d.reader)
	if err != nil {
		return errors.Wrapf(err, "failed to read buffer")
	}
	node, err := d.decode(bytes)
	if err != nil {
		return errors.Wrapf(err, "failed to decode")
	}
	if node == nil {
		return nil
	}
	if err := d.decodeValue(rv.Elem(), node); err != nil {
		return errors.Wrapf(err, "failed to decode value")
	}
	return nil
}
