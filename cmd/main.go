package main

import (
	"flag"
	"runtime"

	"github.com/c4pt0r/rv"
	log "github.com/ngaut/logging"
)

var (
	configFile = flag.String("c", "config.toml", "config file (toml)")
	logLevel   = flag.String("loglevel", "error", "loglevel")
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.SetLevelByString(*logLevel)
	r := rv.NewServer(*configFile)
	r.Serve()
}
