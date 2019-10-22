package errors

import (
	"fmt"

	"github.com/goccy/go-yaml/printer"
	"github.com/goccy/go-yaml/token"
	"golang.org/x/xerrors"
)

var (
	ColoredErr = true
)

type WrapError struct {
	err     error
	nextErr error
	frame   xerrors.Frame
	state   fmt.State
	verb    rune
}

func (e *WrapError) chainStateAndVerb(err error) {
	if wrapErr, ok := err.(*WrapError); ok {
		wrapErr.state = e.state
		wrapErr.verb = e.verb
	} else if syntaxErr, ok := err.(*SyntaxError); ok {
		syntaxErr.state = e.state
		syntaxErr.verb = e.verb
	}
}

func (e *WrapError) FormatError(p xerrors.Printer) error {
	if e.verb == 'v' && e.state.Flag('+') && e.state.Flag('#') {
		// print stack trace for debugging
		p.Print(e.err, "\n")
		if e.state.Flag('+') && e.state.Flag('#') {
			e.frame.Format(p)
		}
		e.chainStateAndVerb(e.nextErr)
		return e.nextErr
	}
	err := e.nextErr
	for {
		if wrapErr, ok := err.(*WrapError); ok {
			err = wrapErr.nextErr
			continue
		}
		break
	}
	if ColoredErr && e.verb == 'v' {
		var yp printer.Printer
		p.Print(yp.PrintErrorMessage("syntax error: ", ColoredErr))
	} else {
		p.Print("syntax error: ")
	}
	e.chainStateAndVerb(err)
	err.(*SyntaxError).FormatError(p)
	return nil
}

type WrapState struct {
	org fmt.State
}

func (s *WrapState) Write(b []byte) (n int, err error) {
	return s.org.Write(b)
}

func (s *WrapState) Width() (wid int, ok bool) {
	return s.org.Width()
}

func (s *WrapState) Precision() (prec int, ok bool) {
	return s.org.Precision()
}

func (s *WrapState) Flag(c int) bool {
	if c == '#' {
		// ignore '#' keyword because xerrors.FormatError doesn't set true to printDetail.
		// ( see https://github.com/golang/xerrors/blob/master/adaptor.go#L39-L43 )
		return false
	}
	return s.org.Flag(c)
}

func (e *WrapError) Format(state fmt.State, verb rune) {
	e.state = state
	e.verb = verb
	xerrors.FormatError(e, &WrapState{org: state}, verb)
}

func (e *WrapError) Error() string {
	return e.err.Error()
}

func Wrapf(err error, msg string, args ...interface{}) error {
	return &WrapError{
		err:     xerrors.Errorf(msg, args...),
		nextErr: err,
		frame:   xerrors.Caller(1),
	}
}

type SyntaxError struct {
	fmt.Formatter
	Msg   string
	Token *token.Token
	frame xerrors.Frame
	state fmt.State
	verb  rune
}

func NewSyntaxError(msg string, tk *token.Token) *SyntaxError {
	return &SyntaxError{Msg: msg, Token: tk, frame: xerrors.Caller(1)}
}

func (e *SyntaxError) FormatError(p xerrors.Printer) error {
	if e.verb == 'v' && e.state.Flag('+') && e.state.Flag('#') {
		p.Print(e.Error())
		if p.Detail() {
			e.frame.Format(p)
		}
	} else if e.verb == 'v' {
		p.Print(e.Error())
	} else {
		p.Print(e.Msg)
	}
	return nil
}

func (e *SyntaxError) Error() string {
	var p printer.Printer
	msg := p.PrintErrorMessage(fmt.Sprintf("%s", e.Msg), ColoredErr)
	err := p.PrintErrorToken(e.Token, ColoredErr)
	return fmt.Sprintf("%s\n%s", msg, err)
}
