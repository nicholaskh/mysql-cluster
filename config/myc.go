package config

import (
	"fmt"

	conf "github.com/nicholaskh/jsconf"
)

type MycConfig struct {
	Server   *ServerConfig
	Mysql    *MysqlConfig
	Sharding map[string]*ShardingConfig
}

func (this *MycConfig) LoadConfig(cf *conf.Conf) {
	this.Server = new(ServerConfig)
	section, err := cf.Section("server")
	if err != nil {
		panic(err)
	}
	this.Server.loadConfig(section)

	this.Mysql = new(MysqlConfig)
	section, err = cf.Section("mysql")
	if err != nil {
		panic(err)
	}
	this.Mysql.loadConfig(section)

	this.Sharding = make(map[string]*ShardingConfig)
	for i, _ := range cf.List("sharding", nil) {
		section, err = cf.Section(fmt.Sprintf("sharding[%d]", i))
		shardingConfig := &ShardingConfig{}
		shardingConfig.loadConfig(section)
		this.Sharding[shardingConfig.TableName] = shardingConfig
	}
}
