package routers

import (
	"fast-gin/api"
	"fast-gin/middlewares"
	"github.com/gin-gonic/gin"
)

func ImageRouter(g *gin.RouterGroup) {
	imageAPI := api.Apis.ImageAPI

	r := g.Group("images").Use(
		middlewares.AdminAuthMiddleware,
	)

	r.POST("upload", imageAPI.Upload)
}
