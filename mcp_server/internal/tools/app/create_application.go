package app

import (
	"context"

	"github.com/chaitin/SafeLine/mcp_server/internal/api/app"
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
)

type CreateApp struct{}

type CreateAppParams struct {
	ServerNames []string `json:"server_names" desc:"domain list" required:"true"`
	Ports       []string `json:"ports" desc:"port list" required:"true"`
	Upstreams   []string `json:"upstreams" desc:"upstream list" required:"true"`
}

func (t *CreateApp) Name() string {
	return "create_http_application"
}

func (t *CreateApp) Description() string {
	return "create a new website or app"
}

func (t *CreateApp) Validate(params CreateAppParams) error {
	return nil
}

func (t *CreateApp) Execute(ctx context.Context, params CreateAppParams) (int64, error) {
	id, err := app.CreateApp(ctx, &app.CreateAppRequest{
		ServerNames: params.ServerNames,
		Ports:       params.Ports,
		Upstreams:   params.Upstreams,
	})
	if err != nil {
		return 0, err
	}
	logger.Info("create app success", logger.Int64("id", id))
	return id, nil
}
