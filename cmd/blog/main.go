package main

import (
	"flag"

	"github.com/narcissus1949/narcissus-blog/cmd/blog/app"
)

var confPath = flag.String("conf", "conf/conf.yaml", "config file path")

// @title           narcissus-blog Swagger API
// @version         1.0
// @description     narcissus-blog Swagger API

// @host      localhost:9090
// @BasePath
func main() {
	flag.Parse()
	app.Run(*confPath)
}
