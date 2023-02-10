package wecom

import (
	"msg-pusher/app"
	"net/http"
)

var self *Model

func Load() bool {
	self = new(Model)
	self.handler = make(map[string]func(w http.ResponseWriter, r *http.Request))
	insertHandler()
	return app.InsertMod(self)
}

func insertHandler() {
	self.handler["/"] = index
	self.handler["/secret"] = secret
}
