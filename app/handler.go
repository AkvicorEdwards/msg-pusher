package app

import (
	"fmt"
	"github.com/AkvicorEdwards/glog"
	"github.com/AkvicorEdwards/util"
	"msg-pusher/config"
	"msg-pusher/db"
	"msg-pusher/mod"
	sendpkg "msg-pusher/send"
	"msg-pusher/stl/pair"
	"msg-pusher/tpl"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func staticFavicon(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write(tpl.Favicon)
}

func index(w http.ResponseWriter, r *http.Request) {
	head, tail := util.SplitPathRepeat(r.URL.Path, 1)
	glog.Debug("[%-4s][%-32s] [%s][%s]", r.Method, "/index", head, tail)
	if !SessionVerify(r) {
		Login(w, r)
		return
	}

	if r.Method == http.MethodGet {
		mods := make([]pair.Url, 0, len(Global.url))
		for _, v := range Global.url {
			mods = append(mods, pair.Url{Url: v.URL(), Name: v.Name()})
		}
		_ = tpl.Index.Execute(w, map[string]interface{}{"title": "Message Pusher", "mods": mods})
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	head, tail := util.SplitPathRepeat(r.URL.Path, 1)
	glog.Debug("[%-4s][%-32s] [%s][%s]", r.Method, "/login", head, tail)

	if r.Method == http.MethodGet {
		_ = tpl.Login.Execute(w, map[string]interface{}{"title": "Login"})
		return
	} else if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == config.Global.Security.Username && password == config.Global.Security.Password {
			glog.Info("Login successful [%s]", username)
			sessionUpdate(w, r, config.Global.Security.Username)
			RespRedirect(w, r, r.URL.String())
		} else {
			glog.Info("Login failed: [%s][%s]", username, password)
			RespRedirect(w, r, r.URL.String())
		}
	}
}

func secret(w http.ResponseWriter, r *http.Request) {
	head, tail := util.SplitPathRepeat(r.URL.Path, 1)
	glog.Debug("[%-4s][%-32s] [%s][%s]", r.Method, "/secret", head, tail)
	if !SessionVerify(r) {
		Login(w, r)
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

	if r.Method == http.MethodGet {
		secrets := db.GetSecrets()
		_ = tpl.Secret.Execute(w, map[string]any{"title": "Secret", "secrets": secrets})
		return
	}
}

func secretInsert(w http.ResponseWriter, r *http.Request) {
	head, tail := util.SplitPathRepeat(r.URL.Path, 1)
	glog.Debug("[%-4s][%-32s] [%s][%s]", r.Method, "/secret/insert", head, tail)
	if !SessionVerify(r) {
		Login(w, r)
		return
	}

	if r.Method == "GET" {
		_ = tpl.SecretInsert.Execute(w, map[string]any{"title": "Secret Insert", "default_time": time.Unix(0, 0).Format("2006-01-02T15:04")})
	} else if r.Method == "POST" {
		callerPst := r.PostFormValue("caller")
		if len(callerPst) < 1 {
			RespAPIInvalidInput(w)
			return
		}
		validityPeriodPst := r.PostFormValue("validity_period")
		validityPeriodTime, err := time.ParseInLocation("2006-01-02T15:04", validityPeriodPst, time.Local)
		if err != nil {
			RespAPIInvalidInput(w)
			return
		}
		data := new(db.SecretModel)
		data.Secret = GenerateSecret()
		data.Caller = callerPst
		data.ValidityPeriod = validityPeriodTime.Unix()
		res := db.InsertSecret(data)
		if !res {
			RespAPIProcessingFailed(w)
			return
		}
		LastPage(w, r)
	}
}

func secretModify(w http.ResponseWriter, r *http.Request) {
	head, tail := util.SplitPathRepeat(r.URL.Path, 1)
	glog.Debug("[%-4s][%-32s] [%s][%s]", r.Method, "/secret/modify", head, tail)
	if !SessionVerify(r) {
		Login(w, r)
		return
	}
	if len(tail) <= 1 {
		return
	}

	id, err := strconv.ParseInt(tail[1:], 10, 64)
	if err != nil {
		return
	}

	sec := db.GetSecretByID(id)
	if sec == nil {
		return
	}

	if r.Method == "GET" {
		validityPeriod := time.Unix(sec.ValidityPeriod, 0).Format("2006-01-02T15:04")
		createTime := time.Unix(sec.CreateTime, 0).Format("2006-01-02T15:04")
		lastUsed := time.Unix(sec.LastUsed, 0).Format("2006-01-02T15:04")
		expired := time.Unix(sec.Expired, 0).Format("2006-01-02T15:04")

		_ = tpl.SecretModify.Execute(w, map[string]any{"title": "Secret Modify", "id": sec.ID, "secret": sec.Secret,
			"caller": sec.Caller, "validity_period": validityPeriod, "create_time": createTime, "last_used": lastUsed, "expired": expired})
	} else if r.Method == "POST" {
		callerPst := r.PostFormValue("caller")
		if len(callerPst) < 1 {
			RespAPIInvalidInput(w)
			return
		}
		sec.Caller = callerPst
		validityPeriodPst := r.PostFormValue("validity_period")
		validityPeriodTime, err := time.ParseInLocation("2006-01-02T15:04", validityPeriodPst, time.Local)
		if err != nil {
			RespAPIInvalidInput(w)
			return
		}
		sec.ValidityPeriod = validityPeriodTime.Unix()
		expiredPst := r.PostFormValue("expired")
		expiredTime, err := time.ParseInLocation("2006-01-02T15:04", expiredPst, time.Local)
		if err != nil {
			RespAPIInvalidInput(w)
			return
		}
		sec.Expired = expiredTime.Unix()

		res := db.ModifySecret(sec)
		if !res {
			RespAPIProcessingFailed(w)
			return
		}
		Reload(w, r)
	}
}

func target(w http.ResponseWriter, r *http.Request) {
	head, tail := util.SplitPathRepeat(r.URL.Path, 1)
	glog.Debug("[%-4s][%-32s] [%s][%s]", r.Method, "/target", head, tail)
	if !SessionVerify(r) {
		Login(w, r)
		return
	}

	if r.Method == http.MethodGet {
		tg := db.GetTargets()
		glog.Debug("%#v", tg)
		targets := make([]string, 0, len(tg))
		Global.mutex.RLock()
		for _, v := range tg {
			m, ok := Global.table[v.TargetMod]
			glog.Debug("[target] mod:[%s] ok:[%v]", v.TargetMod, ok)
			if ok {
				t := m.GetTarget(v.TargetID)
				targets = append(targets, fmt.Sprintf("[%d][%s] %s", v.ID, t.GetKey(), t.GetName()))
			}
		}
		Global.mutex.RUnlock()

		_ = tpl.Target.Execute(w, map[string]any{"title": "Target", "targets": targets})
		return
	}
}

func send(w http.ResponseWriter, r *http.Request) {
	head, tail := util.SplitPathRepeat(r.URL.Path, 1)
	glog.Debug("[%-4s][%-32s] [%s][%s]", r.Method, "/send", head, tail)

	if r.Method == http.MethodPost {
		callerPst := strings.TrimSpace(r.FormValue("caller"))
		targetPst := strings.TrimSpace(r.FormValue("target"))
		secretPst := strings.TrimSpace(r.FormValue("secret"))
		dataPst := strings.TrimSpace(r.FormValue("data"))

		targetStr := strings.Split(targetPst, ",")
		for _, v := range targetStr {
			_, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err != nil {
				RespAPIInvalidInput(w)
				return
			}
		}

		sec := db.GetSecret(secretPst)
		if sec == nil {
			glog.Trace("failed to get secret [%s]", secretPst)
			RespAPIInvalidInput(w)
			return
		}

		glog.Debug("send: caller[%s] target[%s] secret[%s] data[%s]", callerPst, targetPst, secretPst, dataPst)
		msg := mod.ParseMessage([]byte(dataPst))
		if msg == nil {
			glog.Warning("parse message failed")
			RespAPIProcessingFailed(w)
			return
		}
		hid := db.InsertHistory(sec.ID, targetPst, dataPst, callerPst, GetIP(r))
		if hid < 0 {
			glog.Warning("insert history failed")
			RespAPIProcessingFailed(w)
			return
		}
		sendpkg.InsertByHistoryID(hid)

		RespAPIOk(w)
		return
	}
}
