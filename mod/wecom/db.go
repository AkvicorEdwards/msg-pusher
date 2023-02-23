package wecom

import (
	"github.com/AkvicorEdwards/glog"
	"gorm.io/gorm"
	"msg-pusher/db"
	"sync"
	"time"
)

var dbLock = sync.RWMutex{}

type ModWecomModel struct {
	ID             int64  `gorm:"column:id;primaryKey;autoIncrement"`
	CorpID         string `gorm:"column:corp_id"`
	AgentID        int64  `gorm:"column:agent_id"`
	Secret         string `gorm:"column:secret"`
	Name           string `gorm:"column:name"`
	ValidityPeriod int64  `gorm:"column:validity_period"`
	CreateTime     int64  `gorm:"column:create_time;autoCreateTime"`
	LastUsed       int64  `gorm:"column:last_used"`
	Expired        int64  `gorm:"column:expired"`
}

func (ModWecomModel) TableName() string {
	return wecomTable
}

func (w *ModWecomModel) GetName() string {
	return w.Name
}

func (w *ModWecomModel) GetKey() string {
	return wecomKey
}

func (w *ModWecomModel) GetSecret() string {
	return w.Secret
}

func (w *ModWecomModel) GetValidityPeriod() int64 {
	return w.ValidityPeriod
}

func (w *ModWecomModel) GetCreateTime() int64 {
	return w.CreateTime
}

func (w *ModWecomModel) GetLastUsed() int64 {
	return w.LastUsed
}

func (w *ModWecomModel) GetExpired() int64 {
	return w.Expired
}

func InsertDatabaseTable() bool {
	d := db.Connect()
	if d == nil {
		glog.Fatal("can not connect database")
		return false
	}
	err := d.AutoMigrate(&ModWecomModel{})
	if err != nil {
		glog.Fatal("failed to create table [%s]", err.Error())
		return false
	}
	glog.Info("database table insert finished [%s]", ModWecomModel{}.TableName())
	return true
}

func getWecomByID(id int64, isUsed bool) *ModWecomModel {
	d := db.Connect()
	if d == nil {
		return nil
	}
	d = d.Model(&ModWecomModel{})
	dbLock.Lock()
	defer dbLock.Unlock()

	wecom := new(ModWecomModel)
	res := d.Where("id=?", id).First(wecom)
	if res.Error != nil || res.RowsAffected != 1 {
		return nil
	}

	if isUsed {
		now := time.Now().Unix()
		if now > wecom.ValidityPeriod {
			return nil
		}

		res = d.Where("id=?", id).Limit(1).Update("last_used", now)
		if res.Error != nil || res.RowsAffected != 1 {
			glog.Warning("failed to update last_used [%s]", res.Error)
			return nil
		}
	}

	return wecom
}

func getWecom() []ModWecomModel {
	d := db.Connect()
	if d == nil {
		return nil
	}
	d = d.Model(&ModWecomModel{})
	dbLock.RLock()
	defer dbLock.RUnlock()

	wecom := make([]ModWecomModel, 0)
	res := d.Where("expired=0").Find(&wecom)
	if res.Error != nil {
		glog.Warning("failed to get wecom [%s]", res.Error)
		return nil
	}
	return wecom
}

func insertWecom(data *ModWecomModel) bool {
	d := db.Connect()
	if d == nil {
		return false
	}
	dbLock.Lock()
	defer dbLock.Unlock()

	fun := func(d *gorm.DB) int64 {
		res := d.Model(&ModWecomModel{}).Create(data)
		if res.Error != nil || res.RowsAffected != 1 {
			glog.Warning("failed to insert wecom [%s]", res.Error)
			return -1
		}
		return data.ID
	}

	if db.InsertTarget(fun, data.TableName()) {
		return true
	}

	return false
}

func ModifyWecom(data *ModWecomModel) bool {
	d := db.Connect()
	if d == nil {
		return false
	}
	d = d.Model(&ModWecomModel{})
	dbLock.RLock()
	defer dbLock.RUnlock()

	res := d.Where("id=?", data.ID).Limit(1).Updates(data)
	if res.Error != nil || res.RowsAffected != 1 {
		glog.Warning("failed to update wecom [%v] [%v]", res.Error, res.RowsAffected)
		return false
	}

	return true
}
