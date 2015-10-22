package config

import (
	conf "github.com/nicholaskh/jsconf"
)

type VbucketShardingConfig struct {
	VbucketBaseNumber int
}

func (this *VbucketShardingConfig) loadConfig(cf *conf.Conf) {
	this.VbucketBaseNumber = cf.Int("vbucket_base_number", 1024)
}
