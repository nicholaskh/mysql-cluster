package config

import (
	conf "github.com/nicholaskh/jsconf"
)

var Config struct {
	Server *ServerConfig
	Mysql  *MysqlConfig
}

func LoadConfig(cf *conf.Conf) {
	Config.Server = new(ServerConfig)
	section, err := cf.Section("server")
	if err != nil {
		panic(err)
	}
	Config.Server.loadConfig(section)

	Config.Mysql = new(MysqlConfig)
	section, err = cf.Section("mysql")
	if err != nil {
		panic(err)
	}
	Config.Mysql.loadConfig(section)
}
