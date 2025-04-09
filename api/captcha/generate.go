package captcha

import (
	"fast-gin/utils/captcha"
	"fast-gin/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
)

type GenerateCaptchaResponse struct {
	CaptchaID string `json:"captchaID"`
	Captcha   string `json:"captcha"`
}

func (API) GenerateCaptcha(c *gin.Context) {
	var driver = base64Captcha.DriverString{
		Height:          80,
		Width:           240,
		NoiseCount:      4,
		ShowLineOptions: 6,
		Length:          4,
		Source:          "0123456789",
	}

	capt := base64Captcha.NewCaptcha(&driver, captcha.CaptchaStore)
	id, b64s, _, err := capt.Generate()
	if err != nil {
		logrus.Errorf("Failed to generate captcha: %v", err)
		response.FailWithMsg(c, "Failed to generate captcha")
		return
	}

	response.OKWithData(c, GenerateCaptchaResponse{
		CaptchaID: id,
		Captcha:   b64s,
	})
}
