package config

import (
	"testing"
)

func TestConfigInitFromYaml(t *testing.T) {

	ConfigInitFromYaml("/home/chengze/go/src/github.com/multivactech/monitor/config/config.yaml")
	// if Config == nil {
	// 	t.Error("Config is nil")
	// 	return
	// }
	// fmt.Print(Config)
	// fmt.Print(Config.LocalDir)
	// fmt.Print(Config.MailConfig)
	// t.Log(Config)
	// fmt.Print(Config.SshNodes)

}
