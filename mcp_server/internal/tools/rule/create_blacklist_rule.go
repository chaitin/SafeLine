package rule

import (
	"context"

	"github.com/chaitin/SafeLine/mcp_server/internal/api"
	"github.com/chaitin/SafeLine/mcp_server/internal/api/rule"
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
)

type CreateBlacklistRule struct{}

type CreateBlacklistRuleParams struct {
	Name       string   `json:"name" desc:"name" required:"true"`
	IP         []string `json:"ip" desc:"ip" required:"false"`
	URINoQuery []string `json:"uri_no_query" desc:"uri_no_query" required:"false"`
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
	var pattern [][]api.Pattern
	if len(params.IP) > 0 {
		pattern = append(pattern, []api.Pattern{
			{
				K:    api.KeySrcIP,
				Op:   api.OpEq,
				V:    params.IP,
				SubK: "",
			},
		})
	}
	if len(params.URINoQuery) > 0 {
		pattern = append(pattern, []api.Pattern{
			{
				K:    api.KeyURINoQuery,
				Op:   api.OpEq,
				V:    params.URINoQuery,
				SubK: "",
			},
		})
	}

	id, err := rule.CreateRule(ctx, &rule.CreateRuleRequest{
		Name:      params.Name,
		IP:        params.IP,
		IsEnabled: true,
		Action:    int(api.PolicyRuleActionDeny),
		Pattern:   pattern,
	})
	if err != nil {
		return 0, err
	}
	logger.With("id", id).Info("create blacklist rule success")
	return id, nil
}
