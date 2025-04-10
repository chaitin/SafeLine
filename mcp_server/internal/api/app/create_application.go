package app

import (
	"context"

	"github.com/chaitin/SafeLine/mcp_server/internal/api"
	"github.com/chaitin/SafeLine/mcp_server/pkg/errors"
)

type CreateAppRequest struct {
	ServerNames []string `json:"server_names"`
	Ports       []string `json:"ports"`
	Upstreams   []string `json:"upstreams"`
	Comment     string   `json:"comment"`
}

// CreateApp Create new website or app
func CreateApp(ctx context.Context, req *CreateAppRequest) (int64, error) {
	if req == nil {
		return 0, errors.New("request is required")
	}

	var resp api.Response[int64]
	err := api.Service().Post(ctx, "/api/open/site", req, &resp)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create app")
	}

	if resp.Err != nil {
		return 0, errors.New(resp.Msg)
	}

	return resp.Data, nil
}
