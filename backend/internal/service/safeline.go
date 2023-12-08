package service

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
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

// GetExist return ip if id exist
func (s *SafelineService) GetExist(ctx context.Context, id string, token string) (string, error) {
	body := fmt.Sprintf(`{"id":"%s", "token": "%s"}`, id, token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.APIHost+"/api/v1/public/safeline/exist", strings.NewReader(body))
	if err != nil {
		return "", err
	}
	res, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(res.Body)
		return "", errors.New(string(raw))
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return "", err
	}
	if r["code"].(float64) != 0 {
		return "", nil
	}
	return r["data"].(map[string]interface{})["ip"].(string), nil
}
