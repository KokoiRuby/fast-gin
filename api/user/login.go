package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (API) LoginView(c *gin.Context) {
	c.String(http.StatusOK, "Login successfully")
	return
}
