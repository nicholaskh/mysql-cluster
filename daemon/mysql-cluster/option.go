package main

import (
	"flag"
)

var (
	options struct {
		showVersion bool
		configFile  string
	}
)

func parseFlags() {
	flag.BoolVar(&options.showVersion, "v", false, "show version")
	flag.StringVar(&options.configFile, "conf", "etc/myc.cf", "config file path")

	flag.Parse()
}
