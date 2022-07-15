package main

import (
	"GoProject/core"
	"GoProject/global"
	"embed"
	"flag"
)

//go:embed frontend/dist/*
var sources embed.FS

func init() {
	flag.StringVar(&global.Mode, "mode", "product", "application run mode")
	flag.BoolVar(&global.Test, "test", false, "application deployment mode")
	flag.StringVar(&global.Level, "level", "warn", "application log level")
	flag.StringVar(&global.Action, "action", "start", "application service action")
	flag.Parse()

	core.InitLogger()
}

func main() {
	switch global.Mode {
	case "product":
		service, err := core.InitService(&sources)
		if err != nil {
			global.Logger.Panicln("init service panic")
		}
		if global.Action != "" {
			service.Control(global.Action)
		} else {
			core.InitServer(&sources)
		}
	default:
		core.InitServer(&sources)
	}
}
