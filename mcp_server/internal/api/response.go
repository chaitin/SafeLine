package api

import (
	"fmt"
	"strings"
)

// Response is the envelope returned by SafeLine's management API. SafeLine can
// report an application error with HTTP 200, so callers must not infer success
// from the HTTP status alone.
type Response[T any] struct {
	Data T      `json:"data"`
	Err  any    `json:"err"`
	Msg  string `json:"msg"`
}

// ResponseError represents an error carried in a SafeLine response envelope.
type ResponseError struct {
	Code    any
	Message string
}

func (e *ResponseError) Error() string {
	message := strings.TrimSpace(e.Message)
	if message == "" {
		message = "request failed"
	}
	return fmt.Sprintf("SafeLine API error %v: %s", e.Code, message)
}

// Validate returns an error whenever the envelope contains a non-empty err
// value, including responses that otherwise decoded successfully.
func (r *Response[T]) Validate() error {
	return responseError(r.Err, r.Msg)
}

func responseError(code any, message string) error {
	if code == nil {
		return nil
	}
	if text, ok := code.(string); ok && strings.TrimSpace(text) == "" {
		return nil
	}
	return &ResponseError{Code: code, Message: message}
}
