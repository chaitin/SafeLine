package analyze

import (
	"context"
	"net/url"
	"strconv"

	"github.com/chaitin/SafeLine/mcp_server/internal/api"
)

type GetEventListRequest struct {
	InstanceID string `json:"instance_id"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	IP         string `json:"ip"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Start      int64  `json:"start"`
	End        int64  `json:"end"`
}

type GetEventListResponse struct {
	Source   api.InstanceSource `json:"source"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
	Nodes    []Event            `json:"nodes"`
	Total    int64              `json:"total"`
}

type eventListData struct {
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
	client, err := api.ClientFor(req.InstanceID)
	if err != nil {
		return nil, err
	}

	query := make(url.Values)
	query.Set("page", strconv.Itoa(req.Page))
	query.Set("page_size", strconv.Itoa(req.PageSize))
	if req.IP != "" {
		query.Set("ip", req.IP)
	}
	if req.Host != "" {
		query.Set("host", req.Host)
	}
	if req.Port != 0 {
		query.Set("port", strconv.Itoa(req.Port))
	}
	if req.Start != 0 {
		query.Set("start", strconv.FormatInt(req.Start, 10))
	}
	if req.End != 0 {
		query.Set("end", strconv.FormatInt(req.End, 10))
	}

	var resp api.Response[eventListData]
	if err := client.GetAttackEvents(ctx, query, &resp); err != nil {
		return nil, err
	}
	if err := resp.Validate(); err != nil {
		return nil, err
	}
	return &GetEventListResponse{
		Source:   client.Source(),
		Page:     req.Page,
		PageSize: req.PageSize,
		Nodes:    resp.Data.Nodes,
		Total:    resp.Data.Total,
	}, nil
}
