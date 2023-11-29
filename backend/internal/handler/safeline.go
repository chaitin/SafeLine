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
