package mcp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
	markmcp "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mcuadros/go-defaults"
)

const (
	EndpointPath       = "/mcp"
	maxRequestBodySize = 1 << 20
)

type Tool[T any, R any] interface {
	Name() string

	Description() string

	Execute(ctx context.Context, params T) (R, error)

	Validate(params T) error
}

// InputSchemaProvider lets a tool declare constraints that cannot be inferred
// from Go types alone, such as numeric ranges and documented defaults.
type InputSchemaProvider interface {
	InputSchemaOptions() []markmcp.ToolOption
}

type Middleware func(http.Handler) http.Handler

type MCPServer struct {
	server    *server.MCPServer
	transport *server.StreamableHTTPServer
}

func NewMCPServer(name, version, instructions string) *MCPServer {
	s := server.NewMCPServer(
		name,
		version,
		server.WithInstructions(instructions),
		server.WithToolCapabilities(false),
		server.WithInputSchemaValidation(),
		server.WithStrictInputSchemaDefault(),
		server.WithOutputSchemaValidation(),
	)
	transport := server.NewStreamableHTTPServer(
		s,
		server.WithEndpointPath(EndpointPath),
		server.WithStateLess(true),
	)
	return &MCPServer{server: s, transport: transport}
}

// Handler exposes exactly one Streamable HTTP endpoint. Legacy /sse and
// /message routes are intentionally not mounted.
func (s *MCPServer) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle(EndpointPath, http.MaxBytesHandler(s.transport, maxRequestBodySize))
	return mux
}

func (s *MCPServer) Start(addr string, middleware ...Middleware) error {
	var handler http.Handler = s.Handler()
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	return srv.ListenAndServe()
}

func handleToolCall[T any, R any](ctx context.Context, request markmcp.CallToolRequest, tool Tool[T, R]) (result *markmcp.CallToolResult, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			logger.With("tool", tool.Name()).Error("tool execution panicked")
			result = markmcp.NewToolResultError("tool execution failed")
			err = nil
		}
	}()

	var params T
	defaults.SetDefaults(&params)
	if err := request.BindArguments(&params); err != nil {
		return markmcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}
	if err := tool.Validate(params); err != nil {
		return markmcp.NewToolResultError(fmt.Sprintf("invalid arguments: %v", err)), nil
	}

	execResult, err := tool.Execute(ctx, params)
	if err != nil {
		logger.With("tool", tool.Name()).With("error", err).Error("tool execution failed")
		return markmcp.NewToolResultError(err.Error()), nil
	}

	result, err = markmcp.NewToolResultJSON(execResult)
	if err != nil {
		logger.With("tool", tool.Name()).With("error", err).Error("tool result serialization failed")
		return markmcp.NewToolResultError("tool result serialization failed"), nil
	}
	return result, nil
}

func RegisterTool[T any, R any](s *MCPServer, tool Tool[T, R]) error {
	options := []markmcp.ToolOption{
		markmcp.WithDescription(tool.Description()),
	}
	if provider, ok := any(tool).(InputSchemaProvider); ok {
		options = append(options, provider.InputSchemaOptions()...)
	} else {
		options = append(options, markmcp.WithInputSchema[T]())
	}
	options = append(options,
		markmcp.WithOutputSchema[R](),
		markmcp.WithReadOnlyHintAnnotation(true),
		markmcp.WithDestructiveHintAnnotation(false),
		markmcp.WithIdempotentHintAnnotation(true),
		markmcp.WithOpenWorldHintAnnotation(true),
	)
	t := markmcp.NewTool(tool.Name(), options...)

	s.server.AddTool(t, func(ctx context.Context, request markmcp.CallToolRequest) (*markmcp.CallToolResult, error) {
		return handleToolCall(ctx, request, tool)
	})
	return nil
}
