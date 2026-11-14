package controller

import (
	"fmt"
	"strings"
	"testing"
)

func TestGeneratedConfigRejectsBodiesBeyondInspectionLimit(t *testing.T) {
	config := fmt.Sprintf(nginxConfigTpl, "backend_1", "server 127.0.0.1:8080;",
		"listen 0.0.0.0:80;", "server_name example.com;", "", "",
		"http", "backend_1", 1)
	if !strings.Contains(config, "t1k_body_size 1024k;") {
		t.Fatal("inspection limit is missing")
	}
	if !strings.Contains(config, "client_max_body_size 1024k;") {
		t.Fatal("request limit must match the inspection limit")
	}
	if !strings.Contains(config, "proxy_request_buffering on;") {
		t.Fatal("the complete body must be validated before proxying")
	}
}
