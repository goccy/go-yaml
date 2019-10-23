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
	"github.com/goccy/go-yaml/errors"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
)

// Decoder reads and decodes YAML values from an input stream.
type Decoder struct {
	reader              io.Reader
	referenceReaders    []io.Reader
	anchorMap           map[string]interface{}
	opts                []DecodeOption
	referenceFiles      []string
	referenceDirs       []string
	isRecursiveDir      bool
	isResolvedReference bool
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader, opts ...DecodeOption) *Decoder {
	return &Decoder{
		reader:              r,
		anchorMap:           map[string]interface{}{},
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
		d.anchorMap[anchorName] = anchorValue
		return anchorValue
	case *ast.AliasNode:
		aliasName := n.Value.GetToken().Value
		return d.anchorMap[aliasName]
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

func (d *Decoder) docToValue(doc *ast.Document) interface{} {
	for _, node := range doc.Nodes {
		if v := d.nodeToValue(node); v != nil {
			return v
		}
	}
	return nil
}

func (d *Decoder) decodeValue(valueType reflect.Type, value interface{}) (reflect.Value, error) {
	switch valueType.Kind() {
	case reflect.Ptr:
		v, err := d.decodeValue(valueType.Elem(), value)
		if err != nil {
			return reflect.Zero(valueType), errors.Wrapf(err, "failed to decode ptr value")
		}
		return v.Addr(), nil
	case reflect.Interface:
		return reflect.ValueOf(value), nil
	case reflect.Map:
		return d.decodeMap(valueType, value)
	case reflect.Array, reflect.Slice:
		return d.decodeSlice(valueType, value)
	case reflect.Struct:
		return d.decodeStruct(valueType, value)
	}
	return reflect.ValueOf(value).Convert(valueType), nil
}

func (d *Decoder) decodeStruct(structType reflect.Type, value interface{}) (reflect.Value, error) {
	structValue := reflect.New(structType)
	structFieldMap, err := structFieldMap(structType)
	if err != nil {
		return reflect.Zero(structType), errors.Wrapf(err, "failed to create struct field map")
	}
	if value == nil {
		return reflect.Zero(structType), nil
	}
	valueMap, ok := value.(map[string]interface{})
	if !ok {
		return reflect.Zero(structType), errors.Wrapf(err, "value is not struct type: %s", reflect.TypeOf(value).Name())
	}
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if isIgnoredStructField(field) {
			continue
		}
		structField := structFieldMap[field.Name]
		v, exists := valueMap[structField.RenderName]
		if !exists {
			continue
		}
		fieldValue := structValue.Elem().FieldByName(field.Name)
		valueType := fieldValue.Type()
		vv, err := d.decodeValue(valueType, v)
		if err != nil {
			return reflect.Zero(structType), errors.Wrapf(err, "failed to decode value")
		}
		fieldValue.Set(vv)
	}
	return structValue.Elem(), nil
}

func (d *Decoder) decodeSlice(sliceType reflect.Type, value interface{}) (reflect.Value, error) {
	slice := value.([]interface{})
	sliceValue := reflect.MakeSlice(sliceType, 0, len(slice))
	sliceValueType := sliceValue.Type().Elem()
	for _, v := range slice {
		vv, err := d.decodeValue(sliceValueType, v)
		if err != nil {
			return reflect.Zero(sliceType), errors.Wrapf(err, "failed to decode value")
		}
		sliceValue = reflect.Append(sliceValue, vv)
	}
	return sliceValue, nil
}

func (d *Decoder) decodeMap(mapType reflect.Type, value interface{}) (reflect.Value, error) {
	mapValue := reflect.MakeMap(mapType)
	keyType := mapValue.Type().Key()
	valueType := mapValue.Type().Elem()
	for k, v := range value.(map[string]interface{}) {
		castedKey := reflect.ValueOf(k).Convert(keyType)
		vv, err := d.decodeValue(valueType, v)
		if err != nil {
			return reflect.Zero(mapType), errors.Wrapf(err, "failed to decode value")
		}
		mapValue.SetMapIndex(castedKey, vv)
	}
	return mapValue, nil
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

func (d *Decoder) decode(bytes []byte) (interface{}, error) {
	var (
		parser parser.Parser
	)
	tokens := lexer.Tokenize(string(bytes))
	doc, err := parser.Parse(tokens)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse yaml")
	}
	return d.docToValue(doc), nil
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
	value, err := d.decode(bytes)
	if err != nil {
		return errors.Wrapf(err, "failed to decode")
	}
	if value == nil {
		return nil
	}
	decodedValue, err := d.decodeValue(rv.Elem().Type(), value)
	if err != nil {
		return errors.Wrapf(err, "failed to decode value")
	}
	if decodedValue.IsValid() {
		rv.Elem().Set(decodedValue)
	}
	return nil
}
