package config

import (
	"fmt"

	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	Cfg *Config
)

type Config struct {
	Tasks []Tasks `yaml:"tasks"`
}
type Tasks struct {
	Name         string `yaml:"name"`
	Repo         string `yaml:"repo"`
	RepoType     string `yaml:"repoType"`
	Recevier     string `yaml:"recevier"`
	RecevierType string `yaml:"recevierType"`
	PushChannel  string `yaml:"pushChannel"`
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
}

func QueryReceiveIDByRepo(repo string) (string, error) {
	for _, task := range Cfg.Tasks {
		if task.Repo == repo {
			return task.Recevier, nil
		}
	}
	return "", fmt.Errorf("repo not found")
}
