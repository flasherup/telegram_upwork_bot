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

type User struct {
	Id int64 `yaml:"id"`
	ZoneName string `yaml:"zone_name"`
	ZoneOffset int `yaml:"zone_offset"`
}

type Telegram struct {
	Token string  `yaml:"token"`
	Users []User `yaml:"users"`
}

type Filters struct {
	SkipDuration int   `yaml:"skip_duration"`
	Skills       []string `yaml:"skills"`
	Countries    []string `yaml:"countries"`
}

type UpworkBotConfig struct {
	Telegram Telegram  `yaml:"telegram"`
	Upwork   UpworkCfg `yaml:"upwork"`
	Filters  Filters   `yaml:"filters"`
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
