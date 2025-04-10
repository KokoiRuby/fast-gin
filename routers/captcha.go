package routers

import (
	"fast-gin/apis"
	"fast-gin/middlewares"
	"github.com/gin-gonic/gin"
)

func CaptchaRouter(g *gin.RouterGroup) {
	captchaAPI := apis.Apis.CaptchaAPI

	r := g.Group("captcha").Use(
		middlewares.AdminAuthMiddleware,
	)

	r.GET("generate", captchaAPI.GenerateCaptcha)
}
