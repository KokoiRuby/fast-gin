package api

import "fast-gin/api/user"

type APIs struct {
	UserAPI user.API
}

var Apis = new(APIs)
