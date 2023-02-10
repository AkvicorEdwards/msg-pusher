package app

import (
	"github.com/AkvicorEdwards/glog"
	"msg-pusher/mod"
	sendpkg "msg-pusher/send"
)

func InsertMod(mod mod.Model) bool {
	Global.mutex.Lock()
	defer Global.mutex.Unlock()
	_, ok := Global.url[mod.URL()]
	if ok {
		glog.Warning("insert Global.url failed [%s][%s]", mod.URL(), mod.Name())
		return false
	}
	_, ok = Global.table[mod.Table()]
	if ok {
		glog.Warning("insert Global.table failed [%s][%s]", mod.Table(), mod.Name())
		return false
	}
	_, ok = Global.key[mod.Key()]
	if ok {
		glog.Warning("insert Global.key failed [%s][%s]", mod.Key(), mod.Name())
		return false
	}
	if !sendpkg.InsertMod(mod) {
		glog.Warning("insert send.targetMod failed [%s][%s]", mod.Table(), mod.Name())
		return false
	}
	Global.url[mod.URL()] = mod
	Global.table[mod.Table()] = mod
	Global.key[mod.Key()] = mod
	glog.Info("insert mod successful [%s] url[%s] key[%s] table[%s]", mod.Name(), mod.URL(), mod.Key(), mod.Table())
	return true
}
