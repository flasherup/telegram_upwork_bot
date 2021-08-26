package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type UpworkCfg struct {
	Url           string   `yaml:"url"`
	SecurityToken string   `yaml:"security_token"`
	UserUid       string   `yaml:"user_uid"`
	OrgUid        string   `yaml:"org_uid"`
	Feeds         []string `yaml:"feeds"`
}

type UpworkBotConfig struct {
	Token  string    `yaml:"token"`
	Upwork UpworkCfg `yaml:"upwork"`
}

func LoadConfig(path string) (config *UpworkBotConfig, err error) {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	config = &UpworkBotConfig{}

	err = yaml.Unmarshal(c, config)
	if err != nil {
		return
	}

	return
}
