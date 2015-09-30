package main

import (
	"flag"
	"fmt"
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
		args        strSlice
		readTimeout time.Duration
	}
)

type strSlice []string

func (this *strSlice) String() string {
	return fmt.Sprintf("%s", *this)
}

func (this *strSlice) Set(value string) error {
	*this = append(*this, value)
	return nil
}

func parseFlags() {
	flag.StringVar(&options.addr, "addr", ":3253", "server address")
	flag.StringVar(&options.pool, "pool", "test", "shard pool to query")
	flag.StringVar(&options.sql, "sql", "", "sql to query")
	flag.Var(&options.args, "args", "args to execute sql")
	flag.DurationVar(&options.readTimeout, "read_timeout", time.Second*5, "read timeout")

	flag.Parse()
}

func main() {
	parseFlags()

	server.SetupLogging("stdout", "debug", "")

	c := NewClient(options.readTimeout)
	log.Info(options.sql)
	c.Dial(options.addr)
	result, err := c.Query(options.pool, options.sql, options.args)

	if err != nil {
		log.Error(err)
	} else {
		log.Info(result)
	}

	c.Close()

	time.Sleep(time.Second)
}
