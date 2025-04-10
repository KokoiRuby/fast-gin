package routers

import (
	"fast-gin/apis"
	"fast-gin/middlewares"
	"github.com/gin-gonic/gin"
)

func ProbeRouter(g *gin.RouterGroup) {
	probeAPI := apis.Apis.ProbeAPI

	g.GET("/liveness", middlewares.LimitMiddleware(1), probeAPI.LiveView)
	g.GET("/readiness", middlewares.LimitMiddleware(1), probeAPI.ReadyView)

}
