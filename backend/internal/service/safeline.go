package service

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
)

var cacheCount InstallerCount

type InstallerCount struct {
	Total int `json:"total"`
}

type SafelineService struct {
	client  *http.Client
	APIHost string
}

func NewSafelineService(host string) *SafelineService {
	return &SafelineService{
		APIHost: host,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (s *SafelineService) GetInstallerCount(ctx context.Context) (InstallerCount, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.APIHost+"/api/v1/public/safeline/count", nil)
	if err != nil {
		return cacheCount, err
	}
	res, err := s.client.Do(req)
	if err != nil {
		return cacheCount, err
	}
	defer res.Body.Close()
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return cacheCount, err
	}
	if r["code"].(float64) != 0 {
		return cacheCount, nil
	}
	cacheCount = InstallerCount{
		Total: int(r["data"].(map[string]interface{})["total"].(float64)),
	}
	return cacheCount, nil
}
