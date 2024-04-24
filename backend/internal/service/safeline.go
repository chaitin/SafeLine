package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
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

type response[T any] struct {
	Code int `json:"code"`
	Data T   `json:"data"`
}

func (s *SafelineService) request(req *http.Request, data any) error {
	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed, status_code: %d", res.StatusCode)
	}

	if data == nil {
		return nil
	}

	err = json.NewDecoder(res.Body).Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *SafelineService) GetInstallerCount(ctx context.Context) (InstallerCount, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.APIHost+"/api/v1/public/safeline/count", nil)
	if err != nil {
		return cacheCount, err
	}

	var r response[struct {
		Total int `json:"total"`
	}]

	err = s.request(req, &r)
	if r.Code != 0 {
		return cacheCount, nil
	}
	cacheCount = InstallerCount{
		Total: r.Data.Total,
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

	var r response[struct {
		IP string `json:"ip"`
	}]

	err = s.request(req, &r)
	if err != nil {
		return "", err
	}

	if r.Code != 0 {
		return "", nil
	}
	return r.Data.IP, nil
}

type BehaviorType uint64

const (
	BehaviorTypeMin BehaviorType = iota + 1000
	BehaviorTypePurchase
	BehaviorTypeConsult
	BehaviorTypeMax
)

func (s *SafelineService) PostBehavior(ctx context.Context, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.APIHost+"/api/v1/public/safeline/behavior", bytes.NewReader(body))
	if err != nil {
		return err
	}

	err = s.request(req, nil)
	if err != nil {
		return err
	}

	return nil
}
