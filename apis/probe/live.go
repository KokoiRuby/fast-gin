package probe

import (
	"fast-gin/utils/response"
	"github.com/gin-gonic/gin"
)

func (API) LiveView(c *gin.Context) {
	response.OKWithMsg(c, "Alive")
}
