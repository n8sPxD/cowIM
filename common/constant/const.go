package constant

// User Privilege
const (
	USER_ADMIN  = 0
	USER_COMMON = 1
)

// RPC Command Type
const (
	_ = iota
	SINGLE_CHAT_REQ
	GROUP_CHAT_REQ
	SINGLE_CHAT_RESP
	GROUP_CHAT_RESP
	//HEART_BEAT_REQ  = 98
	//HEART_BEAT_RESP = 99
)

// Message Type
const (
	COMMON_MSG = iota
	REPLY_MSG
	IMAGE_MSG
	VIDEO_MSG
	FILE_MSG
)

// Single or group message
const (
	NONE_GROUP  = -1
	NONE_SINGLE = ""
)
