package core

import (
	"GoProject/global"
	"embed"
	"github.com/kardianos/service"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type program struct {
	Fs *embed.FS
}

type ServiceControl struct {
	service service.Service
}

func (sc *ServiceControl) Control(action string) {
	err := service.Control(sc.service, action)
	if err != nil {
		global.Logger.WithFields(logrus.Fields{
			"action": action,
			"err":    err,
		}).Error("service run action error")
		return
	}
	if action == "install" {
		err = service.Control(sc.service, "start")
		if err != nil {
			global.Logger.WithFields(logrus.Fields{
				"action": "start",
				"err":    err,
			}).Error("service start error")
		}
	}
}

var serviceConfig = &service.Config{
	Name:        "goproject",
	DisplayName: "go project service",
	Description: "this is a go web project service",
}

func InitService(fs *embed.FS) (*ServiceControl, error) {
	p := &program{fs}
	s, err := service.New(p, serviceConfig)
	if err != nil {
		global.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("init service Error")
		return nil, err
	}
	return &ServiceControl{s}, nil
}

var eg errgroup.Group

func (p *program) Start(s service.Service) error {
	global.Logger.Info("start service")
	go p.Run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	global.Logger.Info("stop service")
	return nil
}

func (p *program) Run() {
	eg.Go(func() error {
		return InitServer(p.Fs)
	})

	if err := eg.Wait(); err != nil {
		global.Logger.Error("启动系统出错")
	}
}
