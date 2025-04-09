package user

import (
	"fast-gin/utils/response"
	"github.com/gin-gonic/gin"
)

func (API) ListView(c *gin.Context) {
	response.OKWithData(c, "User List")
}
