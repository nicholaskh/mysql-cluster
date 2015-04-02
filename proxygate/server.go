package proxygate

import (
	"fmt"
	"time"

	"github.com/nicholaskh/golib/server"
	log "github.com/nicholaskh/log4go"
	. "github.com/nicholaskh/mysql-cluster/config"
	proto "github.com/nicholaskh/mysql-cluster/proto/go"
)

func LaunchServer() {
	s := server.NewTcpServer("connpool")

	s.LaunchTcpServer(Config.ListenAddr, newClientHandler(), time.Minute*2, 200)
}

type ClientHandler struct {
	client *server.Client
}

func newClientHandler() *ClientHandler {
	return &ClientHandler{}
}

func (this *ClientHandler) OnAccept(cli *server.Client) {
	this.client = cli
}

func (this *ClientHandler) OnRead(input string) {
	q := proto.NewQuery()
	err := q.Parse([]byte(input), len(input))
	if err != nil {
		log.Error("parse query error")
	}
	log.Info("sql: %s\npool: %s", q.Getsql(), q.Getpool())
	this.client.WriteMsg(fmt.Sprintf("received sql: %s\nreceived pool: %s\n", q.Getsql(), q.Getpool()))
}

func (this *ClientHandler) OnClose() {
	this.client.Close()
}

func handleClient(client *server.Client) {
	client.Conn.Write([]byte("connected"))
}
