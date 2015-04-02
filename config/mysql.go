package config

import (
	"bytes"
	"fmt"
	"time"

	conf "github.com/nicholaskh/jsconf"
)

type MysqlConfig struct {
	ConnTimeout  time.Duration
	Connections  int
	MaxStmtCache int

	FailureAllowance uint
	RetryTimeout     time.Duration

	MaxIdleConns int
	MaxOpenConns int

	Pools map[string]map[int]*MysqlInstanceConfig
}

func (this *MysqlConfig) loadConfig(cf *conf.Conf) {
	this.ConnTimeout = cf.Duration("conn_timeout", time.Second*5)
	this.Connections = cf.Int("connections", 3)
	this.MaxStmtCache = cf.Int("max_stmt_cache", 20000)

	this.FailureAllowance = uint(cf.Int("failure_allowance", 5))
	this.RetryTimeout = cf.Duration("retry_timeout", time.Second*10)

	this.MaxIdleConns = cf.Int("max_idle_conns", 10)
	this.MaxOpenConns = cf.Int("max_open_conns", 15)

	this.Pools = make(map[string]map[int]*MysqlInstanceConfig)
	for i, _ := range cf.List("servers", nil) {
		section, _ := cf.Section(fmt.Sprintf("servers[%d]", i))
		mysqlInstanceConfig := &MysqlInstanceConfig{ConnTimeout: this.ConnTimeout}
		mysqlInstanceConfig.loadConfig(section)
		if _, exists := this.Pools[mysqlInstanceConfig.Pool]; !exists {
			this.Pools[mysqlInstanceConfig.Pool] = make(map[int]*MysqlInstanceConfig)
		}
		this.Pools[mysqlInstanceConfig.Pool][mysqlInstanceConfig.ShardId] = mysqlInstanceConfig
	}
}

type MysqlInstanceConfig struct {
	ConnTimeout time.Duration

	Pool    string
	Host    string
	Port    string
	User    string
	Pass    string
	ShardId int
	Db      string
	Charset string

	dsn string
}

func (this *MysqlInstanceConfig) loadConfig(cf *conf.Conf) {
	this.Pool = cf.String("pool", "")
	this.Host = cf.String("host", "")
	this.Port = cf.String("port", "3306")
	this.User = cf.String("user", "root")
	this.Pass = cf.String("pass", "")
	this.ShardId = cf.Int("shard_id", 0)
	this.Charset = cf.String("charset", "utf8")

	if this.ShardId != 0 {
		this.Db = fmt.Sprintf("%s%d", this.Pool, this.ShardId)
	} else {
		this.Db = this.Pool
	}

	if this.Pool == "" ||
		this.Host == "" ||
		this.Port == "" ||
		this.Db == "" {
		panic("miss required field")
	}

	buff := bytes.NewBuffer([]byte{})
	if this.User != "" {
		buff.WriteString(fmt.Sprintf("%s:", this.User))
		if this.Pass != "" {
			buff.WriteString(this.Pass)
		}
	}
	buff.WriteString(fmt.Sprintf("@tcp(%s:%s)/%s?", this.Host, this.Port, this.Db))
	buff.WriteString("autocommit=true") // we are not using transaction
	buff.WriteString(fmt.Sprintf("&timeout=%s", this.ConnTimeout))
	if this.Charset != "utf8" { // driver default utf-8
		buff.WriteString(fmt.Sprintf("&charset=%s", this.Charset))
	}
	buff.WriteString("&parseTime=true") // parse db timestamp automatically

	this.dsn = buff.String()
}

func (this *MysqlInstanceConfig) String() string {
	return this.DSN()
}

func (this *MysqlInstanceConfig) DSN() string {
	return this.dsn
}
