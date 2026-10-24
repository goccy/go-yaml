package yaml_test

import (
	"testing"

	yaml "github.com/goccy/go-yaml"
)

// A struct embedding an unexported type tagged `,omitempty` must marshal
// without panicking. isOmittedByOmitEmptyTag used to call v.Interface()
// without the CanInterface guard its omitzero twin has, so Marshal panicked
// with "reflect: cannot return value obtained from unexported field".
type ccOmitInner struct {
	A int
	B string
}

type ccOmitOuter struct {
	ccOmitInner `yaml:",omitempty"`
}

func TestMarshalOmitEmptyEmbeddedUnexported(t *testing.T) {
	if _, err := yaml.Marshal(ccOmitOuter{}); err != nil {
		t.Fatalf("Marshal of embedded unexported + omitempty: %v", err)
	}
}
