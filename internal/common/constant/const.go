package constant

// Message 类型
const (
	_              = iota
	MSG_COMMON_MSG // 一般消息
	MSG_REPLY_MSG  // 回复消息
	MSG_IMAGE_MSG  // 图片消息
	MSG_VIDEO_MSG  // 视频消息
	MSG_FILE_MSG   // 文件消息
)

// 系统相关 Message 类型
const (
	MSG_ALERT_MSG         = 100 + iota // 在前端直接弹窗出来的消息
	MSG_ACK_MSG                        // 消息重传中使用
	MSG_BIG_GROUP_REQ                  // 详细见 后端设计
	MSG_BIG_GROUP_ALL_REQ              // 详细见 后端设计
)

// 群聊或者单聊
const (
	_              = iota
	SINGLE_CHAT         // 单聊
	GROUP_CHAT          // 群聊
	BIG_GROUP_CHAT      // 大规模群聊
	SYSTEM_INFO    = 99 // 系统聊天，一般和USER_SYSTEM捆绑使用
)

// 存表中区分表的类型
const (
	_              = iota
	MESSAGE_RECORD // 消息记录表
	MESSAGE_SYNC   // 消息同步表，没啥用
	USER_TIMELINE  // 用户Timeline表
)

// 特殊用户
const (
	_           = iota
	USER_SYSTEM // 系统消息
	USER_COMMON // 普通用户
)

// 群聊权限
const (
	_            = iota
	GROUP_MASTER // 群主
	GROUP_ADMIN  // 管理
	GROUP_COMMON // 狗群员
)

// 系统对系统的操作 属于MsgType
const (
	MSG_DUP_CLIENT = 200 + iota // 多客户端登录，强制下线
)

// Websocket 相关代码
const (
	DUP_CLIENT_CODE = 4001 + iota
)

const (
	DUP_CLIENT_ERR = "您已在另一台客户端登陆，即将强制下线"
)
