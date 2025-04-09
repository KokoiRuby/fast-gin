package routers

import (
	"fast-gin/api"
	"fast-gin/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRouter(g *gin.RouterGroup) {
	userAPI := api.Apis.UserAPI

	r := g.Group("users").Use(middlewares.LimitMiddleware(10), middlewares.AuthMiddlewareWithRole)

	r.POST("/login", userAPI.LoginView)
	r.GET("/list", userAPI.ListView)
}
