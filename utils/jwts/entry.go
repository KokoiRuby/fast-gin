package jwts

import (
	"fast-gin/global"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"time"
)

// CustomClaims defines the structure of the token's payload

type ClaimMeta struct {
	UserID uint `json:"userID"`
	RoleID int8 `json:"roleID"`
}

type CustomClaims struct {
	ClaimMeta
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token
func GenerateJWT(meta ClaimMeta) (string, error) {
	// Define claims
	claims := CustomClaims{
		ClaimMeta: meta,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(global.Config.JWT.Expire) * time.Hour)), // Expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    global.Config.JWT.Issuer,
		},
	}

	// Create token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(global.Config.JWT.SecretKey))
	if err != nil {
		logrus.Errorf("Error signing jwt: %v", err)
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT parses and validates a JWT token
func ValidateJWT(tokenString string) (*CustomClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(global.Config.JWT.SecretKey), nil
	})

	if err != nil {
		logrus.Errorf("Error parsing jwt: %v", err)
		return nil, err
	}

	// Extract claims if token is valid
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	logrus.Errorf("Error parsing jwt: %v", err)
	return nil, fmt.Errorf("invalid token")
}
