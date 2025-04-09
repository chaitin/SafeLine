package rule

import (
	"context"

	"github.com/chaitin/SafeLine/mcp_server/internal/api"
	"github.com/chaitin/SafeLine/mcp_server/internal/api/rule"
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
)

type CreateBlacklistRule struct{}

type CreateBlacklistRuleParams struct {
	Name string   `json:"name" desc:"name" required:"true"`
	IP   []string `json:"ip" desc:"ip" required:"false"`
}

func (t *CreateBlacklistRule) Name() string {
	return "create_blacklist_rule"
}

func (t *CreateBlacklistRule) Description() string {
	return "create a new blacklist rule"
}

func (t *CreateBlacklistRule) Validate(params CreateBlacklistRuleParams) error {
	return nil
}

func (t *CreateBlacklistRule) Execute(ctx context.Context, params CreateBlacklistRuleParams) (int64, error) {
	id, err := rule.CreateBlacklistRule(ctx, &rule.CreateBlacklistRuleRequest{
		Name:      params.Name,
		IP:        params.IP,
		IsEnabled: true,
		Pattern: [][]api.Pattern{
			{
				{
					K:    api.KeySrcIP,
					Op:   api.OpEq,
					V:    params.IP,
					SubK: "",
				},
			},
		},
	})
	if err != nil {
		return 0, err
	}
	logger.With("id", id).Info("create blacklist rule success")
	return id, nil
}
