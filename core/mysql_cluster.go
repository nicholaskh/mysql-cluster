package core

import (
	sql_ "database/sql"

	log "github.com/nicholaskh/log4go"
	"github.com/nicholaskh/mysql-cluster/config"
	"github.com/nicholaskh/mysql-cluster/proto/go"
)

type ServerPool map[string]map[int]*mysql // [pool => [shardId => mysql]]

var MysqlClusterInstance *MysqlCluster

type MysqlCluster struct {
	serverPool ServerPool
	selector   Selector
	config     *config.MycConfig
}

func NewMysqlCluster(config *config.MycConfig) *MysqlCluster {
	this := new(MysqlCluster)

	this.serverPool = make(ServerPool)

	for pool, mysqlMap := range config.Mysql.Pools {
		for shardId, mysqlInstanceConfig := range mysqlMap {
			server := newMysql(mysqlInstanceConfig.DSN(), config.Mysql)
			if _, exists := this.serverPool[pool]; !exists {
				this.serverPool[pool] = make(map[int]*mysql)
			}
			if shardId >= 0 {
				this.serverPool[pool][shardId] = server
			} else {
				this.serverPool[pool][0] = server
			}

			// TODO -- open and ping for {{retries}} times
			server.Open()
		}
	}

	switch config.Mysql.ShardingStrategy {
	case "standard":
		this.selector = newStandardSelector(config.StandardSharding, this.serverPool)

	case "vbucket":
		this.selector = newVbucketSelector(config.VbucketSharding, this.serverPool)

	default:
		panic("invalid sharding strategy")
	}

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
	//hintId := q.GetHintId()
	// FIXME -- delete
	hintId := 3
	sql := q.GetSql()

	var server *mysql
	server, err = this.selector.PickServer(pool, hintId, sql)
	if err != nil {
		return
	}

	args := q.GetArgs()
	argsI := make([]interface{}, len(args))

	for i, arg := range args {
		argsI[i] = arg
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
