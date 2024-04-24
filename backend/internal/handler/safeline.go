package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chaitin/SafeLine/internal/service"
)

type SafelineHandler struct {
	safelineService *service.SafelineService
}

func NewSafelineHandler(safelineService *service.SafelineService) *SafelineHandler {
	return &SafelineHandler{
		safelineService: safelineService,
	}
}

// GetInstallerCount
// @Summary get installer count
// @Description get installer count
// @Tags Safeline
// @Accept json
// @Produce json
// @Success 200 {object} service.InstallerCount
// @Router /safeline/count [get]
func (h *SafelineHandler) GetInstallerCount(c *gin.Context) {
	count, err := h.safelineService.GetInstallerCount(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, count)
}

type ExistReq struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}

// Exist return ip if id exist
// @Summary get ip if id exist
// @Description get ip if id exist
// @Tags Safeline
// @Accept json
// @Produce json
// @Param body body ExistReq true "body"
// @Success 200 {object} string
// @Router /exist [post]
func (h *SafelineHandler) Exist(c *gin.Context) {
	req := &ExistReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ip, err := h.safelineService.GetExist(c, req.Id, req.Token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"ip": ip})
}

type BehaviorReq struct {
	Source string               `json:"source"`
	Type   service.BehaviorType `json:"type"`
}

// Behavior record user behavior
// @Summary record user behavior
// @Description record user behavior
// @Tags Safeline
// @Accept json
// @Produce json
// @Param body body BehaviorReq true "body"
// @Success 200 {object} string
// @Router /behavior [post]
func (h *SafelineHandler) Behavior(c *gin.Context) {
	req := &BehaviorReq{}
	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Type >= service.BehaviorTypeMax || req.Type <= service.BehaviorTypeMin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid behavior type"})
		return
	}

	byteReq, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.safelineService.PostBehavior(c, byteReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{})
}
