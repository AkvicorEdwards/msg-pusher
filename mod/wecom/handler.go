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
	case "/modify":
		secretModify(w, r)
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
		_ = tplSecretInsert.Execute(w, map[string]any{"title": wecomName + " Secret Insert", "mod": self.URL(), "default_time": time.Now().AddDate(100, 0, 0).Format("2006-01-02T15:04")})
	} else if r.Method == "POST" {
		namePst := r.PostFormValue("name")
		if len(namePst) < 1 {
			app.RespAPIInvalidInput(w, "invalid name")
			return
		}
		corpIDPst := r.PostFormValue("corp_id")
		if len(corpIDPst) < 1 {
			app.RespAPIInvalidInput(w, "invalid corp id")
			return
		}
		agentIDPst := r.PostFormValue("agent_id")
		if len(agentIDPst) < 1 {
			app.RespAPIInvalidInput(w, "invalid agent id")
			return
		}
		agentIDInt, err := strconv.ParseInt(agentIDPst, 10, 64)
		if err != nil {
			app.RespAPIInvalidInput(w, "invalid agent id")
			return
		}
		secretPst := r.PostFormValue("secret")
		if len(secretPst) < 1 {
			app.RespAPIInvalidInput(w, "invalid secret")
			return
		}
		validityPeriodPst := r.PostFormValue("validity_period")
		validityPeriodTime, err := time.ParseInLocation("2006-01-02T15:04", validityPeriodPst, time.Local)
		if err != nil {
			app.RespAPIInvalidInput(w, "invalid validity period")
			return
		}
		glog.Debug("Date [%s] [%s] [%s]", validityPeriodPst, time.Unix(validityPeriodTime.Unix(), 0).UTC().Format("2006-01-02T15:04"), time.Unix(validityPeriodTime.Unix(), 0).Format("2006-01-02T15:04"))
		data := new(ModWecomModel)
		data.Name = namePst
		data.CorpID = corpIDPst
		data.AgentID = agentIDInt
		data.Secret = secretPst
		data.ValidityPeriod = validityPeriodTime.Unix()
		res := insertWecom(data)
		if !res {
			glog.Warning("failed to get wecom")
			app.RespAPIProcessingFailed(w, "")
			return
		}

		app.LastPage(w, r)
	}

}

func secretModify(w http.ResponseWriter, r *http.Request) {
	head, tail := util.SplitPathRepeat(r.URL.Path, 2)
	glog.Debug("[%-4s][%-32s] [%s][%s]", r.Method, path.Join(wecomUrl, "/secret/modify"), head, tail)

	if len(tail) <= 1 {
		return
	}

	id, err := strconv.ParseInt(tail[1:], 10, 64)
	if err != nil {
		return
	}

	sec := getWecomByID(id, false)
	if sec == nil {
		return
	}

	if r.Method == "GET" {
		validityPeriod := time.Unix(sec.ValidityPeriod, 0).Format("2006-01-02T15:04")
		createTime := time.Unix(sec.CreateTime, 0).Format("2006-01-02T15:04")
		lastUsed := time.Unix(sec.LastUsed, 0).Format("2006-01-02T15:04")
		expired := time.Unix(sec.Expired, 0).Format("2006-01-02T15:04")

		_ = tplSecretModify.Execute(w, map[string]any{"title": wecomName + " Secret Modify", "mod": self.URL(), "sec": sec,
			"validity_period": validityPeriod, "create_time": createTime, "last_used": lastUsed, "expired": expired})
	} else if r.Method == "POST" {
		corpIDPst := r.PostFormValue("corp_id")
		if len(corpIDPst) < 1 {
			app.RespAPIInvalidInput(w, "invalid corp id")
			return
		}
		sec.CorpID = corpIDPst
		agentIDPst := r.PostFormValue("agent_id")
		if len(agentIDPst) < 1 {
			app.RespAPIInvalidInput(w, "invalid agent id")
			return
		}
		agentIDInt, err := strconv.ParseInt(agentIDPst, 10, 64)
		if err != nil {
			app.RespAPIInvalidInput(w, "invalid agent id")
			return
		}
		sec.AgentID = agentIDInt
		secretPst := r.PostFormValue("secret")
		if len(secretPst) < 1 {
			app.RespAPIInvalidInput(w, "invalid secret")
			return
		}
		sec.Secret = secretPst
		namePst := r.PostFormValue("name")
		if len(namePst) < 1 {
			app.RespAPIInvalidInput(w, "invalid name")
			return
		}
		sec.Name = namePst
		validityPeriodPst := r.PostFormValue("validity_period")
		validityPeriodTime, err := time.ParseInLocation("2006-01-02T15:04", validityPeriodPst, time.Local)
		if err != nil {
			app.RespAPIInvalidInput(w, "invalid validity period")
			return
		}
		sec.ValidityPeriod = validityPeriodTime.Unix()
		expiredPst := r.PostFormValue("expired")
		expiredTime, err := time.ParseInLocation("2006-01-02T15:04", expiredPst, time.Local)
		if err != nil {
			app.RespAPIInvalidInput(w, "invalid expired")
			return
		}
		sec.Expired = expiredTime.Unix()

		res := ModifyWecom(sec)
		if !res {
			app.RespAPIProcessingFailed(w, "")
			return
		}
		app.Reload(w, r)
	}
}
