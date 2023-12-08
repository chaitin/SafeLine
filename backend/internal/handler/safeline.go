package handler

import (
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
