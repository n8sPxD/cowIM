// Code generated by goctl. DO NOT EDIT.
package types

type MessageRecordRequest struct {
}

type MessageRecordResponse struct {
}

type TimelineSyncInfo struct {
	SenderID uint32 `json:"senderId"`
	GroupID  uint32 `json:"groupId,omitempty"`
	Message  string `json:"message"`
}

type TimelineSyncRequest struct {
	ID        string `header:"UserID"`
	Timestamp int64  `header:Timestamp` // Unix 时间戳
}

type TimelineSyncResponse struct {
	Infos []TimelineSyncInfo `json:"infos"`
}
