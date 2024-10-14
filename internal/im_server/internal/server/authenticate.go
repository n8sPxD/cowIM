package server

import (
	"net/http"

	"github.com/n8sPxD/cowIM/pkg/jwt"
)

func (s *Server) authenticate(w http.ResponseWriter, r *http.Request) (*jwt.CustomClaims, bool) {
	// 从 URL 查询参数中获取 JWT 令牌
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "Token is required", http.StatusUnauthorized)
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
