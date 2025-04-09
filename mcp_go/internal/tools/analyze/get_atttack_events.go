package analyze

import (
	"context"

	"github.com/chaitin/SafeLine/mcp_server/internal/api/analyze"
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
)

type GetAttackEventsParams struct {
	IP       string `json:"ip" desc:"ip" required:"false"`
	Page     int    `json:"page" desc:"page" required:"false" default:"1"`
	PageSize int    `json:"page_size" desc:"page size" required:"false" default:"10"`
	Start    int64  `json:"start" desc:"start unix timestamp in milliseconds" required:"false"`
	End      int64  `json:"end" desc:"end unix timestamp in milliseconds" required:"false"`
}

type GetAttackEvents struct{}

func (t *GetAttackEvents) Name() string {
	return "get_attack_events"
}

func (t *GetAttackEvents) Description() string {
	return "get attack events"
}

func (t *GetAttackEvents) Validate(params GetAttackEventsParams) error {
	return nil
}

func (t *GetAttackEvents) Execute(ctx context.Context, params GetAttackEventsParams) (analyze.GetEventListResponse, error) {
	resp, err := analyze.GetEventList(ctx, &analyze.GetEventListRequest{
		IP:       params.IP,
		PageSize: params.PageSize,
		Page:     params.Page,
		Start:    params.Start,
		End:      params.End,
	})
	if err != nil {
		return analyze.GetEventListResponse{}, err
	}
	logger.With("total", resp.Total).Info("get attack events")
	return *resp, nil
}
