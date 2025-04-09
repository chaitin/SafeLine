package tools

import "github.com/chaitin/SafeLine/mcp_server/internal/tools/app"

func init() {
	AppendTool(&CalculateSum{})
	AppendTool(&app.CreateApp{})
}
