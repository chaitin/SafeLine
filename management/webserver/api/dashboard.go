package api

import (
	"math"
	"time"

	"github.com/gin-gonic/gin"

	"chaitin.cn/patronus/safeline-2/management/webserver/api/response"
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
)

func GetDashboardCounts(c *gin.Context) {
	var requested int64 = 0
	var intercepted int64 = 0
	db := database.GetDB()

	var statisticTotalReq model.SystemStatistics
	r := db.Where("type = 'total-req'").Where("created_at >= date_trunc('day',now())").First(&statisticTotalReq)
	if r.RowsAffected > 0 {
		requested = statisticTotalReq.Value
	}

	var statisticTotalDenied model.SystemStatistics
	r = db.Where("type = 'total-denied'").Where("created_at >= date_trunc('day',now())").First(&statisticTotalDenied)
	if r.RowsAffected > 0 {
		intercepted = statisticTotalDenied.Value
	}
	response.Success(c, gin.H{"requested": requested, "intercepted": intercepted})
}

func GetDashboardSites(c *gin.Context) {
	response.Success(c, gin.H{"normal": 0, "abnormal": 0})
}

func GetDashboardQps(c *gin.Context) {

	var statistics []model.SystemStatistics

	db := database.GetDB()
	db.Where("type = 'req'").Order("created_at desc").Limit(75).Find(&statistics)

	type Node struct {
		Label string `json:"label"`
		Value int64  `json:"value"`
	}
	var nodes = make([]Node, 0)

	for i := len(statistics) - 1; i >= 0; i-- {
		nodes = append(nodes, Node{
			Label: statistics[i].CreatedAt.Format("2006-01-02 15:04:05"),
			Value: int64(math.Ceil(float64(statistics[i].Value) / 5)),
		})
	}

	response.Success(c, gin.H{"nodes": nodes, "total": len(nodes)})
}

func GetDashboardRequests(c *gin.Context) {

	var statistics []model.SystemStatistics

	db := database.GetDB()
	db.Where("type = 'total-req'").Order("created_at desc").Limit(30).Find(&statistics)

	type Node struct {
		Label string `json:"label"`
		Value int64  `json:"value"`
	}
	var nodes = make([]Node, 0)

	for i := 30; i > len(statistics); i-- {
		nodes = append(nodes, Node{
			Label: time.Now().Add(-time.Duration(i-1) * time.Hour * 24).Format("2006-01-02"),
			Value: 0,
		})
	}

	for i := len(statistics) - 1; i >= 0; i-- {
		nodes = append(nodes, Node{
			Label: statistics[i].CreatedAt.Format("2006-01-02"),
			Value: statistics[i].Value,
		})
	}

	response.Success(c, gin.H{"nodes": nodes, "total": len(nodes)})
}

func GetDashboardIntercepts(c *gin.Context) {
	var statistics []model.SystemStatistics

	db := database.GetDB()
	db.Where("type = 'total-denied'").Order("created_at desc").Limit(30).Find(&statistics)

	type Node struct {
		Label string `json:"label"`
		Value int64  `json:"value"`
	}
	var nodes = make([]Node, 0)

	for i := 30; i > len(statistics); i-- {
		nodes = append(nodes, Node{
			Label: time.Now().Add(-time.Duration(i-1) * time.Hour * 24).Format("2006-01-02"),
			Value: 0,
		})
	}

	for i := len(statistics) - 1; i >= 0; i-- {
		nodes = append(nodes, Node{
			Label: statistics[i].CreatedAt.Format("2006-01-02"),
			Value: statistics[i].Value,
		})
	}

	response.Success(c, gin.H{"nodes": nodes, "total": len(nodes)})
}
