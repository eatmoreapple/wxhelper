package main

import (
	"github.com/eatmoreapple/env"
	"github.com/eatmoreapple/wxhelper/apiserver"
	"log"
)

func main() {
	srv := apiserver.Default()
	log.Fatal(srv.Run(env.Name("RUN_PORT").StringOrElse(":19089")))
}
