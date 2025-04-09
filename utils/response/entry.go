package response

import (
	"fast-gin/utils/validate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func OK(c *gin.Context, data any, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Data: data,
		Msg:  msg,
	})
}

func OKWithData(c *gin.Context, data any) {
	OK(c, data, "Success")
}

func OKWithMsg(c *gin.Context, msg string) {
	OK(c, gin.H{}, msg)
}

func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: gin.H{},
		Msg:  msg,
	})
}

func FailWithMsg(c *gin.Context, msg string) {
	Fail(c, 7, msg)
}

func FailWithErr(c *gin.Context, err error) {
	msg := validate.ValidateError(err)
	Fail(c, 7, msg)
}
