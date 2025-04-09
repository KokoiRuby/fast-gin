package routers

import (
	"fast-gin/api"
	"fast-gin/middlewares"
	"github.com/gin-gonic/gin"
)

func CaptchaRouter(g *gin.RouterGroup) {
	captchaAPI := api.Apis.CaptchaAPI

	r := g.Group("captcha").Use(
		middlewares.AuthMiddlewareWithRole,
	)

	r.GET("/generate", captchaAPI.GenerateCaptcha)
}
