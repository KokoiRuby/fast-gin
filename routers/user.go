package routers

import (
	"fast-gin/api"
	"github.com/gin-gonic/gin"
)

func UserRouter(g *gin.RouterGroup) {
	userAPI := api.Apis.UserAPI

	r := g.Group("users").Use()

	r.POST("/login", userAPI.LoginView)
}
