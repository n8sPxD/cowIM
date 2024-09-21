package inside

import "encoding/json"

type Message struct {
	Type    uint8           `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
