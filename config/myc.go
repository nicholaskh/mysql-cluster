package config

import conf "github.com/nicholaskh/jsconf"

type MycConfig struct {
	Server           *ServerConfig
	Mysql            *MysqlConfig
	StandardSharding *StandardShardingConfig
	VbucketSharding  *VbucketShardingConfig
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

	this.StandardSharding = new(StandardShardingConfig)
	section, err = cf.Section("standard_sharding")
	if err != nil {
		panic(err)
	}
	this.StandardSharding.loadConfig(section)

	this.VbucketSharding = new(VbucketShardingConfig)
	section, err = cf.Section("vbucket_sharding")
	if err != nil {
		panic(err)
	}
	this.VbucketSharding.loadConfig(section)
}
