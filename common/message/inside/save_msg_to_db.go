package inside

import "encoding/json"

// MessageToDB 存库的消息封装体
type MessageToDB struct {
	Type    uint8             `json:"type"`    // 存库消息的类型
	Payload []json.RawMessage `json:"payload"` // 可以直接入库的json marshal后的消息
}
