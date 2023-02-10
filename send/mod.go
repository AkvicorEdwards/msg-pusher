package send

import (
	"msg-pusher/mod"
	"sync"
)

var targetMod = make(map[string]mod.Model)
var targetModLock = sync.RWMutex{}

func InsertMod(mod mod.Model) bool {
	targetModLock.Lock()
	defer targetModLock.Unlock()

	_, ok := targetMod[mod.Table()]
	if ok {
		return false
	}
	targetMod[mod.Table()] = mod
	return true
}
