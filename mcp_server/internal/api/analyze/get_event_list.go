package analyze

import (
	"context"
	"fmt"

	"github.com/chaitin/SafeLine/mcp_server/internal/api"
)

type GetEventListRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	IP       string `json:"ip"`
	Start    int64  `json:"start"`
	End      int64  `json:"end"`
}

type GetEventListResponse struct {
	Nodes []Event `json:"nodes"`
	Total int64   `json:"total"`
}

type Event struct {
	ID        uint   `json:"id"`
	IP        string `json:"ip"`
	Protocol  int    `json:"protocol"`
	Host      string `json:"host"`
	DstPort   uint64 `json:"dst_port"`
	UpdatedAt int64  `json:"updated_at"`
	StartAt   int64  `json:"start_at"`
	EndAt     int64  `json:"end_at"`
	DenyCount int64  `json:"deny_count"`
	PassCount int64  `json:"pass_count"`
	Finished  bool   `json:"finished"`
	Country   string `json:"country"`
	Province  string `json:"province"`
	City      string `json:"city"`
}

func GetEventList(ctx context.Context, req *GetEventListRequest) (*GetEventListResponse, error) {
	var resp api.Response[GetEventListResponse]
	err := api.Service().Get(ctx, fmt.Sprintf("/api/open/events?page=%d&page_size=%d&ip=%s&start=%d&end=%d", req.Page, req.PageSize, req.IP, req.Start, req.End), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
