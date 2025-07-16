package handler

import (
	"livebets/auto_matcher/internal/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	handMatchService   *service.HandMatcherService
	autoMatchService   *service.AutoMatcherService
	onlineMatchService *service.OnlineMatcherService
}

func NewHandler(
	handMatchService *service.HandMatcherService,
	autoMatchService *service.AutoMatcherService,
	onlineMatchService *service.OnlineMatcherService,
) *Handler {
	return &Handler{
		handMatchService:   handMatchService,
		autoMatchService:   autoMatchService,
		onlineMatchService: onlineMatchService,
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

	hand_merge := router.Group("/hand-merge")
	{
		hand_merge.GET("/bookmakers", h.GetBookmakers)
		hand_merge.GET("/sports", h.GetSports)

		leagues := hand_merge.Group("/leagues")
		{
			leagues.GET("/", h.GetAllLeaguesByBookmaker)
			leagues.GET("/get-unmatch", h.GetUnMatchedLeagues)
			leagues.GET("/get-match", h.GetMatchedLeagues)
			leagues.POST("/create-pair", h.CreateNewLeaguePair)
			leagues.GET("/candidates", h.GetLeagueCandidates)
			leagues.GET("/online-unmatch", h.GetOnlineUnmatchLeagues)
			leagues.GET("/online-unmatch-prematch", h.GetOnlineUnmatchLeaguesPrematch)
		}

		teams := hand_merge.Group("/teams")
		{
			teams.GET("/get-unmatch", h.GetUnMatchedTeamsByLeagues)
			teams.GET("/get-match", h.GetMatchedTeamsByLeagues)
			teams.POST("/create-pair", h.CreateNewTeamPair)
			teams.GET("/candidates", h.GetTeamCandidates)
			teams.GET("/online-unmatch", h.GetOnlineUnmatchTeamsByLeagues)
			teams.GET("/online-unmatch-prematch", h.GetOnlineUnmatchTeamsByLeaguesPrematch)
		}
	}

	return router
}
