package api

// Response Common API response structure
type Response[T any] struct {
	// Response data
	Data T `json:"data"`
	// Error message
	Err any `json:"err"`
	// Prompt message
	Msg string `json:"msg"`
}
