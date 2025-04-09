package rule

import (
	"context"

	"github.com/chaitin/SafeLine/mcp_server/internal/api"
	"github.com/chaitin/SafeLine/mcp_server/pkg/errors"
)

type CreateRuleRequest struct {
	Name      string          `json:"name"`
	IP        []string        `json:"ip"`
	IsEnabled bool            `json:"is_enabled"`
	Pattern   [][]api.Pattern `json:"pattern"`
	Action    int             `json:"action"`
}

// CreateRule Create new rule
func CreateRule(ctx context.Context, req *CreateRuleRequest) (int64, error) {
	if req == nil {
		return 0, errors.New("request is required")
	}

	var resp api.Response[int64]
	err := api.Service().Post(ctx, "/api/open/policy", req, &resp)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create policy rule")
	}

	if resp.Err != nil {
		return 0, errors.New(resp.Msg)
	}

	return resp.Data, nil
}
