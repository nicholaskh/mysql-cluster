package core

import (
	"strconv"

	"github.com/nicholaskh/golib/sqlparser"
	log "github.com/nicholaskh/log4go"
	"github.com/nicholaskh/mysql-cluster/config"
)

type Selector struct {
	config map[string]*config.ShardingConfig
}

func NewSelector(config map[string]*config.ShardingConfig) *Selector {
	this := new(Selector)
	this.config = config

	return this
}

func (this *Selector) LookupShardId(sql string) (shardId int, ex error) {
	tree, err := sqlparser.Parse(sql)
	if err != nil {
		shardId = 0
		ex = ErrInvalidSql
		return
	}

	switch treeInst := tree.(type) {
	case *sqlparser.Select:
		tableName := sqlparser.GetTableName(treeInst.From[0].(*sqlparser.AliasedTableExpr).Expr)
		log.Info("table name: %s", tableName)
		colName := sqlparser.GetColName(treeInst.Where.Expr.(*sqlparser.ComparisonExpr).Left)
		colValue := string(treeInst.Where.Expr.(*sqlparser.ComparisonExpr).Right.(sqlparser.NumVal))
		log.Info("col name: %s", colName)
		log.Info("col value: %s", colValue)

		var (
			shardingConfig *config.ShardingConfig
			exists         bool
		)
		if shardingConfig, exists = this.config[tableName]; !exists {
			return
		}

		// TODO -- judge the operate
		if colName == shardingConfig.Pk {
			switch shardingConfig.Privacy.Name {
			case "range":
				intColValue, _ := strconv.Atoi(colValue)
				var (
					listBehavior []interface{}
					ok           bool
				)
				if listBehavior, ok = shardingConfig.Privacy.Behavior.([]interface{}); !ok {
					ex = ErrInvalidPrivacyBehavior
					return
				}
				// TODO -- O(n) room for optimize
				var (
					idx        int
					breakPoint interface{}
				)
				for idx, breakPoint = range listBehavior {
					if intColValue < int(breakPoint.(float64)) {
						shardId = idx + 1
						return
					}
				}
				return
				//shardId = intColValue / int(shardingConfig.Privacy.Behavior.(float64))
				return

			case "custom":
				log.Info("custom sharding privacy")

			default:
				ex = ErrInvalidShardingPrivacy
				log.Error("%s: %s", ex.Error(), shardingConfig.Privacy.Name)
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
	return
}
