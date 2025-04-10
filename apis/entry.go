package apis

import (
	"fast-gin/apis/captcha"
	"fast-gin/apis/image"
	"fast-gin/apis/probe"
	"fast-gin/apis/user"
)

type APIs struct {
	UserAPI    user.API
	ImageAPI   image.API
	CaptchaAPI captcha.API
	ProbeAPI   probe.API
}

var Apis = new(APIs)
