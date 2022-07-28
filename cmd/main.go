package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/config"
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/eventListener"
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/pushChannel/lark"
	"github.com/lyleshaw/ospp-cr-bot/pkg/utils/log"
	"github.com/sirupsen/logrus"
)

var (
	isDebug bool
)

func init() {
	initLog()
	config.InitConfig()
}

func initLog() {
	if isDebug {
		logrus.SetLevel(logrus.DebugLevel)
		log.Infof("Log level is: %s.", logrus.GetLevel())
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	router := gin.Default()
	router.POST("/github/webhook", eventListener.GitHubWebHook)
	router.POST("/api/lark/callback", lark.Callback)
	err := router.Run(":3000")
	if err != nil {
		return
	}
	return
}
