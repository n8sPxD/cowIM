package constant

// Message 类型
const (
	MSG_COMMON_MSG = iota
	MSG_REPLY_MSG
	MSG_IMAGE_MSG
	MSG_VIDEO_MSG
	MSG_FILE_MSG
	MSG_SYSTEM_MSG = 99
)

// 群聊或者单聊
const (
	SINGLE_CHAT = iota
	GROUP_CHAT
	BIG_GROUP_CHAT
	SYSTEM_INFO = 99
)

// 存表中区分表的类型
const (
	MESSAGE_RECORD = iota
	MESSAGE_SYNC
	USER_TIMELINE
)

// 特殊用户
const (
	USER_SYSTEM = iota // 系统消息
	USER_COMMON
)

// 群聊权限
const (
	GROUP_MASTER = iota
	GROUP_ADMIN
	GROUP_COMMON
)
