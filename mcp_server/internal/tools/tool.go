package tools

import (
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
	"github.com/chaitin/SafeLine/mcp_server/pkg/mcp"
)

// By deferring the concretization of generic types to the Register method,
// we avoid type inference issues.

// Each Tool is wrapped in a toolWrapper that knows its concrete type,
// allowing correct passing of generic parameters during registration.
type ToolWrapper interface {
	Register(s *mcp.MCPServer) error
}

var (
	tools = []ToolWrapper{}
)

func AppendTool[T any, R any](tool ...mcp.Tool[T, R]) {
	for _, t := range tool {
		tools = append(tools, &toolWrapper[T, R]{tool: t})
	}
}

func Tools() []ToolWrapper {
	return tools
}

type toolWrapper[T any, R any] struct {
	tool mcp.Tool[T, R]
}

func (w *toolWrapper[T, R]) Register(s *mcp.MCPServer) error {
	logger.Info("Registering tool",
		logger.String("name", w.tool.Name()),
		logger.String("description", w.tool.Description()),
	)
	return mcp.RegisterTool(s, w.tool)
}
