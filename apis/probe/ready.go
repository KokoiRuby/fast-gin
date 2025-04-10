package probe

import (
	"fast-gin/utils/response"
	"github.com/gin-gonic/gin"
)

func (API) ReadyView(c *gin.Context) {
	response.OKWithMsg(c, "Ready")
	// TODO: Dependent services must be up: db, redis
}
