package tools

import (
	"github.com/chaitin/SafeLine/mcp_server/internal/tools/analyze"
	"github.com/chaitin/SafeLine/mcp_server/internal/tools/app"
	"github.com/chaitin/SafeLine/mcp_server/internal/tools/rule"
)

func init() {
	// app
	AppendTool(&app.CreateApp{})

	// rule
	AppendTool(&rule.CreateBlacklistRule{})
	AppendTool(&rule.CreateWhitelistRule{})

	// analyze
	AppendTool(&analyze.GetAttackEvents{})
}
