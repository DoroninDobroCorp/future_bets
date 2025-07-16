package handler

import (
	"livebets/calculator/internal/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	logsService *service.LogsService
}

func NewHandler(
	logsService *service.LogsService,
) *Handler {
	return &Handler{
		logsService: logsService,
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

	router.POST("/log-bet-accept", h.LogBetAccept)
	router.POST("/log-test-bet-accept", h.LogTestBetAccept)
	router.POST("/calc-bet", h.GetCalcBet)

	return router
}
