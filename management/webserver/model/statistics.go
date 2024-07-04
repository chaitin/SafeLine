package model

import "time"

type SystemStatistics struct {
	ID        uint      `json:"id"                gorm:"primarykey"`
	Type      string    `json:"type"              gorm:"index;uniqueIndex:type_website_createdat"`
	Value     int64     `json:"value"`
	CreatedAt time.Time `json:"created_at"        gorm:"index;uniqueIndex:type_website_createdat"`
	Website   string    `json:"website"           gorm:"index;uniqueIndex:type_website_createdat"`
}
