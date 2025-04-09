package api

import (
	"fast-gin/api/captcha"
	"fast-gin/api/image"
	"fast-gin/api/user"
)

type APIs struct {
	UserAPI    user.API
	ImageAPI   image.API
	CaptchaAPI captcha.API
}

var Apis = new(APIs)
