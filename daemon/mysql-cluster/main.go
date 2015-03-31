package main

import (
	"github.com/nicholaskh/golib/server"
	"github.com/nicholaskh/mysql-cluster/config"
	"github.com/nicholaskh/mysql-cluster/connpool"
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

	connpool.LaunchServer()

	var ch chan int = make(chan int)
	<-ch

}
