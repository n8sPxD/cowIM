package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type PayLoad struct {
	ID       uint32 `json:"ID"`
	Username string `json:"username"`
}

type CustomClaims struct {
	PayLoad
	jwt.RegisteredClaims
}

// GenToken 生成JWT token
func GenToken(payLoad PayLoad, accessSecret string, expires int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		PayLoad: payLoad,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expires))),
		},
	})
	return token.SignedString([]byte(accessSecret))
}

// ParseToken 解析JWT token
func ParseToken(tokenString, accessSecret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token非法")
}
