package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var Config *configStyle

type configStyle struct {
	RemoteDir      string
	LocalDir       string
	SshNodes       []SshNode
	MailConfig     MailConfig
	ErrorWhiteList []string
	// mutex          *sync.Mutex
}

type SshNode struct {
	Host   string
	Passwd string
}

type MailConfig struct {
	Sender   string
	Passwd   string
	Receiver string
}

func ConfigInitFromYaml(configDir string) {
	data, err := ioutil.ReadFile(configDir)

	if err != nil {
		panic(fmt.Errorf("read config.yaml failed. %v", err))
	}

	Config = &configStyle{}
	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		panic(fmt.Errorf("parse config.yaml failed. %v", err))
	}
	// log.Print(Config)
}

func ConfigInitFromJson(configDir string) {

	data, err := ioutil.ReadFile(configDir)
	if err != nil {
		panic(fmt.Errorf("read config.json failed. %v", err))
	}
	Config = &configStyle{}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		panic(fmt.Errorf("parse config.json failed. %v", err))
	}
}
