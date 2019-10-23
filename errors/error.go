package errors

import (
	"fmt"

	"github.com/goccy/go-yaml/printer"
	"github.com/goccy/go-yaml/token"
	"golang.org/x/xerrors"
)

var (
	// ColoredErr error with syntax highlight
	ColoredErr = true
	// WithSourceCode error with source code
	WithSourceCode = true
	// ErrDecodeRequiredPointerType error instance for decoding
	ErrDecodeRequiredPointerType = xerrors.New("required pointer type value")
)

// Wrapf wrap error for stack trace
func Wrapf(err error, msg string, args ...interface{}) error {
	return &wrapError{
		baseError: &baseError{},
		err:       xerrors.Errorf(msg, args...),
		nextErr:   err,
		frame:     xerrors.Caller(1),
	}
}

// ErrSyntax create syntax error instance with message and token
func ErrSyntax(msg string, tk *token.Token) *syntaxError {
	return &syntaxError{
		baseError: &baseError{},
		msg:       msg,
		token:     tk,
		frame:     xerrors.Caller(1),
	}
}

type baseError struct {
	state fmt.State
	verb  rune
}

func (e *baseError) Error() string {
	return ""
}

func (e *baseError) chainStateAndVerb(err error) {
	wrapErr, ok := err.(*wrapError)
	if ok {
		wrapErr.state = e.state
		wrapErr.verb = e.verb
	}
	syntaxErr, ok := err.(*syntaxError)
	if ok {
		syntaxErr.state = e.state
		syntaxErr.verb = e.verb
	}
}

type wrapError struct {
	*baseError
	err     error
	nextErr error
	frame   xerrors.Frame
}

func (e *wrapError) FormatError(p xerrors.Printer) error {
	if e.verb == 'v' && e.state.Flag('+') {
		// print stack trace for debugging
		p.Print(e.err, "\n")
		e.frame.Format(p)
		e.chainStateAndVerb(e.nextErr)
		return e.nextErr
	}
	err := e.nextErr
	for {
		if wrapErr, ok := err.(*wrapError); ok {
			err = wrapErr.nextErr
			continue
		}
		break
	}
	e.chainStateAndVerb(err)
	if fmtErr, ok := err.(xerrors.Formatter); ok {
		fmtErr.FormatError(p)
	} else {
		p.Print(err)
	}
	return nil
}

type wrapState struct {
	org fmt.State
}

func (s *wrapState) Write(b []byte) (n int, err error) {
	return s.org.Write(b)
}

func (s *wrapState) Width() (wid int, ok bool) {
	return s.org.Width()
}

func (s *wrapState) Precision() (prec int, ok bool) {
	return s.org.Precision()
}

func (s *wrapState) Flag(c int) bool {
	// set true to 'printDetail' forced because when p.Detail() is false, xerrors.Printer no output any text
	if c == '#' {
		// ignore '#' keyword because xerrors.FormatError doesn't set true to printDetail.
		// ( see https://github.com/golang/xerrors/blob/master/adaptor.go#L39-L43 )
		return false
	}
	return true
}

func (e *wrapError) Format(state fmt.State, verb rune) {
	e.state = state
	e.verb = verb
	xerrors.FormatError(e, &wrapState{org: state}, verb)
}

func (e *wrapError) Error() string {
	return e.err.Error()
}

type syntaxError struct {
	*baseError
	msg   string
	token *token.Token
	frame xerrors.Frame
}

func (e *syntaxError) FormatError(p xerrors.Printer) error {
	if e.verb == 'v' && e.state.Flag('+') {
		// %+v
		// print stack trace for debugging
		p.Print(e.Error())
		e.frame.Format(p)
	} else {
		p.Print(e.Error())
	}
	return nil
}

func (e *syntaxError) Error() string {
	var p printer.Printer
	pos := fmt.Sprintf("[%d:%d] ", e.token.Position.Line, e.token.Position.Column)
	msg := p.PrintErrorMessage(fmt.Sprintf("syntax error: %s%s", pos, e.msg), ColoredErr)
	if WithSourceCode {
		err := p.PrintErrorToken(e.token, ColoredErr)
		return fmt.Sprintf("%s\n%s", msg, err)
	}
	return msg
}
