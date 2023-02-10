package send

import (
	"github.com/AkvicorEdwards/glog"
	"msg-pusher/db"
	"msg-pusher/mod"
	"strconv"
	"strings"
)

var queue *queueModel

func init() {
	queue = new(queueModel)
	queue.data = make([]*queueItem, 128, 128)
	for k, _ := range queue.data {
		queue.data[k] = nil
	}
	queue.index = 0
	queue.updated = make(chan bool, 3)
	queue.killed = make(chan bool, 1)
	queue.dead = false
}

func InsertByHistoryID(historyID int64) bool {
	history := db.GetHistory(historyID)
	if history == nil {
		return false
	}
	return InsertByHistory(history)
}

func InsertByHistory(history *db.HistoryModel) bool {
	targetStr := strings.Split(history.Targets, ",")
	targetsID := make([]int64, len(targetStr))
	for k, v := range targetStr {
		id, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return false
		}
		targetsID[k] = id
	}
	msg := mod.ParseMessage([]byte(history.Data))
	if msg == nil {
		return false
	}

	targets := db.GetTargetIndex(targetsID)
	glog.Debug("targets req[%d] ok[%d] : req[%#v] ok[%#v]", len(targetsID), len(targets), targetsID, targets)

	targetModLock.RLock()
	defer targetModLock.RUnlock()
	for _, v := range targets {
		m, ok := targetMod[v.TargetMod]
		glog.Debug("[send] mod:[%s] ok:[%v]", v.TargetMod, ok)
		if ok {
			queue.push(m.Prepare(v.TargetID, msg), msg, history)
		}
	}
	return true
}
