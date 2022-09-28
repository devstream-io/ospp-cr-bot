package config

import (
	"fmt"
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/constants"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

var (
	Cfg           *Config
	LarkMaps      map[string]Maps
	MsgQueue      map[string]constants.MsgType
	TimeUnread1   time.Duration // PR/Issue 消息第一次发送消息后若未读，经过 TimeUnread1 后重发
	TimeUnread2   time.Duration // PR/Issue 消息第二次发送消息后若未读，经过 TimeUnread2 后发送给上级
	TimeUnread3   time.Duration // PR/Issue 消息第三次发送消息后若未读，经过 TimeUnread3 后抄送群聊
	CommentUnread time.Duration // Comment 消息第一次发送后若未读，经过 CommentUnread 后抄送群聊

)

type Config struct {
	Tasks     []Tasks   `yaml:"tasks"`
	Maps      []Maps    `yaml:"maps"`
	Scheduler Scheduler `yaml:"scheduler"`
}
type Tasks struct {
	Name         string `yaml:"name"`
	Repo         string `yaml:"repo"`
	Receiver     string `yaml:"receiver"`
	ReceiverType string `yaml:"receiverType"`
	PushChannel  string `yaml:"pushChannel"`
}

type Maps struct {
	Github string `yaml:"github"`
	Lark   string `yaml:"lark"`
	Role   string `yaml:"role"`
	Boss   string `yaml:"boss"`
}

type Scheduler struct {
	TimeUnread1   int `yaml:"timeUnread1"`
	TimeUnread2   int `yaml:"timeUnread2"`
	TimeUnread3   int `yaml:"timeUnread3"`
	CommentUnread int `yaml:"commentUnread"`
}

func InitConfig() {
	config, err := ioutil.ReadFile("./common.yaml")
	if err != nil {
		fmt.Print(err)
	}

	err = yaml.Unmarshal(config, &Cfg)
	if err != nil {
		fmt.Println("error")
	}

	if err != nil {
		fmt.Println("error")
	}

	LarkMaps = make(map[string]Maps)
	for _, i := range Cfg.Maps {
		LarkMaps[i.Github] = i
	}

	TimeUnread1 = time.Duration(Cfg.Scheduler.TimeUnread1) * time.Minute     // PR/Issue 消息第一次发送消息后若未读，经过 TimeUnread1 后重发
	TimeUnread2 = time.Duration(Cfg.Scheduler.TimeUnread2) * time.Minute     // PR/Issue 消息第二次发送消息后若未读，经过 TimeUnread2 后发送给上级
	TimeUnread3 = time.Duration(Cfg.Scheduler.TimeUnread3) * time.Minute     // PR/Issue 消息第三次发送消息后若未读，经过 TimeUnread3 后抄送群聊
	CommentUnread = time.Duration(Cfg.Scheduler.CommentUnread) * time.Minute // Comment 消息第一次发送后若未读，经过 CommentUnread 后抄送群聊

}

func QueryReceiveIDByRepo(repo string) (string, error) {
	for _, task := range Cfg.Tasks {
		if task.Repo == repo {
			return task.Receiver, nil
		}
	}
	return "", fmt.Errorf("repo not found")
}

func QueryGithubIdByLarkId(larkId string) (string, error) {
	for _, i := range Cfg.Maps {
		if i.Lark == larkId {
			return i.Github, nil
		}
	}
	return "", fmt.Errorf("larkId not found")
}
