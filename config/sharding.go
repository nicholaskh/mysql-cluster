package config

import (
	conf "github.com/nicholaskh/jsconf"
)

type ShardingConfig struct {
	TableName string
	Pk        string
	Privacy   *PrivacyConfig
}

func (this *ShardingConfig) loadConfig(cf *conf.Conf) {
	this.TableName = cf.String("table_name", "")
	this.Pk = cf.String("pk", "")
	this.Privacy = new(PrivacyConfig)
	section, err := cf.Section("privacy")
	if err != nil {
		panic("error occurred when parsing privacy config")
	}
	this.Privacy.loadConfig(section)

	if this.TableName == "" {
		panic("table name nil in sharding config")
	}
}

type PrivacyConfig struct {
	Name     string
	Behavior interface{}
}

func (this *PrivacyConfig) loadConfig(cf *conf.Conf) {
	this.Name = cf.String("name", "")
	this.Behavior = cf.Interface("behavior", nil)

	if this.Name == "" || this.Behavior == nil {
		panic("name nil in privacy config")
	}
}
