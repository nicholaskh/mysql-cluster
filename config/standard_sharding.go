package config

import (
	"fmt"

	conf "github.com/nicholaskh/jsconf"
)

type StandardShardingConfig struct {
	Servers map[string]*StandardShardingServerConfig
}

func (this *StandardShardingConfig) loadConfig(cf *conf.Conf) {
	this.Servers = make(map[string]*StandardShardingServerConfig)
	for i, _ := range cf.List("servers", nil) {
		section, _ := cf.Section(fmt.Sprintf("servers[%d]", i))
		standardShardingServerConfig := &StandardShardingServerConfig{}
		standardShardingServerConfig.loadConfig(section)
		this.Servers[standardShardingServerConfig.Pool] = standardShardingServerConfig
	}
}

type StandardShardingServerConfig struct {
	Pool     string
	Strategy *StrategyConfig
}

func (this *StandardShardingServerConfig) loadConfig(cf *conf.Conf) {
	this.Pool = cf.String("pool", "")
	this.Strategy = new(StrategyConfig)
	section, err := cf.Section("strategy")
	if err != nil {
		panic("error occurred when parsing strategy config")
	}
	this.Strategy.loadConfig(section)

	if this.Pool == "" {
		panic("pool nil in sharding config")
	}
}

type StrategyConfig struct {
	Name     string
	Behavior interface{}
}

func (this *StrategyConfig) loadConfig(cf *conf.Conf) {
	this.Name = cf.String("name", "")
	this.Behavior = cf.Interface("behavior", nil)

	if this.Name == "" || this.Behavior == nil {
		panic("name nil in strategy config")
	}
}
