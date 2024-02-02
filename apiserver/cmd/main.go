package main

import (
	"github.com/eatmoreapple/wxhelper/apiserver"
	"log"
)

const defaultADDR = ":19089"

func main() {
	srv := apiserver.Default()
	log.Fatal(srv.Run(defaultADDR))
}
