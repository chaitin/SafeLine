package rule

import (
	"context"

	"github.com/chaitin/SafeLine/mcp_server/internal/api"
	"github.com/chaitin/SafeLine/mcp_server/internal/api/rule"
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
)

type CreateWhitelistRule struct{}

type CreateWhitelistRuleParams struct {
	Name       string   `json:"name" desc:"name" required:"true"`
	IP         []string `json:"ip" desc:"ip" required:"false"`
	URINoQuery []string `json:"uri_no_query" desc:"uri_no_query" required:"false"`
}

func (t *CreateWhitelistRule) Name() string {
	return "create_whitelist_rule"
}

func (t *CreateWhitelistRule) Description() string {
	return "create a new whitelist rule"
}

func (t *CreateWhitelistRule) Validate(params CreateWhitelistRuleParams) error {
	return nil
}

func (t *CreateWhitelistRule) Execute(ctx context.Context, params CreateWhitelistRuleParams) (int64, error) {
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
		Action:    int(api.PolicyRuleActionAllow),
		Pattern:   pattern,
	})
	if err != nil {
		return 0, err
	}
	logger.With("id", id).Info("create whitelist rule success")
	return id, nil
}
