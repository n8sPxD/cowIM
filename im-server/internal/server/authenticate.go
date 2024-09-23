package server

import (
	"fmt"
	"net/http"

	"github.com/n8sPxD/cowIM/common/jwt"
)

func (s *Server) authenticate(w http.ResponseWriter, r *http.Request) (*jwt.CustomClaims, bool) {
	// 先鉴权，再进行消息通讯
	// 从Authorization头部获取JWT令牌
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		return nil, false
	}

	// JWT 格式： Bearer <token>
	var tokenString string
	fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
	if tokenString == "" {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return nil, false
	}

	// 验证 JWT
	claims, err := jwt.ParseToken(tokenString, s.config.Auth.AccessSecret)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return nil, false
	}
	return claims, true
}
