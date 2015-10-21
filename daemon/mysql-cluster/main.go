package main

import (
	"github.com/nicholaskh/golib/server"
	"github.com/nicholaskh/mysql-cluster/config"
	"github.com/nicholaskh/mysql-cluster/core"
	server1 "github.com/nicholaskh/mysql-cluster/server"
)

func init() {
	parseFlags()

	if options.showVersion {
		server.ShowVersionAndExit()
	}

	conf := server.LoadConfig(options.configFile)
	config.LoadConfig(conf)

	core.MysqlClusterInstance = core.NewMysqlCluster()
}

func main() {
	server.SetupLogging(options.logFile, options.logLevel, options.crashLogFile)

	server1.LaunchServer()

	var ch chan int = make(chan int)
	<-ch

}
