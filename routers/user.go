package routers

import (
	"fast-gin/api"
	"fast-gin/api/user"
	"fast-gin/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRouter(g *gin.RouterGroup) {
	userAPI := api.Apis.UserAPI

	r := g.Group("users").Use(
		middlewares.LimitMiddleware(1),
		middlewares.BindJsonMiddleware[user.LoginRequest],
		middlewares.AuthMiddlewareWithRole,
	)

	r.POST("/login", userAPI.LoginView)
	r.GET("/list", userAPI.ListView)
}
