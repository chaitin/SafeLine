package tools

import (
	"github.com/chaitin/SafeLine/mcp_server/internal/tools/app"
	"github.com/chaitin/SafeLine/mcp_server/internal/tools/rule"
)

func init() {
	AppendTool(&app.CreateApp{})
	AppendTool(&rule.CreateBlacklistRule{})
}
