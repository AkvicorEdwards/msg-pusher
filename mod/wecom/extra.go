package wecom

import "encoding/json"

type ExtraModel struct {
	Type    string `json:"type"`
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty"`
	ToTag   string `json:"totag"`
	Url     string `json:"url"`
	Btn     string `json:"btn"`
}

func (m *ExtraModel) String() string {
	extraStr, err := json.Marshal(m)
	if err != nil {
		return `{}`
	}
	return string(extraStr)
}

func ParseExtra(data map[string]string) *ExtraModel {
	extra := new(ExtraModel)
	extra.Type = data["type"]
	if len(extra.Type) == 0 {
		extra.Type = defaultMessageType
	}
	extra.ToUser = data["touser"]
	extra.ToParty = data["toparty"]
	extra.ToTag = data["totag"]
	if len(extra.ToUser) == 0 && len(extra.ToParty) == 0 && len(extra.ToTag) == 0 {
		extra.ToUser = defaultMessageToUser
	}
	extra.Url = data["url"]
	if len(extra.Url) == 0 {
		extra.Url = defaultMessageURL
	}
	extra.Btn = data["btn"]
	if len(extra.Btn) == 0 {
		extra.Btn = defaultMessageButton
	}
	return extra
}
