package db

import (
	"fmt"
	"github.com/AkvicorEdwards/glog"
	"gorm.io/gorm"
	"sync"
	"time"
)

var historyLock = sync.RWMutex{}

type HistoryModel struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement"`
	SecretID int64  `gorm:"column:secret_id"`
	Targets  string `gorm:"column:targets"`
	Data     string `gorm:"column:data"`
	Caller   string `gorm:"column:caller"`
	IP       string `gorm:"column:ip"`
	Ready    int64  `gorm:"column:ready;autoCreateTime"`
	Sent     int64  `gorm:"column:sent"`
}

func (HistoryModel) TableName() string {
	return "history"
}

func (h *HistoryModel) Format() HistoryFormatModel {
	his := HistoryFormatModel{
		ID:       h.ID,
		SecretID: h.SecretID,
		Targets:  h.Targets,
		Data:     h.Data,
		Caller:   h.Caller,
		IP:       h.IP,
		Ready:    time.Unix(h.Ready, 0).Format("2006-01-02 15:04:05"),
	}
	if h.Sent == 0 {
		his.Sent = "Pending"
	} else if h.Sent < 0 {
		his.Sent = fmt.Sprintf("Discard: [%s]", time.Unix(h.Sent, 0).Format("2006-01-02 15:04:05"))
	} else {
		his.Sent = fmt.Sprintf("Sent: [%s]", time.Unix(h.Sent, 0).Format("2006-01-02 15:04:05"))
	}
	return his
}

type HistoryFormatModel struct {
	ID       int64
	SecretID int64
	Targets  string
	Data     string
	Caller   string
	IP       string
	Ready    string
	Sent     string
}

func GetHistory(id int64) *HistoryModel {
	d := Connect()
	if d == nil {
		return nil
	}
	d = d.Model(&HistoryModel{})
	historyLock.RLock()
	defer historyLock.RUnlock()

	history := new(HistoryModel)
	res := d.Where("id=?", id).First(history)
	if res.Error != nil || res.RowsAffected != 1 {
		glog.Warning("get history failed [%v] [%v]", res.Error, res.RowsAffected)
		return nil
	}
	return history
}

func GetAllHistories() []HistoryModel {
	d := Connect()
	if d == nil {
		return nil
	}
	d = d.Model(&HistoryModel{})
	historyLock.RLock()
	defer historyLock.RUnlock()

	histories := make([]HistoryModel, 0)
	res := d.Find(&histories)
	if res.Error != nil {
		glog.Warning("get all histories failed [%v] [%v]", res.Error, res.RowsAffected)
		return nil
	}
	return histories
}

func GetPendingHistories() []HistoryModel {
	d := Connect()
	if d == nil {
		return nil
	}
	d = d.Model(&HistoryModel{})
	historyLock.RLock()
	defer historyLock.RUnlock()

	histories := make([]HistoryModel, 0)
	res := d.Where("sent=0").Find(&histories)
	if res.Error != nil {
		glog.Warning("get pending histories failed [%v] [%v]", res.Error, res.RowsAffected)
		return nil
	}
	return histories
}

func InsertHistory(secretID int64, targets string, data string, caller string, ip string) int64 {
	d := Connect()
	if d == nil {
		return -1
	}
	d = d.Model(&HistoryModel{})
	historyLock.RLock()
	defer historyLock.RUnlock()

	now := time.Now().Unix()
	his := &HistoryModel{
		SecretID: secretID,
		Targets:  targets,
		Data:     data,
		Caller:   caller,
		IP:       ip,
		Ready:    now,
		Sent:     0,
	}
	res := d.Create(his)
	if res.Error != nil || res.RowsAffected != 1 {
		glog.Warning("insert history failed [%v] [%v]", res.Error, res.RowsAffected)
		return -1
	}
	return his.ID
}

func UpdateHistory(id int64, enforce, discard bool) {
	d := Connect()
	if d == nil {
		return
	}
	d = d.Model(&HistoryModel{})
	historyLock.Lock()
	defer historyLock.Unlock()

	var f int64 = 1
	if discard {
		f = -1
	}

	var res *gorm.DB
	if enforce {
		res = d.Where("id=?", id).Update("sent", f*time.Now().Unix())
	} else {
		res = d.Where("id=? AND sent=0", id).Update("sent", f*time.Now().Unix())
	}
	if res.Error != nil || res.RowsAffected != 1 {
		glog.Warning("update history failed [%d] [%v] [%v]", id, res.Error, res.RowsAffected)
		return
	}
}
