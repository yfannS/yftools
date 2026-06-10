package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtSecret string

// SetSecret 初始化时由外部调用，将配置文件中的 secret 注入
func SetSecret(secret string) {
	jwtSecret = secret
}

func getSecret() []byte {
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-me"
	}
	return []byte(jwtSecret)
}

func ResolveExpire(expireStr string) time.Duration {
	if expireStr == "" {
		expireStr = "168h"
	}

	expire, err := time.ParseDuration(expireStr)
	if err != nil || expire <= 0 {
		return 7 * 24 * time.Hour
	}

	return expire
}

func GenerateToken(userID int64, username string, expireStr string) (string, error) {
	expire := ResolveExpire(expireStr)

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecret())
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getSecret(), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
