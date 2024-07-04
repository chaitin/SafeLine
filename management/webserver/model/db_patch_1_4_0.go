package model

import (
	"strings"

	"chaitin.cn/dev/go/errors"
	"gorm.io/gorm"
)

type sqlResult struct {
	Ids string `json:"ids"`
}

func DBPatch140(tx *gorm.DB) error {
	if !tx.Migrator().HasTable(&SystemStatistics{}) {
		return nil
	}

	var result []sqlResult
	//SELECT string_agg(id::text, ',') as ids FROM mgt_system_statistics GROUP BY (type, website, created_at) HAVING COUNT(*) > 1) as tmp
	res := tx.Model(&SystemStatistics{}).Select("string_agg(id::text, ',') as ids").Group("type, website, created_at").Having("COUNT(*) > 1").Find(&result)
	if res.Error != nil {
		return errors.Wrap(res.Error, "Failed to select data")
	}

	if len(result) <= 0 {
		return nil
	}

	var deleteIds []string
	for _, s := range result {
		tmpIds := strings.Split(s.Ids, ",")
		deleteIds = append(deleteIds, tmpIds...)
	}
	deleteRes := tx.Delete(&SystemStatistics{}, deleteIds)

	return errors.Wrap(deleteRes.Error, "Failed to delete")
}
