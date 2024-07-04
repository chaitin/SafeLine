package model

import "gorm.io/datatypes"

type PolicyRule struct {
	Base
	Action    int            `gorm:"action"                     json:"action"`
	Comment   string         `gorm:"comment"                    json:"comment"`
	Pattern   datatypes.JSON `gorm:"pattern"                    json:"pattern"`
	IsEnabled bool           `gorm:"is_enabled;default=true"    json:"is_enabled"`
}

type PolicyRulePattern struct {
	K  string `json:"k"`
	Op string `json:"op"`
	V  string `json:"v"`
}

const (
	KeySrcIP = "src_ip"
	KeyURI   = "uri"
	KeyHost  = "host"

	OpEq     = "eq"     // 完全相等
	OpMatch  = "match"  // 模糊匹配
	OpCIDR   = "cidr"   // CIDR
	OpHas    = "has"    // 关键字
	OpPrefix = "prefix" // 前缀关键字
	OpRe     = "re"     // 正则
)
