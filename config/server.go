package config

import (
	conf "github.com/nicholaskh/jsconf"
	"time"
)

type ServerConfig struct {
	ListenAddr  string
	SessTimeout time.Duration
}

func (this *ServerConfig) loadConfig(cf *conf.Conf) {
	this.ListenAddr = cf.String("listen_addr", ":3253")
	this.SessTimeout = cf.Duration("sess_timeout", time.Second*5)
}
