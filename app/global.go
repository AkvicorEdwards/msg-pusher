package app

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"msg-pusher/config"
	"msg-pusher/mod"
	"net/http"
	"sync"
)

var Global *app

func Generate() {
	Global = new(app)
	Global.session = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))
	Global.session.Options.Domain = config.Global.Session.Domain
	Global.session.Options.Path = config.Global.Session.Path
	Global.session.Options.MaxAge = config.Global.Session.MaxAge
	Global.handler = make(map[string]func(w http.ResponseWriter, r *http.Request))
	Global.url = make(map[string]mod.Model)
	Global.table = make(map[string]mod.Model)
	Global.key = make(map[string]mod.Model)
	Global.mutex = sync.RWMutex{}

	Global.handler["/favicon.ico"] = staticFavicon

	Global.handler["/"] = index
	Global.handler["/login"] = Login
	Global.handler["/secret"] = secret
	Global.handler["/target"] = target
	Global.handler["/history"] = history
	Global.handler["/send"] = send
}
