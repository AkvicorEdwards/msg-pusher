package db

import (
	"github.com/AkvicorEdwards/glog"
	"github.com/AkvicorEdwards/util"
	"msg-pusher/config"
)

func CreateDatabase() {
	if util.FileStat(config.Global.Database.Path) == 2 {
		glog.Fatal("database file exist!")
	}
	d := Connect()
	if d == nil {
		glog.Fatal("con not connect to database!")
	}
	err := db.AutoMigrate(&SecretModel{}, &TargetModel{}, &HistoryModel{})
	if err != nil {
		glog.Fatal(err.Error())
	}
	glog.Info("database create finished")
}
