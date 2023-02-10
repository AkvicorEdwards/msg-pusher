package wecom

import (
	"fmt"
	"github.com/AkvicorEdwards/glog"
	"github.com/AkvicorEdwards/util"
	"msg-pusher/app"
	"msg-pusher/stl/pair"
	"net/http"
	"path"
	"strconv"
	"time"
)

func index(w http.ResponseWriter, r *http.Request) {
	head, tail := util.SplitPathRepeat(r.URL.Path, 2)
	glog.Debug("[%-4s][%-32s] [%s][%s]", r.Method, path.Join(wecomUrl, "/index"), head, tail)
	if !app.SessionVerify(r) {
		app.Login(w, r)
		return
	}

	_ = tplIndex.Execute(w, map[string]any{"title": wecomName, "mod_url": self.URL()})
}

func secret(w http.ResponseWriter, r *http.Request) {
	head, tail := util.SplitPathRepeat(r.URL.Path, 2)
	glog.Debug("[%-4s][%-32s] [%s][%s]", r.Method, path.Join(wecomUrl, "/secret"), head, tail)
	if !app.SessionVerify(r) {
		app.Login(w, r)
		return
	}

	switch head {
	case "/insert":
		secretInsert(w, r)
		return
	}

	if r.Method == "GET" {
		wecom := getWecom()
		secrets := make([]pair.Url, len(wecom))
		for k, v := range wecom {
			secrets[k].Name = v.Name
			secrets[k].Url = fmt.Sprint(v.ID)
		}
		_ = tplSecret.Execute(w, map[string]any{"title": wecomName + " Secret", "mod": self.URL(), "secrets": secrets})
	}

}

func secretInsert(w http.ResponseWriter, r *http.Request) {
	head, tail := util.SplitPathRepeat(r.URL.Path, 2)
	glog.Debug("[%-4s][%-32s] [%s][%s]", r.Method, path.Join(wecomUrl, "/secret/insert"), head, tail)

	if r.Method == "GET" {
		_ = tplSecretInsert.Execute(w, map[string]any{"title": wecomName + " Secret Insert", "mod": self.URL()})
	} else if r.Method == "POST" {
		namePst := r.PostFormValue("name")
		if len(namePst) < 1 {
			app.RespAPIInvalidInput(w)
			return
		}
		corpIDPst := r.PostFormValue("corp_id")
		if len(corpIDPst) < 1 {
			app.RespAPIInvalidInput(w)
			return
		}
		agentIDPst := r.PostFormValue("agent_id")
		if len(agentIDPst) < 1 {
			app.RespAPIInvalidInput(w)
			return
		}
		agentIDInt, err := strconv.ParseInt(agentIDPst, 10, 64)
		if err != nil {
			app.RespAPIInvalidInput(w)
			return
		}
		secretPst := r.PostFormValue("secret")
		if len(secretPst) < 1 {
			app.RespAPIInvalidInput(w)
			return
		}
		validityPeriodPst := r.PostFormValue("validity_period")
		validityPeriodTime, err := time.ParseInLocation("2006-01-02T00:00", validityPeriodPst, time.Local)
		if err != nil {
			app.RespAPIInvalidInput(w)
			return
		}
		data := new(ModWecomModel)
		data.Name = namePst
		data.CorpID = corpIDPst
		data.AgentID = agentIDInt
		data.Secret = secretPst
		data.ValidityPeriod = validityPeriodTime.Unix()
		res := insertWecom(data)
		if !res {
			glog.Warning("failed to get wecom")
			app.RespAPIProcessingFailed(w)
			return
		}

		app.LastPage(w, r)
	}

}
