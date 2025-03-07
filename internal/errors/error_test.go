package errors

import (
	"reflect"
	"testing"

	"github.com/goccy/go-yaml/token"
)

func TestErrorWithMessage(t *testing.T) {
	tok := token.New("foo", "foo", &token.Position{})
	errs := []ErrorWithSource{
		ErrSyntax("original message", tok),
		ErrOverflow(reflect.TypeOf(0), "1", tok),
		ErrTypeMismatch(reflect.TypeOf(0), reflect.TypeOf(""), tok),
		ErrDuplicateKey("original message", tok),
		ErrUnknownField("original message", tok),
		ErrUnexpectedNodeType(0, 1, tok),
	}
	want := "[0:0] new message\n>  0 | foo\n      ^\n"
	for _, err := range errs {
		got := err.WithMessage("new message").Error()
		if got != want {
			t.Fatalf("unexpected error message: %s, wanted: %s", got, want)
		}
	}
}
