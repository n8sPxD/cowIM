// Code generated by goctl. DO NOT EDIT.
package types

type ChatListInfo struct {
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	RecentMsg string `json:"recentMsg"`
}

type ChatListRequest struct {
	UserID uint32 `header:X-User-ID`
}

type ChatListResponse struct {
	List []ChatListInfo `json:"list"`
}
