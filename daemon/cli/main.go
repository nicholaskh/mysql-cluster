package main

import (
	"flag"
	"time"

	"github.com/nicholaskh/golib/server"
	log "github.com/nicholaskh/log4go"
	. "github.com/nicholaskh/mysql-cluster/proto/go"
)

var (
	options struct {
		addr        string
		pool        string
		sql         string
		readTimeout time.Duration
	}
)

func parseFlags() {
	flag.StringVar(&options.addr, "addr", ":3253", "server address")
	flag.StringVar(&options.pool, "pool", "test", "shard pool to query")
	flag.StringVar(&options.sql, "sql", "", "sql to query")
	flag.DurationVar(&options.readTimeout, "read_timeout", time.Second*5, "read timeout")

	flag.Parse()
}

func main() {
	parseFlags()

	server.SetupLogging("stdout", "debug", "")

	c := NewClient(options.readTimeout)
	c.Dial(options.addr)
	result, _ := c.Query(options.pool, options.sql)

	log.Info(result)

	c.Close()

	time.Sleep(time.Second * 10)
}
