package model

import "gorm.io/datatypes"

type Website struct {
	Base
	Comment     string         `gorm:"comment"           json:"comment"`
	ServerNames datatypes.JSON `gorm:"server_names"      json:"server_names"`
	Ports       datatypes.JSON `gorm:"ports"             json:"ports"`
	Upstreams   datatypes.JSON `gorm:"upstreams"         json:"upstreams"`

	CertFilename string `gorm:"cert_filename"    json:"cert_filename"`
	KeyFilename  string `gorm:"key_filename"     json:"key_filename"`

	IsEnabled bool `gorm:"is_enabled;default=true"       json:"is_enabled"`
}
