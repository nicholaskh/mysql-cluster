package config

import (
	conf "github.com/nicholaskh/jsconf"
)

var Config struct {
	ListenAddr string
}

func LoadConfig(cf *conf.Conf) {
	Config.ListenAddr = cf.String("listen_addr", ":3253")
}
