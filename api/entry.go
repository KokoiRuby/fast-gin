package api

import (
	"fast-gin/api/image"
	"fast-gin/api/user"
)

type APIs struct {
	UserAPI  user.API
	ImageAPI image.API
}

var Apis = new(APIs)
