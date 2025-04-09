package user

import (
	"fast-gin/utils/response"
	"github.com/gin-gonic/gin"
)

func (API) LoginView(c *gin.Context) {
	response.OKWithMsg(c, "Login successfully")
	return
}
