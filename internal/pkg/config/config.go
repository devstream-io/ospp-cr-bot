package config

import (
	"fmt"
	"github.com/lyleshaw/ospp-cr-bot/internal/pkg/constants"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	Cfg      *Config
	LarkMaps map[string]Maps
	MsgQueue map[string]constants.MsgType
)

type Config struct {
	Tasks []Tasks `yaml:"tasks"`
	Maps  []Maps  `yaml:"maps"`
}
type Tasks struct {
	Name         string `yaml:"name"`
	Repo         string `yaml:"repo"`
	RepoType     string `yaml:"repoType"`
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
