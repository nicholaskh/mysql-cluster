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

}

func main() {
	s := server.NewTcpServer("mysql-cluster")
	s.LoadConfig(options.configFile)

	config.LoadConfig(s.Conf)

	InitGlobal()

	server.SetupLogging(options.logFile, options.logLevel, options.crashLogFile)

	LaunchServer()

	var ch chan int = make(chan int)
	<-ch

}
