package core

import (
	"GoProject/database"
	"GoProject/global"
	"GoProject/util"
	"embed"
	"github.com/sirupsen/logrus"
	"net/http"
)

func InitServer(sources *embed.FS) {
	util.InitCache()
	database.InitDb()
	r := InitRoutes(sources)

	s := &http.Server{
		Addr:           global.WebServeAddr,
		Handler:        r,
		ReadTimeout:    global.DefaultEventTimeout,
		WriteTimeout:   global.DefaultEventTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if err != nil {
		global.Logger.WithFields(logrus.Fields{
			"err":  err.Error(),
			"port": global.WebServeAddr,
		}).Panic("start listen panic")
	}
}
