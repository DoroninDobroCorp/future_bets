package handler

import (
	"livebets/analazer/internal/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	priceService         *service.PriceService
	pairsMatchingService *service.PairsMatchingService
}

func NewHandler(
	priceService *service.PriceService,
	pairsMatchingService *service.PairsMatchingService,
) *Handler {
	return &Handler{
		priceService:         priceService,
		pairsMatchingService: pairsMatchingService,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "Cookie", "Content-Length", "X-CSRF-Token", "Accept-Encoding", "Cache-Control"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 4 * time.Hour,
	}))

	router.GET("/prices-by-time", h.GetPriceRecordsByTime) // for calculator
	router.GET("/online-match-data", h.GetOnlineMatchData) // for automatcher

	router.GET("/match-data", h.GetMatchData)
	router.GET("/cache-keys", h.GetCacheKeys)
	router.GET("/cache-pairs", h.GetCachePairs)
	router.GET("/pairs", h.GetPairs)

	return router
}
