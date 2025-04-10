package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
)

var (
	// Common errors
	ErrInternal     = New("internal error")
	ErrInvalidParam = New("invalid parameter")
	ErrNotFound     = New("resource not found")
	ErrUnauthorized = New("unauthorized")
	ErrForbidden    = New("forbidden")
	ErrTimeout      = New("timeout")
)

// Error Custom error structure
type Error struct {
	err      error
	stack    []string
	msg      string
	location string
}

// Error Implement error interface
func (e *Error) Error() string {
	if e.msg != "" {
		return fmt.Sprintf("%s: %v (at %s)", e.msg, e.err, e.location)
	}
	return fmt.Sprintf("%v (at %s)", e.err, e.location)
}

// Unwrap Return original error
func (e *Error) Unwrap() error {
	if e.err == nil {
		return nil
	}
	if wrapped, ok := e.err.(*Error); ok {
		return wrapped.Unwrap()
	}
	return e.err
}

// Stack Return error stack
func (e *Error) Stack() []string {
	return e.stack
}

// Location Return error location
func (e *Error) Location() string {
	return e.location
}

// getCallerLocation Get caller location
func getCallerLocation(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// WrapL Wrap error and print log
func WrapL(err error, msg string) error {
	if err == nil {
		return nil
	}

	// Get stack trace information
	var stack []string
	for i := 1; i < 32; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			break
		}
		name := fn.Name()
		if strings.Contains(name, "runtime.") {
			break
		}
		stack = append(stack, fmt.Sprintf("%s:%d", file, line))
	}

	wrappedErr := &Error{
		err:      err,
		stack:    stack,
		msg:      msg,
		location: getCallerLocation(2),
	}

	// Print error information and stack using logger
	logger.With("error", err).
		With("location", wrappedErr.location).
		With("stack", strings.Join(stack, "\n")).
		Error(msg)

	return wrappedErr
}

// Is Check error type
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As Type assertion
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Wrap Wrap error without printing log
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &Error{
		err:      err,
		msg:      msg,
		location: getCallerLocation(2),
	}
}

// New Create new error
func New(text string) error {
	return &Error{
		err:      errors.New(text),
		location: getCallerLocation(2),
	}
}
