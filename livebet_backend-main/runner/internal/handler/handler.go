package handler

import (
	"livebets/runner/internal/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	commandService *service.CommandService
	statusService  *service.StatusService
}

func NewHandler(
	commandService *service.CommandService,
	statusService *service.StatusService,
) *Handler {
	return &Handler{
		statusService:  statusService,
		commandService: commandService,
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

	router.POST("/set-command", h.SetCommand)
	router.GET("/status", h.GetStatuses)

	return router
}
