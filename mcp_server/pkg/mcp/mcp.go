package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chaitin/SafeLine/mcp_server/pkg/errors"
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mcuadros/go-defaults"
)

type Tool[T any, R any] interface {
	Name() string

	Description() string

	Execute(ctx context.Context, params T) (R, error)

	Validate(params T) error
}

type SSEServer struct {
	sse    *server.SSEServer
	secret string
}

func (s *SSEServer) Start(addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: s,
	}
	return srv.ListenAndServe()
}

func (s *SSEServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.secret == "" {
		s.sse.ServeHTTP(w, r)
		return
	}
	messagePath := s.sse.CompleteMessagePath()
	if messagePath != "" && r.URL.Path == messagePath {
		secret := r.Header.Get("Secret")
		if secret != s.secret {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
	s.sse.ServeHTTP(w, r)
}

type MCPServer struct {
	server *server.MCPServer
	sse    *SSEServer
}

func NewMCPServer(name, version string, secret string) *MCPServer {
	s := server.NewMCPServer(
		name,
		version,
		server.WithLogging(),
	)
	return &MCPServer{
		server: s,
		sse:    &SSEServer{sse: server.NewSSEServer(s), secret: secret},
	}
}

func (s *MCPServer) Start(addr string) error {
	return s.sse.Start(addr)
}

func handleToolCall[T any, R any](ctx context.Context, request mcp.CallToolRequest, tool Tool[T, R]) (result *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()

	var raw []byte
	raw, err = json.Marshal(request.Params.Arguments)
	if err != nil {
		return nil, errors.Wrap(err, "marshal arguments failed")
	}
	var params T
	defaults.SetDefaults(&params)
	if err = json.Unmarshal(raw, &params); err != nil {
		return nil, errors.Wrap(err, "unmarshal parameters failed")
	}

	if err = tool.Validate(params); err != nil {
		return nil, err
	}

	var execResult R
	execResult, err = tool.Execute(ctx, params)
	if err != nil {
		return nil, err
	}

	v := any(execResult)
	switch v := v.(type) {
	case string:
		return mcp.NewToolResultText(v), nil
	case []byte:
		return mcp.NewToolResultText(string(v)), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return mcp.NewToolResultText(json.Number(fmt.Sprint(v)).String()), nil
	case bool:
		return mcp.NewToolResultText(strconv.FormatBool(v)), nil
	default:
		bytes, err := json.Marshal(v)
		if err != nil {
			return nil, errors.New("invalid result type")
		}
		return mcp.NewToolResultText(string(bytes)), nil
	}
}

func RegisterTool[T any, R any](s *MCPServer, tool Tool[T, R]) error {
	var v T
	opts, err := SchemaToOptions(v)
	if err != nil {
		return err
	}
	opts = append(opts, mcp.WithDescription(tool.Description()))
	t := mcp.NewTool(tool.Name(),
		opts...,
	)

	s.server.AddTool(t, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleToolCall(ctx, request, tool)
		if err != nil {
			logger.With("error", err).Error("handle tool call failed")
			if wrapped, ok := err.(*errors.Error); ok {
				return nil, wrapped.Unwrap()
			}
			return nil, err
		}
		return result, nil
	})

	return nil
}
