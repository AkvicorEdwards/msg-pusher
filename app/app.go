package app

import (
	"msg-pusher/mod"
	"net/http"
	"sync"

	_ "msg-pusher/tpl"

	"github.com/AkvicorEdwards/glog"
	"github.com/AkvicorEdwards/util"
	"github.com/gorilla/sessions"
)

type app struct {
	mutex   sync.RWMutex
	session *sessions.CookieStore
	handler map[string]func(w http.ResponseWriter, r *http.Request)
	url     map[string]mod.Model
	table   map[string]mod.Model
	// Only used as a set to prevent duplication
	key map[string]mod.Model
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	head, tail := util.SplitPath(r.URL.Path)
	glog.Debug("[%-4s][%-32s] [%s][%s]", r.Method, "/", head, tail)

	var handler func(w http.ResponseWriter, r *http.Request)
	var m mod.Model
	var ok bool

	handler, ok = a.handler[head]
	if ok {
		handler(w, r)
		return
	}

	a.mutex.RLock()
	m, ok = a.url[head]
	a.mutex.RUnlock()
	if ok {
		m.Handle(w, r)
		return
	}
	glog.Debug("Unhandled [%-4s][%-32s] [%s][%s]", r.Method, "/", head, tail)
}
