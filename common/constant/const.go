package constant

// Message 类型
const (
	_ = iota
	MSG_COMMON_MSG
	MSG_REPLY_MSG
	MSG_IMAGE_MSG
	MSG_VIDEO_MSG
	MSG_FILE_MSG
)

// 系统相关 Message 类型
const (
	MSG_SYSTEM_MSG = 100 + iota
	MSG_ACK_MSG    // 消息重传中使用
)

// 群聊或者单聊
const (
	_ = iota
	SINGLE_CHAT
	GROUP_CHAT
	BIG_GROUP_CHAT
	SYSTEM_INFO = 99
)

// 存表中区分表的类型
const (
	_ = iota
	MESSAGE_RECORD
	MESSAGE_SYNC
	USER_TIMELINE
)

// 特殊用户
const (
	_           = iota
	USER_SYSTEM // 系统消息
	USER_COMMON
)

// 群聊权限
const (
	_ = iota
	GROUP_MASTER
	GROUP_ADMIN
	GROUP_COMMON
)
