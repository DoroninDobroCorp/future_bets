package handler

import (
	"livebets/calculator/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReqCalcBet struct {
	UserID string                `json:"userId"`
	Pair   entity.PairOneOutcome `json:"pair"`
}

type ResCalcBet struct {
	UsersCount int                  `json:"usersCount"`
	CalcBet    entity.CalculatedBet `json:"calcBet"`
}

func (h *Handler) LogBetAccept(c *gin.Context) {
	var input entity.AcceptBet

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := h.logsService.LogBetAccept(c, input); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetCalcBet(c *gin.Context) {
	var input ReqCalcBet

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	calcBet, usersCount := h.logsService.CalcSumBet(c, input.UserID, input.Pair)

	c.JSON(http.StatusOK, &ResCalcBet{
		UsersCount: usersCount,
		CalcBet:    calcBet,
	})
}

func (h *Handler) LogTestBetAccept(c *gin.Context) {
	var input entity.AcceptBet

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := h.logsService.LogTestBetAccept(c, input); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
