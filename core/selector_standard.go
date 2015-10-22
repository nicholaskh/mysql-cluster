package core

import (
	"github.com/nicholaskh/golib/sqlparser"
	log "github.com/nicholaskh/log4go"
	"github.com/nicholaskh/mysql-cluster/config"
)

type StandardSelector struct {
	config     *config.StandardShardingConfig
	serverPool ServerPool
}

func newStandardSelector(config *config.StandardShardingConfig, serverPool ServerPool) *StandardSelector {
	this := new(StandardSelector)
	this.config = config
	this.serverPool = serverPool

	return this
}

func (this *StandardSelector) PickServer(pool string, hintId int, sql string) (server *mysql, ex error) {
	if _, exists := this.serverPool[pool]; !exists {
		ex = ErrInvalidPool
		return
	}

	tree, err := sqlparser.Parse(sql)
	if err != nil {
		ex = ErrInvalidSql
		return
	}

	switch treeInst := tree.(type) {
	case *sqlparser.Select:
		//tableName := sqlparser.GetTableName(treeInst.From[0].(*sqlparser.AliasedTableExpr).Expr)
		//colName := sqlparser.GetColName(treeInst.Where.Expr.(*sqlparser.ComparisonExpr).Left)
		//colValue := string(treeInst.Where.Expr.(*sqlparser.ComparisonExpr).Right.(sqlparser.NumVal))

		var (
			shardingConfig *config.StandardShardingServerConfig
			exists         bool
		)
		if shardingConfig, exists = this.config.Servers[pool]; !exists {
			return
		}

		// TODO -- judge the operate
		if hintId > 0 {
			switch shardingConfig.Strategy.Name {
			case "range":
				var (
					listBehavior []interface{}
					ok           bool
				)
				if listBehavior, ok = shardingConfig.Strategy.Behavior.([]interface{}); !ok {
					ex = ErrInvalidStrategyBehavior
					return
				}
				// TODO -- O(n) room for optimize
				var (
					idx        int
					breakPoint interface{}
				)
				for idx, breakPoint = range listBehavior {
					if hintId < int(breakPoint.(float64)) {
						server, ex = this.ServerByBucket(pool, idx)
						return
					}
				}
				server, ex = this.ServerByBucket(pool, idx+1) // now idx points the last break point
				return

			case "custom":
				log.Info("custom sharding privacy")

			default:
				ex = ErrInvalidShardingStrategy
				log.Error("%s: %s", ex.Error(), shardingConfig.Strategy.Name)
				return
			}
		} else {
			// query all shards
			return
		}

	default:
		ex = ErrInvalidSqlType
		log.Error("%s: %s", ex.Error(), treeInst)
		return
	}

	// never arrive here
	ex = ErrUnknown
	return
}

func (this *StandardSelector) ServerByBucket(pool string, shardId int) (*mysql, error) {
	server, exists := this.serverPool[pool][shardId]
	if !exists {
		return nil, ErrServerNotFound
	}
	return server, nil
}
