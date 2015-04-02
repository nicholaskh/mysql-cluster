package proxygate

import (
	"fmt"

	"github.com/nicholaskh/mysql-cluster/config"
	"github.com/nicholaskh/mysql-cluster/proto/go"
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
			if shardId != 0 {
				this.pools[fmt.Sprintf("%s%d", pool, shardId)] = mysql
			} else {
				this.pools[pool] = mysql
			}
			mysql.Open()
		}
	}
	return this
}

func (this *ProxyGate) Execute(q *proto.QueryStruct) (res string, err error) {
	pool := q.Getpool()
	sql := q.Getsql()
	args := q.Getallargs()
	argsI := make([]interface{}, len(args))

	for i, arg := range args {
		argsI[i] = arg
	}

	rows, err := this.pools[pool].Query(sql, argsI...)
	defer rows.Close()
	for rows.Next() {
		var id string
		var name string
		rows.Scan(&id, &name)
		res += fmt.Sprintf("%s %s\n", id, name)
	}
	return
}
