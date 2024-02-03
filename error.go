package yaml

import (
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/internal/errors"
	"github.com/goccy/go-yaml/token"
	"golang.org/x/xerrors"
)

var (
	ErrInvalidQuery               = xerrors.New("invalid query")
	ErrInvalidPath                = xerrors.New("invalid path instance")
	ErrInvalidPathString          = xerrors.New("invalid path string")
	ErrNotFoundNode               = xerrors.New("node not found")
	ErrUnknownCommentPositionType = xerrors.New("unknown comment position type")
	ErrInvalidCommentMapValue     = xerrors.New("invalid comment map value. it must be not nil value")
)

func ErrUnsupportedHeadPositionType(node ast.Node) error {
	return xerrors.Errorf("unsupported comment head position for %s", node.Type())
}

func ErrUnsupportedLinePositionType(node ast.Node) error {
	return xerrors.Errorf("unsupported comment line position for %s", node.Type())
}

func ErrUnsupportedFootPositionType(node ast.Node) error {
	return xerrors.Errorf("unsupported comment foot position for %s", node.Type())
}

// IsInvalidQueryError whether err is ErrInvalidQuery or not.
func IsInvalidQueryError(err error) bool {
	return xerrors.Is(err, ErrInvalidQuery)
}

// IsInvalidPathError whether err is ErrInvalidPath or not.
func IsInvalidPathError(err error) bool {
	return xerrors.Is(err, ErrInvalidPath)
}

// IsInvalidPathStringError whether err is ErrInvalidPathString or not.
func IsInvalidPathStringError(err error) bool {
	return xerrors.Is(err, ErrInvalidPathString)
}

// IsNotFoundNodeError whether err is ErrNotFoundNode or not.
func IsNotFoundNodeError(err error) bool {
	return xerrors.Is(err, ErrNotFoundNode)
}

// IsInvalidTokenTypeError whether err is ast.ErrInvalidTokenType or not.
func IsInvalidTokenTypeError(err error) bool {
	return xerrors.Is(err, ast.ErrInvalidTokenType)
}

// IsInvalidAnchorNameError whether err is ast.ErrInvalidAnchorName or not.
func IsInvalidAnchorNameError(err error) bool {
	return xerrors.Is(err, ast.ErrInvalidAnchorName)
}

// IsInvalidAliasNameError whether err is ast.ErrInvalidAliasName or not.
func IsInvalidAliasNameError(err error) bool {
	return xerrors.Is(err, ast.ErrInvalidAliasName)
}

// TokenScopedError represents an error associated with a specific [token.Token].
type TokenScopedError struct {
	// Msg is the underlying error message.
	Msg string
	// Token is the [token.Token] associated with this error.
	Token *token.Token
	// err is the underlying, unwraped error.
	err error
}

// Error implements the error interface.
// It returns the unwraped error returned by go-yaml.
func (s TokenScopedError) Error() string {
	return s.err.Error()
}

// AsTokenScopedError checks if the error is associated with a specific token.
// If so, it returns
// Otherwise, it returns nil.
func AsTokenScopedError(err error) *TokenScopedError {
	var syntaxError *errors.SyntaxError
	if xerrors.As(err, &syntaxError) {
		return &TokenScopedError{
			Msg:   syntaxError.GetMessage(),
			Token: syntaxError.GetToken(),
			err:   err,
		}
	}
	var typeError *errors.TypeError
	if xerrors.As(err, &typeError) {
		return &TokenScopedError{
			Msg:   typeError.Error(),
			Token: typeError.Token,
			err:   err,
		}
	}
	return nil
}
