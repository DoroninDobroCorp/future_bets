package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetOnlineUnmatchLeagues(c *gin.Context) {
	sportName := c.Query("sportName")
	if sportName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	firstBookmakerName := c.Query("firstBookmakerName")
	if firstBookmakerName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	secondBookmakerName := c.Query("secondBookmakerName")
	if secondBookmakerName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if firstBookmakerName == secondBookmakerName {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	leagues, err := h.onlineMatchService.GetOnlineUnmatchLeagues(c, sportName, firstBookmakerName, secondBookmakerName)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": leagues,
	})
}

func (h *Handler) GetOnlineUnmatchTeamsByLeagues(c *gin.Context) {
	sportName := c.Query("sportName")
	if sportName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	firstBookmakerName := c.Query("firstBookmakerName")
	if firstBookmakerName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	secondBookmakerName := c.Query("secondBookmakerName")
	if secondBookmakerName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if firstBookmakerName == secondBookmakerName {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	teams, err := h.onlineMatchService.GetOnlineUnmatchTeams(c, sportName, firstBookmakerName, secondBookmakerName)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": teams,
	})
}

func (h *Handler) GetOnlineUnmatchLeaguesPrematch(c *gin.Context) {
	sportName := c.Query("sportName")
	if sportName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	firstBookmakerName := c.Query("firstBookmakerName")
	if firstBookmakerName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	secondBookmakerName := c.Query("secondBookmakerName")
	if secondBookmakerName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if firstBookmakerName == secondBookmakerName {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	leagues, err := h.onlineMatchService.GetOnlineUnmatchLeaguesPrematch(c, sportName, firstBookmakerName, secondBookmakerName)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": leagues,
	})
}

func (h *Handler) GetOnlineUnmatchTeamsByLeaguesPrematch(c *gin.Context) {
	sportName := c.Query("sportName")
	if sportName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	firstBookmakerName := c.Query("firstBookmakerName")
	if firstBookmakerName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	secondBookmakerName := c.Query("secondBookmakerName")
	if secondBookmakerName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if firstBookmakerName == secondBookmakerName {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	teams, err := h.onlineMatchService.GetOnlineUnmatchTeamsPrematch(c, sportName, firstBookmakerName, secondBookmakerName)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": teams,
	})
}