package main

import (
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/eventListener"
	"github.com/lyleshaw/ospp-cr-bot/pkg/utils/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	isDebug bool
)

func init() {
	initLog()
	initConfig()
}
func initConfig() {
	viper.AutomaticEnv()
	if err := viper.BindEnv("lark_access_token"); err != nil {
		log.Fatal(err)
	}
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
	err := eventListener.Init()
	if err != nil {
		log.Fatal(err.Error())
	}
}
