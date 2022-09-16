package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/config"
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/constants"
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
	config.MsgQueue = make(map[string]constants.MsgType)
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
	router.POST("/api/lark/cardCallback", lark.CardCallback)

	err := router.Run(":3000")
	if err != nil {
		return
	}
	return
}
