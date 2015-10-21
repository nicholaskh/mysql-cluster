package main

import (
	"github.com/nicholaskh/golib/server"
	"github.com/nicholaskh/mysql-cluster/config"
	"github.com/nicholaskh/mysql-cluster/core"
	server1 "github.com/nicholaskh/mysql-cluster/server"
)

var mycConfig *config.MycConfig

func init() {
	parseFlags()

	if options.showVersion {
		server.ShowVersionAndExit()
	}

	conf := server.LoadConfig(options.configFile)
	mycConfig = new(config.MycConfig)
	mycConfig.LoadConfig(conf)

	core.MysqlClusterInstance = core.NewMysqlCluster(mycConfig)
}

func main() {
	server.SetupLogging(options.logFile, options.logLevel, options.crashLogFile)

	server1.LaunchServer(mycConfig.Server)

	var ch chan int = make(chan int)
	<-ch

}
