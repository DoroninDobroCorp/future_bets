package handler

import (
	"livebets/analazer/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RespTimeRecords struct {
	ISave   int                          `json:"isave"`
	Records []entity.ResponsePriceRecord `json:"records"`
}

func (h *Handler) GetPriceRecordsByTime(c *gin.Context) {
	var input entity.ReqGetPriceRecordsByTime

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	iSave, records := h.priceService.GetPriceRecordsByTime(input)

	if iSave == -1 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, RespTimeRecords{iSave, records})
}

func (h *Handler) GetMatchData(c *gin.Context) {
	matchData := h.pairsMatchingService.GetMatchData(c)

	c.JSON(http.StatusOK, map[string]interface{}{
		"len":  len(matchData),
		"data": matchData,
	})
}

func (h *Handler) GetCacheKeys(c *gin.Context) {
	cacheKeys := h.pairsMatchingService.GetCacheKeys(c)

	c.JSON(http.StatusOK, map[string]interface{}{
		"len":  len(cacheKeys),
		"data": cacheKeys,
	})
}

func (h *Handler) GetCachePairs(c *gin.Context) {
	keys, values := h.pairsMatchingService.GetCachePairs(c)

	c.JSON(http.StatusOK, map[string]interface{}{
		"lenKeys":    len(keys),
		"lenValues":  len(values),
		"dataKeys":   keys,
		"dataValues": values,
	})
}

func (h *Handler) GetPairs(c *gin.Context) {
	pairs := h.pairsMatchingService.GetPairs(c)

	c.JSON(http.StatusOK, map[string]interface{}{
		"len":  len(pairs),
		"data": pairs,
	})
}

func (h *Handler) GetOnlineMatchData(c *gin.Context) {
	matchData := h.pairsMatchingService.GetOnlineMatchData(c)

	c.JSON(http.StatusOK, matchData)
}