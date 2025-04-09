package tools

import (
	"context"

	"github.com/chaitin/SafeLine/mcp_server/pkg/errors"
	"github.com/chaitin/SafeLine/mcp_server/pkg/logger"
)

type CalculateSum struct{}

func (t *CalculateSum) Name() string {
	return "calculate_sum"
}

func (t *CalculateSum) Description() string {
	return "Add two numbers together"
}

type MyToolInput struct {
	A int `json:"a" desc:"number a" required:"true"`
	B int `json:"b" desc:"number b" required:"true"`
}

type MyToolOutput struct {
	C int `json:"c"`
}

func (t *CalculateSum) Validate(params MyToolInput) error {
	return nil
}

func (t *CalculateSum) Execute(ctx context.Context, params MyToolInput) (MyToolOutput, error) {
	logger.With("a", params.A).
		With("b", params.B).
		Debug("Executing calculation")

	result := MyToolOutput{
		C: params.A + params.B,
	}

	logger.With("result", result.C).
		Debug("Calculation completed")

	return result, errors.New("test error")
}
