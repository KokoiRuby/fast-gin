package user

import (
	"fast-gin/global"
	"fast-gin/middlewares"
	"fast-gin/models"
	"fast-gin/utils/captcha"
	"fast-gin/utils/jwts"
	"fast-gin/utils/pwd"
	"fast-gin/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LoginRequest struct {
	Username   string `json:"username" binding:"required" label:"uname"`
	Password   string `json:"password" binding:"required" label:"pwd"`
	CaptchaID  string `json:"captchaID"`
	CaptchaAns string `json:"captchaAns"`
}

func (API) LoginView(c *gin.Context) {
	req := middlewares.GetBind[LoginRequest](c)

	// 1. Validate captcha
	if global.Config.Site.Login.Captcha {
		if req.CaptchaID == "" || req.CaptchaAns == "" {
			response.FailWithMsg(c, "Captcha is required")
			return
		}
		if !captcha.CaptchaStore.Verify(req.CaptchaID, req.CaptchaAns, true) {
			response.FailWithMsg(c, "Failed to validate captcha")
			return
		}
	}

	// 2. Get user from DB
	var user models.UserModel
	err := global.DB.Take(&user, "username = ?", req.Username).Error
	if err != nil {
		response.FailWithMsg(c, "Username or Password is incorrect")
		return
	}

	// 3. Validate password
	if !pwd.Validate(user.Password, req.Password) {
		response.FailWithMsg(c, "Username or Password is incorrect")
		return
	}

	// 4. Issue token
	token, err := jwts.GenerateJWT(jwts.ClaimMeta{
		UserID: user.ID,
		RoleID: user.RoleID,
	})
	if err != nil {
		logrus.Errorf("Failed to generate JWT token: %v", err)
		response.FailWithMsg(c, "Failed to login")
	}

	response.OKWithData(c, token)
	return
}
