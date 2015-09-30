package main

import (
	"github.com/nicholaskh/golib/server"
	"github.com/nicholaskh/mysql-cluster/config"
	. "github.com/nicholaskh/mysql-cluster/proxygate"
)

func init() {
	parseFlags()

	if options.showVersion {
		server.ShowVersionAndExit()
	}

	conf := server.LoadConfig(options.configFile)
	config.LoadConfig(conf)

}

func main() {
	InitGlobal()

	server.SetupLogging(options.logFile, options.logLevel, options.crashLogFile)

	LaunchServer()

	var ch chan int = make(chan int)
	<-ch

}
