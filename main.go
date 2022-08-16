package main

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/event"
	"altsub/modules/source"
	"altsub/server"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func watchingSignal() {
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("Program Exit...", s)
				base.NewLog("warn", nil, fmt.Sprintf("Program exit with code: %d", s), "watchingSignal()")
				graceFullExit(99)
			case syscall.SIGUSR1:
				fmt.Println("usr1 signal", s)
			case syscall.SIGUSR2:
				fmt.Println("usr2 signal", s)
			default:
				fmt.Println("other signal", s)
			}
		}
	}()
}

func graceFullExit(code int) {
	os.Exit(code)
}

func parseFlags() {
	var cfgFile = flag.String("config", "./configs/config.toml", "server config file")

	flag.Parse()
	base.ConfigFile = *cfgFile
}

func main() {
	parseFlags()
	base.Init()
	base.MigrateDB(&models.MSource{}, &models.MSchema{}, &models.MEvent{}, &models.MRule{}, &models.MRuleClause{}, &models.MReceiver{}, &models.MSubscribe{}, &models.MSchemaedEvent{}, &models.MDuty{}, &models.MDutyGroup{}, &models.MDutyAt{})
	if err := base.DB().SetupJoinTable(&models.MRule{}, "Receivers", &models.MSubscribe{}); err != nil {
		base.NewLog("fatal", err, "db join table setup error", "main()")
	}
	if err := base.DB().SetupJoinTable(&models.MReceiver{}, "Rules", &models.MSubscribe{}); err != nil {
		base.NewLog("fatal", err, "db join table setup error", "main()")
	}
	ss, err := source.Fetch(base.DB(), nil)
	if err != nil {
		base.NewLog("fatal", err, "couldn't fetch sources from mysql", "main()")
	}
	var sources []string
	for _, s := range ss.All {
		sources = append(sources, s.Name)
	}
	base.InitKafka(sources...)
	event.ReadAndParseEventFromBufferForever(sources...)
	watchingSignal()
	var httpServer = server.NewServer(base.Config.MainConfig.Listen, gin.DebugMode)
	httpServer.RegisterRoutes()
	httpServer.RunForever()
}
