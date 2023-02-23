package db

import (
	"github.com/AkvicorEdwards/glog"
	"sync"
	"time"
)

var secretLock = sync.RWMutex{}

type SecretModel struct {
	ID             int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Secret         string `gorm:"column:secret;unique"`
	Caller         string `gorm:"column:caller"`
	ValidityPeriod int64  `gorm:"column:validity_period"`
	CreateTime     int64  `gorm:"column:create_time;autoCreateTime"`
	LastUsed       int64  `gorm:"column:last_used"`
	Expired        int64  `gorm:"column:expired"`
}

func (SecretModel) TableName() string {
	return "secret"
}

func InsertSecret(data *SecretModel) bool {
	d := Connect()
	if d == nil {
		return false
	}
	d = d.Model(&SecretModel{})
	secretLock.Lock()
	defer secretLock.Unlock()

	res := d.Create(data)
	if res.Error != nil || res.RowsAffected != 1 {
		glog.Warning("insert secret failed [%v] [%v]", res.Error, res.RowsAffected)
		return false
	}

	return true
}

func ModifySecret(data *SecretModel) bool {
	d := Connect()
	if d == nil {
		return false
	}
	d = d.Model(&SecretModel{})
	secretLock.Lock()
	defer secretLock.Unlock()

	res := d.Where("id=?", data.ID).Limit(1).Updates(data)
	if res.Error != nil || res.RowsAffected != 1 {
		glog.Warning("failed to update secret [%v] [%v]", res.Error, res.RowsAffected)
		return false
	}

	return true
}

func GetSecrets() []SecretModel {
	d := Connect()
	if d == nil {
		return nil
	}
	d = d.Model(&SecretModel{})
	secretLock.RLock()
	defer secretLock.RUnlock()
	secret := make([]SecretModel, 0)
	res := d.Where("expired=0").Find(&secret)
	if res.Error != nil {
		glog.Warning("get secrets failed [%v] [%v]", res.Error, res.RowsAffected)
		return nil
	}
	return secret
}

func GetSecret(secret string) *SecretModel {
	d := Connect()
	if d == nil {
		return nil
	}
	d = d.Model(&SecretModel{})
	secretLock.Lock()
	defer secretLock.Unlock()

	sec := new(SecretModel)
	res := d.Where("secret=?", secret).First(sec)
	if res.Error != nil || res.RowsAffected != 1 {
		return nil
	}

	now := time.Now().Unix()
	if now > sec.ValidityPeriod {
		return nil
	}

	res = d.Where("secret=?", secret).Limit(1).Update("last_used", now)
	if res.Error != nil || res.RowsAffected != 1 {
		glog.Warning("failed to update last_used [%v] [%v]", res.Error, res.RowsAffected)
		return nil
	}

	return sec
}

func GetSecretByID(id int64) *SecretModel {
	d := Connect()
	if d == nil {
		return nil
	}
	d = d.Model(&SecretModel{})
	secretLock.Lock()
	defer secretLock.Unlock()

	sec := new(SecretModel)
	res := d.Where("id=?", id).First(sec)
	if res.Error != nil || res.RowsAffected != 1 {
		return nil
	}

	return sec
}
