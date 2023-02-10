package main

import (
	"flag"
	"fmt"
	"github.com/AkvicorEdwards/glog"
	"github.com/AkvicorEdwards/util"
	"msg-pusher/app"
	"msg-pusher/config"
	"msg-pusher/db"
	"msg-pusher/mod/wecom"
	"msg-pusher/send"
	_ "msg-pusher/tpl"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var err error

	isInit := flag.Bool("i", false, "init database")
	c := flag.String("c", "config.ini", "path to config file")
	flag.Parse()

	if util.FileStat(*c) != 2 {
		glog.Fatal("missing config [%s]!", *c)
	}

	config.Load(*c)

	glogMask := 0
	if config.Global.Log.MaskUnknown {
		glogMask |= glog.MaskUNKNOWN
	}
	if config.Global.Log.MaskDebug {
		glogMask |= glog.MaskDEBUG
	}
	if config.Global.Log.MaskTrace {
		glogMask |= glog.MaskTRACE
	}
	if config.Global.Log.MaskInfo {
		glogMask |= glog.MaskINFO
	}
	if config.Global.Log.MaskWarning {
		glogMask |= glog.MaskWARNING
	}
	if config.Global.Log.MaskError {
		glogMask |= glog.MaskERROR
	}
	if config.Global.Log.MaskFatal {
		glogMask |= glog.MaskFATAL
	}
	glog.SetMask(glogMask)
	if config.Global.Prod {
		glog.SetFlag(glog.FlagStdFlag)
	} else {
		glog.SetFlag(glog.FlagStdFlag | glog.FlagShortFile)
	}

	if *isInit {
		db.CreateDatabase()
		wecom.InsertDatabaseTable()
		os.Exit(0)
	}

	if util.FileStat(config.Global.Database.Path) != 2 {
		glog.Fatal("missing database [%s]!", config.Global.Database.Path)
	}

	if config.Global.Log.LogToFile {
		err = glog.SetLogFile(config.Global.Log.FilePath)
		if err != nil {
			glog.Fatal("failed to set log file [%s]", err.Error())
		}
	}

	EnableShutDownListener()

	app.Generate()
	ok := wecom.Load()
	if !ok {
		glog.Fatal("wecom load failed")
	}

	send.EnableServer()

	addr := fmt.Sprintf("%s:%d", config.Global.Server.HTTPAddr, config.Global.Server.HTTPPort)
	if config.Global.Server.EnableHTTPS {
		glog.Info("ListenAndServe: https://%s", addr)
		err = http.ListenAndServeTLS(addr, config.Global.Server.SSLCert, config.Global.Server.SSLKey, app.Global)
	} else {
		glog.Info("ListenAndServe: http://%s", addr)
		err = http.ListenAndServe(addr, app.Global)
	}
	if err != nil {
		glog.Fatal("failed to listen and serve [%s]", err.Error())
	}
}

func EnableShutDownListener() {
	go func() {
		down := make(chan os.Signal, 1)
		signal.Notify(down, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-down
		go func() {
			ticker := time.NewTicker(3 * time.Second)
			<-ticker.C
			glog.Fatal("Ticker Finished")
		}()

		glog.Info("close send server")
		send.KillServer()

		glog.Info("close log file")
		if config.Global.Log.LogToFile {
			glog.CloseFile()
		}
		glog.Info("log file closed")

		os.Exit(0)
	}()
}
