package constant

// User 权利
const (
	USER_ADMIN  = 0
	USER_COMMON = 1
	USER_SYSTEM = 99
)

// Message 类型
const (
	MSG_COMMON_MSG = iota
	MSG_REPLY_MSG
	MSG_IMAGE_MSG
	MSG_VIDEO_MSG
	MSG_FILE_MSG
	MSG_SYSTEM_MSG
)
