package misc

import (
	"fmt"

	"golang.org/x/xerrors"
)

type wrapError struct {
	message string
	next    error
	frame   xerrors.Frame
}

func (e *wrapError) Unwrap() error {
	return e.next
}

func (e *wrapError) Error() string {
	if e.next == nil {
		return e.message
	}
	return fmt.Sprintf("%s: %v", e.message, e.next)
}

func (e *wrapError) Format(f fmt.State, c rune) {
	xerrors.FormatError(e, f, c)
}

func (e *wrapError) FormatError(p xerrors.Printer) error {
	p.Print(e.message)
	if p.Detail() {
		e.frame.Format(p)
	}
	return e.next
}

func wrap(err error, message string, skip int) error {
	if err == nil {
		return nil
	}
	return &wrapError{
		message: message,
		next:    err,
		frame:   xerrors.Caller(skip),
	}
}

// Wrap returns a error annotating `err` with `message` and the caller's frame.
// Wrap returns nil if `err` is nil.
func ErrorWrap(err error, message string) error {
	return wrap(err, message, 2)
}

// Wrapf returns a error annotating `err` with `message` formatted and the caller's frame.
func ErrorWrapf(err error, message string, args ...interface{}) error {
	return wrap(err, fmt.Sprintf(message, args...), 2)
}
