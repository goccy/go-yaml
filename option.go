package yaml

import "io"

// DecodeOption functional option type for Decoder
type DecodeOption func(d *Decoder) error

// ReferenceReaders pass to Decoder that reference to anchor defined by passed readers
func ReferenceReaders(readers ...io.Reader) DecodeOption {
	return func(d *Decoder) error {
		d.referenceReaders = append(d.referenceReaders, readers...)
		return nil
	}
}

// ReferenceFiles pass to Decoder that reference to anchor defined by passed files
func ReferenceFiles(files ...string) DecodeOption {
	return func(d *Decoder) error {
		d.referenceFiles = files
		return nil
	}
}

// ReferenceDirs pass to Decoder that reference to anchor defined by files under the passed dirs
func ReferenceDirs(dirs ...string) DecodeOption {
	return func(d *Decoder) error {
		d.referenceDirs = dirs
		return nil
	}
}

// RecursiveDir search yaml file recursively from passed dirs by ReferenceDirs option
func RecursiveDir(isRecursive bool) DecodeOption {
	return func(d *Decoder) error {
		d.isRecursiveDir = isRecursive
		return nil
	}
}

// EncodeOption functional option type for Encoder
type EncodeOption func(e *Encoder) error

// Indent change indent number
func Indent(spaces int) EncodeOption {
	return func(e *Encoder) error {
		e.indent = spaces
		return nil
	}
}
