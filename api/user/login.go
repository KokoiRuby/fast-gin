package user

import (
	"fast-gin/middlewares"
	"fast-gin/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required" label:"uname"`
	Password string `json:"password" binding:"required" label:"pwd"`
}

func (API) LoginView(c *gin.Context) {
	req := middlewares.GetBind[LoginRequest](c)
	fmt.Println(req)

	response.OKWithMsg(c, "Login successfully")
	return
}
