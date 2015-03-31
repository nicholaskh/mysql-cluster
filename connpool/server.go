package connpool

import (
	"fmt"
	"time"

	"github.com/nicholaskh/golib/server"
	. "github.com/nicholaskh/mysql-cluster/config"
)

func LaunchServer() {
	s := server.NewTcpServer("connpool")

	s.LaunchTcpServer(Config.ListenAddr, handleClient, time.Minute*2, 200)
}

func handleClient(client *server.Client) {
	fmt.Println("connpool")
}
