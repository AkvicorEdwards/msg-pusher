package wecom

import "encoding/json"

const TypeTextCard = "textcard"

type TextCardModel struct {
	ToUser   string               `json:"touser"`
	ToParty  string               `json:"toparty"`
	ToTag    string               `json:"totag"`
	MsgType  string               `json:"msgtype"`
	AgentId  int64                `json:"agentid"`
	TextCard TextCardContentModel `json:"textcard"`
}

type TextCardContentModel struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	BtnTxt      string `json:"btntxt"`
}

func NewTextCard(toUser, toParty, toTag string, agentId int64, title, description, url, btn string) *TextCardModel {
	return &TextCardModel{
		ToUser:  toUser,
		ToParty: toParty,
		ToTag:   toTag,
		MsgType: "textcard",
		AgentId: agentId,
		TextCard: TextCardContentModel{
			Title:       title,
			Description: description,
			URL:         url,
			BtnTxt:      btn,
		},
	}
}

func (t *TextCardModel) String() string {
	msgStr, err := json.Marshal(t)
	if err != nil {
		return `{}`
	}
	return string(msgStr)
}

func (t *TextCardModel) Bytes() []byte {
	msg, err := json.Marshal(t)
	if err != nil {
		return []byte(`{}`)
	}
	return msg
}
