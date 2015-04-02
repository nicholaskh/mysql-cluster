package main

import "flag"

var (
	options struct {
		showVersion  bool
		configFile   string
		logFile      string
		logLevel     string
		crashLogFile string
	}
)

func parseFlags() {
	flag.BoolVar(&options.showVersion, "v", false, "show version")
	flag.StringVar(&options.configFile, "conf", "etc/myc.cf", "config file path")
	flag.StringVar(&options.logFile, "log", "stdout", "log file")
	flag.StringVar(&options.logLevel, "level", "info", "log level")
	flag.StringVar(&options.crashLogFile, "crashLog", "panic.dump", "crash log file")

	flag.Parse()
}
