package main

import (
	"github.com/multivactech/monitor/logger"

	"github.com/multivactech/monitor/config"
)

func main() {

	configDir := "/home/chengze/go/src/github.com/multivactech/monitor/config/config.yaml"

	logger.SetLogDebug()
	
	
	logger.InitLogRotator("log/output.log")
	config.ConfigInitFromYaml(configDir)

	logger.Log.Debug(config.Config)

	// var monitor monitor.Monitor

	// os.MkdirAll("log", os.ModePerm)
	// logFile, err := os.OpenFile("log/output.log", os.O_RDWR|os.O_CREATE, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer logFile.Close()
	// log.SetOutput(logFile)

	// monitor.Start()

}
