package main

import (
	"flag"

	"github.com/c4pt0r/rv"
)

var (
	configFile = flag.String("c", "config.toml", "config file (toml)")
	verbose    = flag.Bool("v", false, "verbose")
)

func main() {
	flag.Parse()

	rv.SetVerbose(*verbose)

	r := rv.NewServer(*configFile)
	r.Serve()
}
