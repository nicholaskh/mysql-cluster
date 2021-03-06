package server

import (
	"fmt"
	"io"
	"net"
	"time"

	proto1 "github.com/golang/protobuf/proto"
	"github.com/nicholaskh/golib/server"
	log "github.com/nicholaskh/log4go"
	. "github.com/nicholaskh/mysql-cluster/config"
	"github.com/nicholaskh/mysql-cluster/core"
	proto "github.com/nicholaskh/mysql-cluster/proto/go"
)

func LaunchServer(config *ServerConfig) {
	s := server.NewTcpServer("myc")

	s.LaunchTcpServer(config.ListenAddr, newClientHandler(config), time.Minute*2, 200)
}

type ClientHandler struct {
	serverConfig *ServerConfig
}

func newClientHandler(serverConfig *ServerConfig) *ClientHandler {
	return &ClientHandler{serverConfig}
}

func (this *ClientHandler) OnAccept(c *server.Client) {
	proto := server.NewProtocol()
	proto.SetConn(c.Conn)
	client := newClient(c, proto)
	for {
		if this.serverConfig.SessTimeout.Nanoseconds() > int64(0) {
			client.Proto.SetReadDeadline(time.Now().Add(this.serverConfig.SessTimeout))
		}

		data, err := client.proto.Read()

		if err != nil {
			err_, ok := err.(net.Error)
			if ok {
				if err_.Temporary() {
					log.Info("Temporary failure: %s", err_.Error())
					break
				}
			}
			if err == io.EOF {
				log.Info("Client %s closed the connection", client.Proto.RemoteAddr().String())
				break
			} else {
				log.Error(err.Error())
				break
			}
		}

		go this.OnRead(data, client)
	}
	client.Close()
}

func (this *ClientHandler) OnRead(input []byte, client *Client) {
	q := &proto.QueryStruct{}
	err := proto1.Unmarshal([]byte(input), q)
	if err != nil {
		log.Error("parse query error")
	}
	log.Info("sql: %s\npool: %s", q.GetSql(), q.GetPool())
	cols, rows, err := core.MysqlClusterInstance.Query(q)
	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info(cols)
		log.Info(rows)
		client.WriteMsg(fmt.Sprintf("%s\n", rows))
	}
}

func handleClient(client *server.Client) {
	client.Conn.Write([]byte("connected"))
}

type Client struct {
	*server.Client
	proto *server.Protocol
}

func newClient(c *server.Client, proto *server.Protocol) *Client {
	this := new(Client)
	this.Client = c
	this.proto = proto

	return this
}
