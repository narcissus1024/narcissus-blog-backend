package main

import (
	"flag"

	"github.com/narcissus1949/narcissus-blog/cmd"
)

var confPath = flag.String("conf", "conf/conf.yaml", "config file path")

func main() {
	flag.Parse()
	cmd.Run(*confPath)
}
