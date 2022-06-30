package core

import (
	"GoProject/global"
	"GoProject/util"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func InitLogger() {
	global.Logger = logrus.New()
	global.Logger.ReportCaller = true

	level, err := logrus.ParseLevel(global.Level)
	if err != nil {
		global.Logger.WithFields(logrus.Fields{
			"level":        global.Level,
			"allLevels":    logrus.AllLevels,
			"defaultLevel": "warn",
			"err":          err,
		}).Warn("parse log level failed, will use default level")
	} else {
		global.Logger.SetLevel(level)
	}

	if global.Mode == "product" {
		global.Logger.SetFormatter(&logrus.JSONFormatter{})
		logPath := filepath.Join(util.GetExecutePath(), global.LogPath)
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			global.Logger.SetOutput(os.Stdout)
			global.Logger.WithFields(logrus.Fields{
				"path": logPath,
				"err":  err,
			}).Warn("set log out failed, will use default out")
		} else {
			global.Logger.SetOutput(file)
			global.Logger.WithFields(logrus.Fields{
				"path": logPath,
			}).Info("set log out success")
		}
	} else {
		global.Logger.SetFormatter(&logrus.TextFormatter{})
		global.Logger.SetOutput(os.Stdout)
	}
}
