package wecom

import "encoding/json"

const TypeText = "text"

type TextModel struct {
	ToUser  string           `json:"touser"`
	ToParty string           `json:"toparty"`
	ToTag   string           `json:"totag"`
	MsgType string           `json:"msgtype"`
	AgentId int64            `json:"agentid"`
	Text    TextContentModel `json:"text"`
}

type TextContentModel struct {
	Content string `json:"content"`
}

func NewText(toUser, toParty, toTag string, agentId int64, content string) *TextModel {
	return &TextModel{
		ToUser:  toUser,
		ToParty: toParty,
		ToTag:   toTag,
		MsgType: "text",
		AgentId: agentId,
		Text:    TextContentModel{Content: content},
	}
}

func (t *TextModel) String() string {
	msgStr, err := json.Marshal(t)
	if err != nil {
		return `{}`
	}
	return string(msgStr)
}

func (t *TextModel) Bytes() []byte {
	msg, err := json.Marshal(t)
	if err != nil {
		return []byte(`{}`)
	}
	return msg
}
