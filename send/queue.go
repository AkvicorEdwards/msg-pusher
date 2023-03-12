package send

import (
	"github.com/AkvicorEdwards/glog"
	"msg-pusher/db"
	"msg-pusher/mod"
	"sync"
	"time"
)

type queueItem struct {
	Package mod.Package
	Message *mod.MessageModel
	History *db.HistoryModel
}

type queueModel struct {
	data    []*queueItem
	index   int
	updated chan bool
	killed  chan bool
	dead    bool
	sync.RWMutex
}

func (q *queueModel) push(pkg mod.Package, msg *mod.MessageModel, history *db.HistoryModel) {
	defer func() {
		q.updated <- true
	}()
	q.Lock()
	defer q.Unlock()
	item := new(queueItem)
	item.Package = pkg
	item.Message = msg
	item.History = history
	defer func() {
		q.index++
	}()

	if q.index >= len(q.data) {
		q.data = append(q.data, nil)
	}

	if q.index == 0 {
		q.data[0] = item
		return
	}

	i := q.index
	for ; i > 0; i-- {
		next := q.data[i-1]
		if next.Message.TimeSend > msg.TimeSend {
			break
		} else if next.Message.TimeSend == msg.TimeSend {
			if next.Message.Urgency > msg.Urgency {
				break
			}
		}
		q.data[i] = next
	}
	q.data[i] = item
}

// pop() returns the message to be sent within 30 seconds
//    if not, returns the wait time for the next message to send minus 15 seconds
func (q *queueModel) pop() (*queueItem, time.Duration) {
	if q.index <= 0 {
		return nil, 10 * time.Minute
	}
	q.Lock()
	defer q.Unlock()
	var item = q.data[q.index-1]
	limit := time.Now().Unix() + 30
	if item.Message.TimeSend > limit {
		glog.Trace("pop: need[%d] now[%d] x[%d] ret[%d]", item.Message.TimeSend, limit-30, item.Message.TimeSend-limit+30, item.Message.TimeSend-limit+15)
		return nil, time.Duration(item.Message.TimeSend-limit+15) * time.Second
	}
	q.index--
	q.data[q.index] = nil
	return item, 0
}

// EnableServer must be executed after all modules are loaded
func EnableServer() {
	go service()
	histories := db.GetPendingHistories()
	if histories != nil {
		for _, v := range histories {
			h := v
			InsertByHistory(&h)
		}
	}
}

func KillServer() {
	queue.killed <- true
}

func service() {
	defer func() {
		glog.Warning("send server killed")
	}()
	var duration = trySend()
	for {
		glog.Trace("next ticker wait [%v]", duration)
		select {
		case <-queue.killed:
			glog.Trace("killed")
			queue.dead = true
			return
		case <-queue.updated:
			glog.Trace("updated")
			duration = trySend()
		case <-time.After(duration):
			glog.Trace("ticker")
			duration = trySend()
		}
	}
}

func trySend() time.Duration {
	item, wait := queue.pop()
	for item != nil {
		glog.Trace("try send [%s]", item.Message.String())
		go waitSend(item)
		item, wait = queue.pop()
	}
	return wait
}

func waitSend(item *queueItem) {
	wait := func(send int64) time.Duration {
		send -= time.Now().Unix()
		if send < 0 {
			send = 0
		}
		return time.Duration(send) * time.Second
	}
	var ok bool
	limit := 1
	w := wait(item.Message.TimeSend)
	glog.Trace("send wait [%v]", w)

	select {
	case <-time.After(w):
		if queue.dead {
			glog.Trace("dead [%s]", item.Message.String())
			return
		}
		glog.Trace("send [%s]", item.Message.String())
		for limit <= repeatLimit {
			ok = item.Package.Send()
			if ok {
				glog.Trace("sent [%s]", item.Message.String())
				db.UpdateHistory(item.History.ID, false, false)
				return
			}
			limit *= 2
			time.Sleep(time.Duration(limit) * time.Second)
		}
		db.UpdateHistory(item.History.ID, false, true)
		glog.Trace("send failed [%s]", item.Message.String())
	}
}
