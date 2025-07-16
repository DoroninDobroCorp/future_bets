package handler

import (
	"livebets/runner/internal/entity"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SetCommand(c *gin.Context) {
	bookmaker := c.Query("bookmaker")
	if bookmaker == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	runRaw := c.Query("run")
	if runRaw == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	run, err := strconv.ParseBool(runRaw)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	h.commandService.SetCommand(c, entity.Command{Name: bookmaker, Run: run})
	c.Status(http.StatusOK)
}

func (h *Handler) GetStatuses(c *gin.Context) {
	statuses := h.statusService.GetStatuses(c)
	if len(statuses) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, statuses)
}
