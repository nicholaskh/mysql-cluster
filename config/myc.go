package config

import (
	conf "github.com/nicholaskh/jsconf"
)

var Config struct {
	ListenAddr string
	Mysql      *MysqlConfig
}

func LoadConfig(cf *conf.Conf) {
	Config.ListenAddr = cf.String("listen_addr", ":3253")

	Config.Mysql = new(MysqlConfig)
	section, err := cf.Section("mysql")
	if err != nil {
		panic(err)
	}
	Config.Mysql.loadConfig(section)
}
