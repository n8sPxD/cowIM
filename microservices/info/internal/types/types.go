// Code generated by goctl. DO NOT EDIT.
package types

type MessageRecordRequest struct {
}

type MessageRecordResponse struct {
}

type TimelineSyncInfo struct {
	SenderID  uint32 `json:"senderId"`
	GroupID   uint32 `json:"groupId,omitempty"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type TimelineSyncRequest struct {
	ID        uint32 `header:"UserID"`
	Timestamp int64  `form:Timestamp` // Unix 时间戳
}

type TimelineSyncResponse struct {
	Infos []TimelineSyncInfo `json:"infos"`
}
