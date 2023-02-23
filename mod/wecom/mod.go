package wecom

import (
	"github.com/AkvicorEdwards/glog"
	"github.com/AkvicorEdwards/util"
	"msg-pusher/app"
	"msg-pusher/mod"
	"net/http"
)

type Model struct {
	handler map[string]func(w http.ResponseWriter, r *http.Request)
}

func (m *Model) URL() string {
	return wecomUrl
}

func (m *Model) Table() string {
	return wecomTable
}

func (m *Model) Name() string {
	return wecomName
}

func (m *Model) Key() string {
	return wecomKey
}

func (m *Model) Handle(w http.ResponseWriter, r *http.Request) {
	if !app.SessionVerify(r) {
		app.Login(w, r)
		return
	}

	var handler func(w http.ResponseWriter, r *http.Request)
	var ok bool

	head, tail := util.SplitPathRepeat(r.URL.Path, 1)
	glog.Debug("[%-4s][%-32s] [%s][%s]", r.Method, wecomUrl, head, tail)

	handler, ok = self.handler[head]
	if ok {
		handler(w, r)
		return
	}
}

func (m *Model) Prepare(id int64, data *mod.MessageModel) mod.Package {
	wecom := getWecomByID(id, true)
	if wecom == nil {
		return nil
	}
	glog.Debug("prepare wecom: [%#v]", wecom)
	glog.Debug("prepare data: [%s]", data.String())

	extra := ParseExtra(data.Extra[m.Key()])

	var msg Message
	switch extra.Type {
	case TypeText:
		msg = NewText(extra.ToUser, extra.ToParty, extra.ToTag, wecom.AgentID, data.Content)
	case TypeTextCard:
		fallthrough
	default:
		msg = NewTextCard(extra.ToUser, extra.ToParty, extra.ToTag, wecom.AgentID, data.Title, data.Content, extra.Url, extra.Btn)
	}
	glog.Debug("prepare [%s]", msg.String())
	pkg := &Package{
		wecom: wecom,
		msg:   msg.Bytes(),
	}
	return pkg
}

func (m *Model) GetTarget(id int64) mod.Target {
	return getWecomByID(id, false)
}
