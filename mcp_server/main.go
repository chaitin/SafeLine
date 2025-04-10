package main

import (
	"flag"
	"fmt"

	"github.com/chaitin/SafeLine/mcp_server/internal/api"
	"github.com/chaitin/SafeLine/mcp_server/internal/config"
	"github.com/chaitin/SafeLine/mcp_server/internal/tools"
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
	"github.com/chaitin/SafeLine/mcp_server/pkg/mcp"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	if err := config.Load(*configPath); err != nil {
		panic(fmt.Errorf("failed to load config: %v", err))
	}

	logConfig := config.GetLogger()
	if err := logger.Init(&logger.Config{
		Level:       logConfig.Level,
		FilePath:    logConfig.FilePath,
		Console:     logConfig.Console,
		Caller:      logConfig.Caller,
		Development: logConfig.Development,
	}); err != nil {
		panic(fmt.Errorf("failed to init logger: %v", err))
	}

	logger.With("base_url", config.GetAPI().BaseURL).Info("Initializing API service...")
	if err := api.Init(config.GetAPI()); err != nil {
		panic(fmt.Errorf("failed to init API service: %v", err))
	}

	logger.Info("Starting MCP Server...")
	serverConfig := config.GetServer()
	s := mcp.NewMCPServer(
		serverConfig.Name,
		serverConfig.Version,
		serverConfig.Secret,
	)

	logger.Info("Registering tools...")
	for _, tool := range tools.Tools() {
		if err := tool.Register(s); err != nil {
			logger.With("error", err).
				Error("Failed to register tool")
			panic(err)
		}
	}

	addr := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)
	logger.With("addr", addr).Info("Starting server")
	if err := s.Start(addr); err != nil {
		logger.With("error", err).
			Error("Server failed to start")
		panic(err)
	}
}
