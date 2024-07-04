package model

// WebsiteConfig is supposed to be same with webserver/model/website.go
type WebsiteConfig struct {
	Id           int      `json:"id"`
	ServerNames  []string `json:"server_names"`
	Ports        []string `json:"ports"`
	Upstreams    []string `json:"upstreams"`
	CertFilename string   `json:"cert_filename"`
	KeyFilename  string   `json:"key_filename"`
}
