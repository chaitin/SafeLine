package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chaitin/SafeLine/mcp_server/internal/api"
	"github.com/chaitin/SafeLine/mcp_server/internal/auth"
	"github.com/chaitin/SafeLine/mcp_server/internal/config"
	"github.com/chaitin/SafeLine/mcp_server/internal/tools"
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
	mcpserver "github.com/chaitin/SafeLine/mcp_server/pkg/mcp"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) > 0 && args[0] == "auth" {
		return runAuth(args[1:], stdout, stderr)
	}
	return runServer(args, stderr)
}

func runAuth(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: mcp-server auth <init|rotate> [--state-file path]")
	}

	command := args[0]
	flags := flag.NewFlagSet("auth "+command, flag.ContinueOnError)
	flags.SetOutput(stderr)
	stateFile := flags.String("state-file", auth.StateFileFromEnv(), "authentication state file")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected auth arguments: %s", strings.Join(flags.Args(), " "))
	}

	switch command {
	case "init":
		token, created, err := auth.Initialize(*stateFile)
		if err != nil {
			return err
		}
		if !created {
			_, _ = fmt.Fprintf(stdout, "MCP authentication is already initialized at %s; token was not changed.\n", *stateFile)
			return nil
		}
		printNewToken(stdout, token)
		return nil
	case "rotate":
		token, err := auth.Rotate(*stateFile)
		if err != nil {
			return err
		}
		printNewToken(stdout, token)
		return nil
	default:
		return fmt.Errorf("unknown auth command %q; expected init or rotate", command)
	}
}

func printNewToken(output io.Writer, token string) {
	_, _ = fmt.Fprintln(output, "MCP bearer token (shown once; save it now):")
	_, _ = fmt.Fprintln(output, token)
}

func buildServerInstructions(instances []*config.InstanceConfig) string {
	var instructions strings.Builder
	instructions.WriteString("This server exposes read-only SafeLine data. The get_attack_events tool requires an instance_id.\n")
	instructions.WriteString("Configured SafeLine instance mappings:\n")
	for _, instance := range instances {
		instructions.WriteString(fmt.Sprintf(
			"- display_name %q -> instance_id %q\n",
			instance.DisplayName,
			instance.ID,
		))
	}
	instructions.WriteString("When the user refers to an instance by display_name, use exactly the mapped instance_id. Never guess an instance_id or substitute another instance.")
	return instructions.String()
}

func runServer(args []string, stderr io.Writer) error {
	flags := flag.NewFlagSet("mcp-server", flag.ContinueOnError)
	flags.SetOutput(stderr)
	configPath := flags.String("config", "config.yaml", "path to config file")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(flags.Args(), " "))
	}

	if err := config.Load(*configPath); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logConfig := config.GetLogger()
	if err := logger.Init(&logger.Config{
		Level:       logConfig.Level,
		FilePath:    logConfig.FilePath,
		Console:     logConfig.Console,
		Caller:      logConfig.Caller,
		Development: logConfig.Development,
	}); err != nil {
		return fmt.Errorf("initialize logger: %w", err)
	}

	instances := config.GetInstances()
	logger.With("instances", len(instances)).Info("initializing SafeLine API clients")
	if err := api.InitInstances(instances); err != nil {
		return fmt.Errorf("initialize SafeLine API clients: %w", err)
	}

	authEnabled, err := auth.EnabledFromEnv()
	if err != nil {
		return err
	}
	var middleware []mcpserver.Middleware
	if authEnabled {
		stateFile := auth.StateFileFromEnv()
		authenticator, err := auth.Load(stateFile)
		if err != nil {
			return fmt.Errorf("MCP authentication is enabled but not initialized; run `mcp-server auth init --state-file %s`: %w", stateFile, err)
		}
		middleware = append(middleware, mcpserver.Middleware(authenticator.Middleware))
	}
	logger.With("enabled", authEnabled).Info("MCP bearer authentication configured")

	serverConfig := config.GetServer()
	s := mcpserver.NewMCPServer(serverConfig.Name, serverConfig.Version, buildServerInstructions(instances))
	for _, tool := range tools.Tools() {
		if err := tool.Register(s); err != nil {
			return fmt.Errorf("register tool: %w", err)
		}
	}

	addr := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)
	logger.With("addr", addr).With("endpoint", mcpserver.EndpointPath).Info("starting MCP server")
	return s.Start(addr, middleware...)
}
