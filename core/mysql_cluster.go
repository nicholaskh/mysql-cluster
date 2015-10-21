package core

import (
	"fmt"

	sql_ "database/sql"
	log "github.com/nicholaskh/log4go"
	"github.com/nicholaskh/mysql-cluster/config"
	"github.com/nicholaskh/mysql-cluster/proto/go"
)

var MysqlClusterInstance *MysqlCluster

type MysqlCluster struct {
	pools    map[string]*mysql //[db(pool + shardId) => mysql]
	selector *Selector
	config   *config.MycConfig
}

func NewMysqlCluster(config *config.MycConfig) *MysqlCluster {
	this := new(MysqlCluster)

	this.pools = make(map[string]*mysql)

	for pool, mysqlMap := range config.Mysql.Pools {
		for shardId, mysqlInstanceConfig := range mysqlMap {
			mysql := newMysql(mysqlInstanceConfig.DSN(), config.Mysql)
			if shardId != 0 {
				this.pools[fmt.Sprintf("%s%d", pool, shardId)] = mysql
			} else {
				this.pools[pool] = mysql
			}
			mysql.Open()
		}
	}

	this.selector = NewSelector(config.Sharding)
	this.config = config
	return this
}

func (this *MysqlCluster) Query(q *proto.QueryStruct) (cols []string, rows [][]string, err error) {
	rows = make([][]string, 0)
	var (
		rawRowValues []sql_.RawBytes
		scanArgs     []interface{}
		rowValues    []string
	)

	pool := q.GetPool()
	sql := q.GetSql()

	var shardId int
	shardId, err = this.selector.LookupShardId(sql)
	if err != nil {
		return
	}

	args := q.GetArgs()
	argsI := make([]interface{}, len(args))

	for i, arg := range args {
		argsI[i] = arg
	}

	var (
		server *mysql
		exists bool
	)
	log.Info("server id: %s%d", pool, shardId)
	if server, exists = this.pools[fmt.Sprintf("%s%d", pool, shardId)]; !exists {
		err = ErrMysqlServerNotFound
		return
	}

	rs, err := server.Query(sql, argsI...)
	if err != nil {
		log.Error(err)
		return
	}
	defer rs.Close()

	// initialize the vars only once
	cols, err = rs.Columns()
	if err != nil {
		rs.Close()
		return
	}

	rawRowValues = make([]sql_.RawBytes, len(cols))
	scanArgs = make([]interface{}, len(cols))
	for i, _ := range cols {
		scanArgs[i] = &rawRowValues[i]
	}

	for rs.Next() {
		if err := rs.Scan(scanArgs...); err != nil {
			break
		}

		rowValues = make([]string, len(cols))
		// TODO O(N), room for optimization, allow_nullable_columns
		for i, raw := range rawRowValues {
			if raw == nil {
				rowValues[i] = "NULL"
			} else {
				rowValues[i] = string(raw)
			}
		}

		rows = append(rows, rowValues)
	}
	return
}
