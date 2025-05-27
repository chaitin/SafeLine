# SafeLine MCP Server

SafeLine MCP Server is an implementation of the [Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction) that provides complete management and control capabilities for SafeLine WAF.

[![Docker](https://img.shields.io/badge/Docker-Supported-2496ED?style=flat-square&logo=docker&logoColor=white)](docker-compose.yml)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)](go.mod)

## Use Cases

- Automated management and control of SafeLine WAF instances
- WAF configuration and policy management through API
- Building AI-based security protection tools and applications

## Prerequisites

1. Install [Docker](https://www.docker.com/) (if running in container)
2. Configure SafeLine API Token (obtained from SafeLine console)

## Features

- Complete MCP (Management Control Protocol) server implementation
- Support for SafeLine WAF instance management and control
- Flexible configuration system supporting file configuration and environment variables
- Docker containerization support
- Secure API communication

## Quick Start

### Environment Variables

| Environment Variable | Description | Default Value | Required |
|---------|------|--------|-----|
| LISTEN_PORT | Service listening port | 5678 | No |
| LISTEN_ADDRESS | Service listening address | 0.0.0.0 | No |
| SAFELINE_SECRET | SSE server secret | - | No |
| SAFELINE_ADDRESS | SafeLine API address | - | Yes |
| SAFELINE_API_TOKEN | SafeLine API authentication token | - | Yes |

### Using Docker

#### Method 1: Using docker run

```bash
docker run -d \
  --name safeline-mcp \
  -p 5678:5678 \
  -e SAFELINE_API_TOKEN="your_api_token" \
  -e SAFELINE_ADDRESS="https://your.safeline.com" \
  -e LISTEN_PORT=5678 \
  -e LISTEN_ADDRESS="0.0.0.0" \
  chaitin/safeline-mcp:latest
```

#### Method 2: Using docker-compose

```bash
# 1. Clone repository
git clone https://github.com/chaitin/safeline-mcp.git
cd safeline-mcp

# 2. Edit docker-compose.yml to configure environment variables
# Example docker-compose.yml:
# version: '3'
# services:
#   mcp:
#     image: chaitin/safeline-mcp:latest
#     container_name: safeline-mcp
#     ports:
#       - "5678:5678"
#     environment:
#       - SAFELINE_API_TOKEN=your_api_token
#       - SAFELINE_ADDRESS=https://your.safeline.com
#       - LISTEN_PORT=5678
#       - LISTEN_ADDRESS=0.0.0.0

# 3. Start service
docker compose -f docker-compose.yml up -d
```

#### Method 3: Using Go

```bash
# 1. Clone repository
git clone https://github.com/chaitin/SafeLine.git
cd safeline-mcp

# 2. Install dependencies
go mod download

# 3. Configure config.yaml
cp config.yaml.example config.yaml
# Edit config.yaml with necessary configurations

# 4. Run service
go run main.go
```

For more API details, please refer to the [API Documentation](https://demo.waf.chaitin.com:9443/swagger/index.html).

## Tools

### Application Management

- **create_application**

### Rule Management
- **create_blacklist_rule**
- **create_whitelist_rule**

### Analyze
- **get_attack_events**

## Development Guide

The Go API in this project is currently under development, and APIs may change. If you have specific requirements, please submit an Issue for discussion.

### Directory Structure

```
internal/
├── api/              # API implementation
│   ├── app/         # Application-related APIs
│   │   └── create_application.go
│   └── rule/        # Rule-related APIs
│       └── create_rule.go
└── tools/           # MCP tool implementation
    ├── app/         # Application-related tools
    │   └── create_application.go
    └── rule/        # Rule-related tools
        └── create_rule.go
```

### Adding New Tools

1. **Create Tool File**
   - Create corresponding directory and file under `internal/tools`
   - File name should match tool name
   - Use separate file for each tool
   - Example: `internal/tools/app/create_application.go`

2. **Tool Implementation Template**
```go
package app

type ToolName struct{}

type ToolParams struct {
    // Parameter definitions
    Param1 string `json:"param1" desc:"parameter description" required:"true"`
    Param2 int    `json:"param2" desc:"parameter description" required:"false"`
}

type ToolResult struct {
    Field1 string `json:"field1"`
}

func (t *ToolName) Name() string {
    return "tool_name"
}

func (t *ToolName) Description() string {
    return "tool description"
}

func (t *ToolName) Validate(params ToolParams) error {
    // Parameter validation logic
    return nil
}

func (t *ToolName) Execute(ctx context.Context, params ToolParams) (result ToolResult, err error) {
    // Tool execution logic
    return result, nil
}
```

3. **[Optional]Create API Implementation**

If you need to use some APIs that have not been implemented yet, you need to create corresponding files in the api directory for implementation
   - Create same directory structure under `internal/api`
   - File name should match tool func
   - Example: `internal/api/app/create_application.go`

**API Implementation Template**
```go
package app

type RequestType struct {
    // Request parameter definitions
    Param1 string `json:"param1"`
    Param2 int    `json:"param2"`
}

func APIName(ctx context.Context, req *RequestType) (ResultType, error) {
    if req == nil {
        return nil, errors.New("request is required")
    }

    var resp api.Response[ResultType]
    err := api.Service().Post(ctx, "/api/path", req, &resp)
    if err != nil {
        return nil, errors.Wrap(err, "failed to execute")
    }

    if resp.Err != nil {
        return nil, errors.New(resp.Msg)
    }

    return resp.Data, nil
}
```
4. **Tool Registration (init.go)**

The tool registration file `internal/tools/init.go` is used to centrally manage all tool registrations
  - Register all tools uniformly in the `init()` function
  - Use the `AppendTool()` method for registration
  - Example:
    ```go
    // Register create application tool
    AppendTool(&app.CreateApp{})
    
    // Register create blacklist rule tool
    AppendTool(&rule.CreateBlacklistRule{})
    ```

### Development Standards

1. **Naming Conventions**
   - Use lowercase letters and underscores for tool names
   - File names should match tool names

2. **Directory Organization**
   - Divide directories by functional modules (e.g., app, rule, etc.)
   - Maintain consistent structure between tools and api directories
   - Keep related functionality in the same directory

3. **Code Standards**
   - Follow Go standard code conventions
   - Add necessary parameter validation
   - Use unified error handling approach
   - Add appropriate logging

4. **Documentation Requirements**
   - Provide clear functional description in tool Description
   - Add detailed description for parameters
   - Update API toolkit documentation in README

### Example

Refer to the implementation of the `create_application` tool:
- Tool implementation: `internal/tools/app/create_application.go`
- API implementation: `internal/api/app/create_application.go`





