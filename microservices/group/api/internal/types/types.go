// Code generated by goctl. DO NOT EDIT.
package types

type GroupCreateRequest struct {
	Groupname string `json:"groupname"`
}

type GroupCreateResponse struct {
	GroupID uint32 `json:"groupID"`
}

type GroupInfoRequest struct {
	GroupID uint32 `json:"groupID`
}

type GroupInfoResponse struct {
	Groupname    string   `json:"groupname"`
	GroupMembers []uint32 `json:"groupMembers"`
}

type GroupInviteRequest struct {
	UserID  uint32   `header:UserID`
	GroupID uint32   `json:"groupId"`
	Members []uint32 `json:"members"`
}

type GroupInviteResponse struct {
}

type GroupJoinedRequest struct {
	UserID uint32 `header:UserID`
}

type GroupJoinedResponse struct {
	GroupID []uint32 `json:"groupId"`
}
