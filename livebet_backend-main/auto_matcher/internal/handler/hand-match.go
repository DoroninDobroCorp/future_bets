package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllLeaguesByBookmaker(c *gin.Context) {
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

	leagues, err := h.handMatchService.GetAllLeaguesByBookmaker(c, sportName, firstBookmakerName, secondBookmakerName)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if len(leagues) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": leagues,
	})
}

func (h *Handler) GetUnMatchedLeagues(c *gin.Context) {
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

	leagues, err := h.handMatchService.GetUnMachedLeagues(c, sportName, firstBookmakerName, secondBookmakerName)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if len(leagues) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": leagues,
	})
}

func (h *Handler) GetMatchedLeagues(c *gin.Context) {
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

	pairs, err := h.handMatchService.GetMatchedLeagues(c, sportName, firstBookmakerName, secondBookmakerName)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if len(pairs) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": pairs,
	})
}

func (h *Handler) GetUnMatchedTeamsByLeagues(c *gin.Context) {
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

	pairs, err := h.handMatchService.GetUnMatchedTeamsByLeagues(c, sportName, firstBookmakerName, secondBookmakerName)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if len(pairs) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": pairs,
	})
}

func (h *Handler) GetMatchedTeamsByLeagues(c *gin.Context) {
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

	pairs, err := h.handMatchService.GetMatchedTeamsByLeagues(c, sportName, firstBookmakerName, secondBookmakerName)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if len(pairs) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": pairs,
	})
}

func (h *Handler) CreateNewTeamPair(c *gin.Context) {
	strFirstTeamID := c.Query("firstTeamID")
	if strFirstTeamID == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	firstTeamID, err := strconv.ParseInt(strFirstTeamID, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	strSecondTeamID := c.Query("secondTeamID")
	if strSecondTeamID == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	secondTeamID, err := strconv.ParseInt(strSecondTeamID, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if firstTeamID == secondTeamID {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	success, err := h.handMatchService.CreateTeamsPair(c, firstTeamID, secondTeamID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !success {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) CreateNewLeaguePair(c *gin.Context) {
	strFirstLeagueID := c.Query("firstLeagueID")
	if strFirstLeagueID == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	firstLeagueID, err := strconv.ParseInt(strFirstLeagueID, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	strSecondLeagueID := c.Query("secondLeagueID")
	if strSecondLeagueID == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	secondLeagueID, err := strconv.ParseInt(strSecondLeagueID, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if firstLeagueID == secondLeagueID {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	success, err := h.handMatchService.CreateLeaguesPair(c, firstLeagueID, secondLeagueID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !success {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	c.Status(http.StatusOK)
}
