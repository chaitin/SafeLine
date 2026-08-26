package analyze

import (
	"context"
	"fmt"
	"strings"

	"github.com/chaitin/SafeLine/mcp_server/internal/api/analyze"
	"github.com/chaitin/SafeLine/mcp_server/pkg/errors"
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
	markmcp "github.com/mark3labs/mcp-go/mcp"
)

type GetAttackEventsParams struct {
	InstanceID string `json:"instance_id" jsonschema_description:"Configured SafeLine instance ID"`
	Page       int    `json:"page,omitempty" jsonschema_description:"Page number, starting at 1" default:"1"`
	PageSize   int    `json:"page_size,omitempty" jsonschema_description:"Events per page, from 1 through 100" default:"10"`
	IP         string `json:"ip,omitempty" jsonschema_description:"Source IP address, CIDR, or partial IP filter"`
	Host       string `json:"host,omitempty" jsonschema_description:"Protected host name filter"`
	Port       int    `json:"port,omitempty" jsonschema_description:"Protected destination port filter, from 1 through 65535"`
	Start      int64  `json:"start,omitempty" jsonschema_description:"Start Unix timestamp in milliseconds"`
	End        int64  `json:"end,omitempty" jsonschema_description:"End Unix timestamp in milliseconds"`
}

type GetAttackEvents struct{}

func (t *GetAttackEvents) Name() string {
	return "get_attack_events"
}

func (t *GetAttackEvents) Description() string {
	return "Query aggregated attack events from one configured SafeLine instance"
}

func (t *GetAttackEvents) InputSchemaOptions() []markmcp.ToolOption {
	return []markmcp.ToolOption{
		markmcp.WithString("instance_id",
			markmcp.Required(),
			markmcp.MinLength(1),
			markmcp.Description("Configured SafeLine instance ID"),
		),
		markmcp.WithInteger("page",
			markmcp.Min(1),
			markmcp.DefaultNumber(1),
			markmcp.Description("Page number, starting at 1"),
		),
		markmcp.WithInteger("page_size",
			markmcp.Min(1),
			markmcp.Max(100),
			markmcp.DefaultNumber(10),
			markmcp.Description("Events per page, from 1 through 100"),
		),
		markmcp.WithString("ip", markmcp.Description("Source IP address, CIDR, or partial IP filter")),
		markmcp.WithString("host", markmcp.Description("Protected host name filter")),
		markmcp.WithInteger("port",
			markmcp.Min(1),
			markmcp.Max(65535),
			markmcp.Description("Protected destination port filter, from 1 through 65535"),
		),
		markmcp.WithInteger("start",
			markmcp.Min(0),
			markmcp.Description("Start Unix timestamp in milliseconds"),
		),
		markmcp.WithInteger("end",
			markmcp.Min(0),
			markmcp.Description("End Unix timestamp in milliseconds"),
		),
		markmcp.WithSchemaAdditionalProperties(false),
	}
}

func (t *GetAttackEvents) Validate(params GetAttackEventsParams) error {
	if strings.TrimSpace(params.InstanceID) == "" {
		return errors.New("instance_id is required")
	}
	if params.Page < 1 {
		return errors.New("page must be at least 1")
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		return errors.New("page_size must be between 1 and 100")
	}
	if params.Port < 0 || params.Port > 65535 {
		return errors.New("port must be between 1 and 65535 when provided")
	}
	if params.Start < 0 || params.End < 0 {
		return errors.New("start and end must be non-negative Unix timestamps")
	}
	if params.Start > 0 && params.End > 0 && params.Start > params.End {
		return errors.New("start must be less than or equal to end")
	}
	return nil
}

func (t *GetAttackEvents) Execute(ctx context.Context, params GetAttackEventsParams) (analyze.GetEventListResponse, error) {
	resp, err := analyze.GetEventList(ctx, &analyze.GetEventListRequest{
		InstanceID: strings.TrimSpace(params.InstanceID),
		Page:       params.Page,
		PageSize:   params.PageSize,
		IP:         strings.TrimSpace(params.IP),
		Host:       strings.TrimSpace(params.Host),
		Port:       params.Port,
		Start:      params.Start,
		End:        params.End,
	})
	if err != nil {
		logger.With("instance_id", strings.TrimSpace(params.InstanceID)).With("error", err).Error("SafeLine event query failed")
		return analyze.GetEventListResponse{}, fmt.Errorf("SafeLine query failed for instance %q", strings.TrimSpace(params.InstanceID))
	}
	logger.With("total", resp.Total).Info("get attack events")
	return *resp, nil
}
