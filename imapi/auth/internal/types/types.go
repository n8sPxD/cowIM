// Code generated by goctl. DO NOT EDIT.
package types

type LoginRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type RegisterResponse struct {
}
