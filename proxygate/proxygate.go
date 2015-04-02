package proxygate

import (
	"fmt"

	"github.com/nicholaskh/mysql-cluster/config"
)

type ProxyGate struct {
	pools map[string]*mysql //[db(pool + shardId) => mysql]
}

func NewProxyGate() *ProxyGate {
	this := new(ProxyGate)
	this.pools = make(map[string]*mysql)

	for pool, mysqlMap := range config.Config.Mysql.Pools {
		for shardId, mysqlInstanceConfig := range mysqlMap {
			mysql := newMysql(mysqlInstanceConfig.DSN(), config.Config.Mysql.MaxStmtCache,
				config.Config.Mysql.MaxIdleConns, config.Config.Mysql.MaxOpenConns)
			this.pools[fmt.Sprintf("%s%d", pool, shardId)] = mysql
		}
	}
	return this
}

func (this *ProxyGate) Execute(sql string, params ...interface{}) {

}
