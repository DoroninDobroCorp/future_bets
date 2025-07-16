package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetLeagueCandidates(c *gin.Context) {
	candidates := h.autoMatchService.GetLeagueCandidates()

	if len(candidates) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": candidates,
	})
}

func (h *Handler) GetTeamCandidates(c *gin.Context) {
	candidates := h.autoMatchService.GetTeamCandidates()

	if len(candidates) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": candidates,
	})
}
