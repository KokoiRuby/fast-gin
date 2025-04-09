package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (API) ListView(c *gin.Context) {
	c.String(http.StatusOK, "User List")
}
