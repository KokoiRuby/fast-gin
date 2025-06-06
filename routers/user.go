package routers

import (
	"fast-gin/apis"
	"fast-gin/apis/user"
	"fast-gin/middlewares"
	"fast-gin/models"
	"github.com/gin-gonic/gin"
)

func UserRouter(g *gin.RouterGroup) {
	userAPI := apis.Apis.UserAPI

	r := g.Group("users").Use(
		middlewares.LimitMiddleware(1),
		middlewares.AdminAuthMiddleware,
	)

	r.POST("login", middlewares.BindJsonMiddleware[user.LoginRequest], userAPI.LoginView)
	r.POST("logout", userAPI.LogoutView)

	r.GET("list", middlewares.BindQueryMiddleware[models.PageInfo], userAPI.ListView)
}
