package svc_redis

import (
	"context"
	"fast-gin/global"
	"fast-gin/utils/jwts"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

func Logout(token string) {
	claims, err := jwts.ValidateJWT(token)
	if err != nil {
		logrus.Errorf("Failed to validate JWT: %v", err)
		return
	}
	key := fmt.Sprintf("logout_%s", token)
	expiration := claims.ExpiresAt.Sub(time.Now())

	res, err := global.Redis.Set(context.Background(), key, "", expiration).Result()
	if err != nil {
		logrus.Errorf("Failed to set token in Redis: %v", err)
	}
	logrus.Debugf("Token blacklisted: %v", res)
}

func HasLoggedOut(token string) bool {
	key := fmt.Sprintf("logout_%s", token)
	_, err := global.Redis.Get(context.Background(), key).Result()
	if err == nil {
		return true
	}
	return false
}
