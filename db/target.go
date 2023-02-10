package db

import (
	"github.com/AkvicorEdwards/glog"
	"gorm.io/gorm"
	"sync"
)

var targetLock = sync.RWMutex{}

type TargetModel struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement"`
	TargetMod string `gorm:"column:target_mod"`
	TargetID  int64  `gorm:"column:target_id"`
}

func (TargetModel) TableName() string {
	return "target"
}

func GetTargetIndex(ids []int64) []TargetModel {
	d := Connect()
	if d == nil {
		return nil
	}
	d = d.Model(&TargetModel{})
	targetLock.RLock()
	defer targetLock.RUnlock()

	sec := make([]TargetModel, 0, len(ids))
	res := d.Where("id IN ?", ids).Find(&sec)
	if res.Error != nil {
		glog.Warning("get target failed [%v] [%v]", res.Error, res.RowsAffected)
		return nil
	}

	return sec
}

func GetTargets() []TargetModel {
	d := Connect()
	if d == nil {
		return nil
	}
	d = d.Model(&TargetModel{})
	targetLock.RLock()
	defer targetLock.RUnlock()

	sec := make([]TargetModel, 0)
	res := d.Find(&sec)
	if res.Error != nil {
		glog.Warning("get targets failed [%v] [%v]", res.Error, res.RowsAffected)
		return nil
	}

	return sec
}

// TargetFunc returning a negative number means that the function execution failed.
// when the function executes successfully, it returns the id of the table where the record is located
type TargetFunc func(d *gorm.DB) int64

func InsertTarget(fun TargetFunc, table string) bool {
	d := Connect()
	if d == nil {
		return false
	}
	targetLock.Lock()
	defer targetLock.Unlock()
	dx := d.Begin()
	id := fun(dx)
	if id < 0 {
		dx.Rollback()
		return false
	}

	res := dx.Model(&TargetModel{}).Create(&TargetModel{TargetMod: table, TargetID: id})
	if res.Error != nil || res.RowsAffected != 1 {
		dx.Rollback()
		glog.Warning("failed to insert target [%s]", res.Error)
		return false
	}
	dx.Commit()
	return true
}
