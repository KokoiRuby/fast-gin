package routers

import (
	"fast-gin/apis"
	"fast-gin/middlewares"
	"github.com/gin-gonic/gin"
)

func ImageRouter(g *gin.RouterGroup) {
	imageAPI := apis.Apis.ImageAPI

	r := g.Group("images").Use(
		middlewares.AdminAuthMiddleware,
	)

	r.POST("upload", imageAPI.UploadView)
}
