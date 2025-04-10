package user

import (
	"fast-gin/global"
	"fast-gin/service/svc_redis"
	"fast-gin/utils/response"
	"github.com/gin-gonic/gin"
)

func (API) LogoutView(c *gin.Context) {
	token := c.GetHeader("token")
	if global.Redis == nil {
		response.OKWithMsg(c, "Logout successfully")
		return
	}

	svc_redis.Logout(token)
	response.OKWithMsg(c, "Logout successfully")
	return
}
