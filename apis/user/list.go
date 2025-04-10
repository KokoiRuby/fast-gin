package user

import (
	"fast-gin/middlewares"
	"fast-gin/models"
	"fast-gin/service/common"
	"fast-gin/utils/response"
	"github.com/gin-gonic/gin"
)

func (API) ListView(c *gin.Context) {
	req := middlewares.GetBind[models.PageInfo](c)

	list, count, _ := common.QueryList(models.UserModel{}, common.QueryOption{
		PageInfo: req,
		Likes:    []string{"username", "nickname"},
		Debug:    true,
	})

	response.OKWithList(c, list, count)
}
