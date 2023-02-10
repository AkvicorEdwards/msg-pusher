package mod

import (
	"encoding/json"
)

type MessageUrgency int

const (
	UrgencyInstruction MessageUrgency = 0
	UrgencyUrgent      MessageUrgency = 1
	UrgencyImportant   MessageUrgency = 2
	UrgencyGeneral     MessageUrgency = 3
)

type MessageModel struct {
	Title      string         `json:"title"`
	Content    string         `json:"content"`
	Sender     string         `json:"sender"`
	Urgency    MessageUrgency `json:"urgency"`
	TimeCreate int64          `json:"time_create"`
	// TimeSend is expected delivery time of the message
	TimeSend int64 `json:"time_send"`
	// Extra is map[mod.key]map[key]value
	Extra map[string]map[string]string `json:"extra"`
}

func (m *MessageModel) String() string {
	msgStr, err := json.Marshal(m)
	if err != nil {
		return `{}`
	}
	return string(msgStr)
}

func ParseMessage(data []byte) *MessageModel {
	msg := new(MessageModel)
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil
	}
	return msg
}
