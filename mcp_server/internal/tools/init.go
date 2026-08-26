package tools

import (
	"github.com/chaitin/SafeLine/mcp_server/internal/tools/analyze"
)

func init() {
	AppendTool(&analyze.GetAttackEvents{})
}
